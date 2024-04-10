package api

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gre-ory/games-go/internal/game/share/model"
	"github.com/gre-ory/games-go/internal/util"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// cookie server

type CookieCallback func(cookie *model.Cookie)

type CookieServer interface {
	util.Server
	NewCookie() *model.Cookie
	GetCookieOrDefault(r *http.Request) *model.Cookie
	GetCookie(r *http.Request) (*model.Cookie, error)
	GetValidCookie(r *http.Request) (*model.Cookie, error)
	SetCookie(w http.ResponseWriter, cookie *model.Cookie) error
	ClearCookie(w http.ResponseWriter) error
}

// //////////////////////////////////////////////////
// constructor

func NewCookieServer(logger *zap.Logger, key string, maxAge int, cookieSecret string, cookieCallback CookieCallback) CookieServer {

	// TODO replace by proto
	empty := &model.Cookie{}
	gob.Register(empty)

	// encrypter
	block, err := aes.NewCipher([]byte(cookieSecret))
	if err != nil {
		panic(err)
	}
	encrypter, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	return &cookieServer{
		logger:    logger.With(zap.String("cookie", key)),
		key:       key,
		maxAge:    maxAge,
		encrypter: encrypter,
		onCookie:  cookieCallback,
		hxServer:  util.NewHxServer(logger, ShareTpl),
	}
}

type cookieServer struct {
	logger    *zap.Logger
	key       string
	maxAge    int
	encrypter cipher.AEAD
	onCookie  CookieCallback
	hxServer  util.HxServer
}

// //////////////////////////////////////////////////
// register routes

func (s *cookieServer) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/htmx/user", s.htmx_get_user)
	router.HandlerFunc(http.MethodPut, "/htmx/user", s.htmx_set_user)
	router.HandlerFunc(http.MethodGet, "/htmx/user-avatar-modal", s.htmx_user_avatar_modal)
	router.HandlerFunc(http.MethodGet, "/htmx/user-name-modal", s.htmx_user_name_modal)
	router.HandlerFunc(http.MethodGet, "/htmx/user-language-modal", s.htmx_user_language_modal)
}

// //////////////////////////////////////////////////
// get cookie

func (s *cookieServer) NewCookie() *model.Cookie {
	return model.NewCookie()
}

func (s *cookieServer) GetCookie(r *http.Request) (*model.Cookie, error) {

	cookie, err := r.Cookie(s.key)
	if err != nil {
		s.logger.Info(fmt.Sprintf("unable to get cookie: %s", err.Error()))
		return nil, err
	}
	cookieBase64 := cookie.Value
	// c.logger.Info("get cookie", zap.String("cookie-base64", cookieBase64), zap.Any("cookie", cookie))

	cookieEncrypted, err := s.decodeBase64(cookieBase64)
	if err != nil {
		s.logger.Info(fmt.Sprintf("unable to base64 decode: %s", err.Error()), zap.String("value", cookie.Value))
		return nil, err
	}
	// c.logger.Info("decode 64", zap.Binary("cookie-encrypted", cookieEncrypted))

	cookieEncoded, err := s.decrypt(cookieEncrypted)
	if err != nil {
		s.logger.Info(fmt.Sprintf("unable to decrypt: %s", err.Error()))
		return nil, err
	}
	// c.logger.Info("decrypt", zap.Binary("cookie-encoded", cookieEncoded))

	value, err := s.decode(cookieEncoded)
	if err != nil {
		s.logger.Info(fmt.Sprintf("unable to decode: %s", err.Error()), zap.Binary("cookie-encoded", cookieEncoded))
		return nil, err
	}
	// s.logger.Info("get cookie", zap.Any("value", value), zap.Any("cookie", cookie))

	// sanitize
	value.Sanitize()

	return value, nil
}

func (s *cookieServer) GetCookieOrDefault(r *http.Request) *model.Cookie {
	cookie, err := s.GetCookie(r)
	if err != nil {
		s.logger.Info("cookie not found >>> create default one!", zap.Error(err))
		cookie = s.NewCookie()
	}
	if err := cookie.Validate(); err != nil {
		s.logger.Info("invalid cookie >>> create default one!", zap.Error(err))
		cookie = s.NewCookie()
	}
	return cookie
}

func (s *cookieServer) GetValidCookie(r *http.Request) (*model.Cookie, error) {
	cookie, err := s.GetCookie(r)
	if err != nil {
		return nil, err
	}
	if err := cookie.Validate(); err != nil {
		return nil, err
	}
	return cookie, nil
}

// //////////////////////////////////////////////////
// set cookie

func (s *cookieServer) SetCookie(w http.ResponseWriter, cookie *model.Cookie) error {

	cookieEncoded, err := s.encode(cookie)
	if err != nil {
		s.logger.Warn("unable to encode", zap.Error(err))
		return err
	}
	// c.logger.Info("encode", zap.Any("cookie-value", value), zap.Binary("cookie-encoded", cookieEncoded))

	cookieEncrypted, err := s.encrypt(cookieEncoded)
	if err != nil {
		s.logger.Warn("unable to encrypt cookie", zap.Error(err))
		return err
	}
	// c.logger.Info("encrypt", zap.Binary("cookie-encrypted", cookieEncrypted))

	cookieBase64 := s.encodeBase64(cookieEncrypted)
	// c.logger.Info("base64", zap.String("cookie-base64", cookieBase64))

	err = s.validateCookieValue(cookieBase64)
	if err != nil {
		s.logger.Warn("invalid cookie value", zap.Error(err))
		return err
	}

	httpCookie := s.newCookie(cookieBase64, s.maxAge)
	s.logger.Info("set cookie", zap.Any("cookie", cookie), zap.Any("http-cookie", httpCookie))
	http.SetCookie(w, httpCookie)
	return nil
}

// //////////////////////////////////////////////////
// clear cookie

func (s *cookieServer) ClearCookie(w http.ResponseWriter) error {
	cookie := s.newCookie("", 0)
	s.logger.Info("clear cookie", zap.String("value", cookie.Value), zap.Int("age", cookie.MaxAge))
	http.SetCookie(w, cookie)
	return nil
}

// //////////////////////////////////////////////////
// new cookie

func (s *cookieServer) newCookie(value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     s.key,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   false, // TODO: support HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

// //////////////////////////////////////////////////
// validate

var (
	ErrValueTooLong = errors.New("cookie value is too long")
)

func (c *cookieServer) validateCookieValue(value string) error {
	if len(value) > 4096 {
		return ErrValueTooLong
	}
	return nil
}

// //////////////////////////////////////////////////
// encoder

func (s *cookieServer) encode(value *model.Cookie) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(value)
	if err != nil {
		return nil, fmt.Errorf("unable to gob encode: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *cookieServer) decode(encoded []byte) (*model.Cookie, error) {
	value := &model.Cookie{}
	reader := bytes.NewReader(encoded)
	// reader := strings.NewReader(encodedValue)
	if err := gob.NewDecoder(reader).Decode(value); err != nil {
		return nil, fmt.Errorf("unable to gob decode: %w", err)
	}
	return value, nil
}

// //////////////////////////////////////////////////
// base 64

var (
	ErrInvalidBase64Value = fmt.Errorf("invalid base64 value")
)

func (s *cookieServer) encodeBase64(value []byte) string {
	return base64.URLEncoding.EncodeToString(value)
}

func (s *cookieServer) decodeBase64(value string) ([]byte, error) {
	encryptedValue, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		s.logger.Warn("base64 decode: error", zap.Error(err))
		return nil, ErrInvalidBase64Value
	}
	return encryptedValue, nil
}

// //////////////////////////////////////////////////
// encrypter

var (
	ErrNonceSize             = fmt.Errorf("invalid nonce size")
	ErrInvalidEncryptedValue = fmt.Errorf("invalid encrypted value")
	ErrInvalidEncryptedKey   = fmt.Errorf("invalid encrypted key")
)

func (s *cookieServer) encrypt(value []byte) ([]byte, error) {
	// decryptedValue := fmt.Sprintf("%s:%s", c.key, value)
	// encryptedValue := c.encrypter.Seal(c.nonce, c.nonce, []byte(decryptedValue), nil)
	// c.logger.Info("encrypt", zap.Binary("value", value), zap.String("decryptedValue", decryptedValue), zap.String("encryptedValue", string(encryptedValue)))

	// nonce
	nonce := make([]byte, s.encrypter.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	encryptedValue := s.encrypter.Seal(nonce, nonce, value, nil)
	return encryptedValue, nil
}

func (s *cookieServer) decrypt(value []byte) ([]byte, error) {

	nonceSize := s.encrypter.NonceSize()
	if len(value) < nonceSize {
		return nil, ErrNonceSize
	}

	// Split apart the nonce from the actual encrypted data.
	nonce, encryptedValue := value[:nonceSize], value[nonceSize:]

	decryptedValue, err := s.encrypter.Open(nil, []byte(nonce), []byte(encryptedValue), nil)
	if err != nil {
		return nil, err
	}
	return decryptedValue, nil

	// values := strings.SplitN(string(decryptedValue), ":", 1)
	// if len(values) != 2 {
	// 	c.logger.Warn("decrypt: error", zap.Error(ErrInvalidEncryptedValue))
	// 	return nil, ErrInvalidEncryptedValue
	// }

	// if values[0] != c.key {
	// 	c.logger.Warn("decrypt: error", zap.Error(ErrInvalidEncryptedKey))
	// 	return "", ErrInvalidEncryptedKey
	// }
	// c.logger.Info("decrypt", zap.String("values[1]", values[1]), zap.String("decryptedValue", string(decryptedValue)), zap.String("encryptedValue", string(value)))
	// return values[1], nil
}

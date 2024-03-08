package util

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

	"go.uber.org/zap"
)

// //////////////////////////////////////////////////
// cookie helper

type CookieHelper[T any] interface {
	GetCookie(r *http.Request) (*T, error)
	GetCookieOrDefault(r *http.Request, defaultCookieFn func() *T) *T
	SetCookie(w http.ResponseWriter, value *T) error
	ClearCookie(w http.ResponseWriter) error
}

// //////////////////////////////////////////////////
// constructor

func NewCookieHelper[T any](logger *zap.Logger, key string, maxAge int, cookieSecret string) CookieHelper[T] {

	var empty T
	gob.Register(&empty)

	// encrypter
	block, err := aes.NewCipher([]byte(cookieSecret))
	if err != nil {
		panic(err)
	}
	encrypter, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	return &cookieHelper[T]{
		logger:    logger.With(zap.String("cookie", key)),
		key:       key,
		maxAge:    maxAge,
		encrypter: encrypter,
	}
}

type cookieHelper[T any] struct {
	logger    *zap.Logger
	key       string
	maxAge    int
	encrypter cipher.AEAD
}

// //////////////////////////////////////////////////
// get cookie

func (c *cookieHelper[T]) GetCookie(r *http.Request) (*T, error) {

	cookie, err := r.Cookie(c.key)
	if err != nil {
		c.logger.Info(fmt.Sprintf("unable to get cookie: %s", err.Error()))
		return nil, err
	}
	cookieBase64 := cookie.Value
	// c.logger.Info("get cookie", zap.String("cookie-base64", cookieBase64), zap.Any("cookie", cookie))

	cookieEncrypted, err := c.decodeBase64(cookieBase64)
	if err != nil {
		c.logger.Info(fmt.Sprintf("unable to base64 decode: %s", err.Error()), zap.String("value", cookie.Value))
		return nil, err
	}
	// c.logger.Info("decode 64", zap.Binary("cookie-encrypted", cookieEncrypted))

	cookieEncoded, err := c.decrypt(cookieEncrypted)
	if err != nil {
		c.logger.Info(fmt.Sprintf("unable to decrypt: %s", err.Error()))
		return nil, err
	}
	// c.logger.Info("decrypt", zap.Binary("cookie-encoded", cookieEncoded))

	value, err := c.decode(cookieEncoded)
	if err != nil {
		c.logger.Info(fmt.Sprintf("unable to decode: %s", err.Error()), zap.Binary("cookie-encoded", cookieEncoded))
		return nil, err
	}
	c.logger.Info("get cookie", zap.Any("value", value), zap.Any("cookie", cookie))

	return value, nil
}

func (c *cookieHelper[T]) GetCookieOrDefault(r *http.Request, defaultCookieFn func() *T) *T {
	cookie, err := c.GetCookie(r)
	if err != nil {
		c.logger.Info("cookie not found >>> create default one!", zap.Error(err))
		cookie = defaultCookieFn()
	}
	return cookie
}

// //////////////////////////////////////////////////
// set cookie

func (c *cookieHelper[T]) SetCookie(w http.ResponseWriter, value *T) error {

	cookieEncoded, err := c.encode(value)
	if err != nil {
		c.logger.Warn("unable to encode", zap.Error(err))
		return err
	}
	// c.logger.Info("encode", zap.Any("cookie-value", value), zap.Binary("cookie-encoded", cookieEncoded))

	cookieEncrypted, err := c.encrypt(cookieEncoded)
	if err != nil {
		c.logger.Warn("unable to encrypt cookie", zap.Error(err))
		return err
	}
	// c.logger.Info("encrypt", zap.Binary("cookie-encrypted", cookieEncrypted))

	cookieBase64 := c.encodeBase64(cookieEncrypted)
	// c.logger.Info("base64", zap.String("cookie-base64", cookieBase64))

	err = c.validateCookieValue(cookieBase64)
	if err != nil {
		c.logger.Warn("invalid cookie value", zap.Error(err))
		return err
	}

	cookie := c.newCookie(cookieBase64, c.maxAge)
	c.logger.Info("set cookie", zap.Any("value", value), zap.Any("cookie", cookie))
	http.SetCookie(w, cookie)
	return nil
}

// //////////////////////////////////////////////////
// clear cookie

func (c *cookieHelper[T]) ClearCookie(w http.ResponseWriter) error {
	cookie := c.newCookie("", 0)
	c.logger.Info("clear cookie", zap.String("value", cookie.Value), zap.Int("age", cookie.MaxAge))
	http.SetCookie(w, cookie)
	return nil
}

// //////////////////////////////////////////////////
// new cookie

func (c *cookieHelper[T]) newCookie(value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     c.key,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

// //////////////////////////////////////////////////
// validate

var (
	ErrValueTooLong = errors.New("cookie value is too long")
)

func (c *cookieHelper[T]) validateCookieValue(value string) error {
	if len(value) > 4096 {
		return ErrValueTooLong
	}
	return nil
}

// //////////////////////////////////////////////////
// encoder

func (c *cookieHelper[T]) encode(value *T) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(value)
	if err != nil {
		return nil, fmt.Errorf("unable to gob encode: %w", err)
	}
	return buf.Bytes(), nil
}

func (c *cookieHelper[T]) decode(encoded []byte) (*T, error) {
	value := new(T)
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

func (c *cookieHelper[T]) encodeBase64(value []byte) string {
	return base64.URLEncoding.EncodeToString(value)
}

func (c *cookieHelper[T]) decodeBase64(value string) ([]byte, error) {
	encryptedValue, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		c.logger.Warn("base64 decode: error", zap.Error(err))
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

func (c *cookieHelper[T]) encrypt(value []byte) ([]byte, error) {
	// decryptedValue := fmt.Sprintf("%s:%s", c.key, value)
	// encryptedValue := c.encrypter.Seal(c.nonce, c.nonce, []byte(decryptedValue), nil)
	// c.logger.Info("encrypt", zap.Binary("value", value), zap.String("decryptedValue", decryptedValue), zap.String("encryptedValue", string(encryptedValue)))

	// nonce
	nonce := make([]byte, c.encrypter.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	encryptedValue := c.encrypter.Seal(nonce, nonce, value, nil)
	return encryptedValue, nil
}

func (c *cookieHelper[T]) decrypt(value []byte) ([]byte, error) {

	nonceSize := c.encrypter.NonceSize()
	if len(value) < nonceSize {
		return nil, ErrNonceSize
	}

	// Split apart the nonce from the actual encrypted data.
	nonce, encryptedValue := value[:nonceSize], value[nonceSize:]

	decryptedValue, err := c.encrypter.Open(nil, []byte(nonce), []byte(encryptedValue), nil)
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

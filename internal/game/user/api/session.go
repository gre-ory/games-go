package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util"

	"github.com/gre-ory/games-go/internal/game/user/model"
	"github.com/gre-ory/games-go/internal/game/user/service"
)

// //////////////////////////////////////////////////
// session server

type SessionServer interface {
	util.Server
	SetSessionHeaders(w http.ResponseWriter, session *model.Session)
	ExtractSessionFromRequest(r *http.Request) (*model.Session, error)
	ExtractSessionToken(r *http.Request) (model.Token, error)
	ExtractSession(token model.Token) (*model.Session, error)
	AttachContext(token model.Token, context map[string]string) (*model.Session, error)
}

func NewSessionServer(logger *zap.Logger, service service.SessionService) SessionServer {
	return &sessionServer{
		logger:  logger,
		service: service,
	}
}

type sessionServer struct {
	logger  *zap.Logger
	service service.SessionService
}

// //////////////////////////////////////////////////
// register

func (s *sessionServer) RegisterRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPut, "/user/session/new", s.session_create)
}

// //////////////////////////////////////////////////
// headers

const (
	SessionTokenHeader  = "X-Session-Token"
	SessionExpireHeader = "X-Session-Expire"
	ExpireFormat        = "2006-01-02T15:04:05.999Z"
)

func (s *sessionServer) SetSessionHeaders(w http.ResponseWriter, session *model.Session) {
	if session == nil {
		return
	}
	if session.Token != "" {
		w.Header().Set(SessionTokenHeader, string(session.Token))
	}
	if !session.Expire.IsZero() {
		w.Header().Set(SessionExpireHeader, string(session.Expire.UTC().Format(ExpireFormat)))
	}
}

func (s *sessionServer) ExtractSessionToken(r *http.Request) (model.Token, error) {
	if value := r.Header.Get(SessionTokenHeader); value != "" {
		return model.Token(value), nil
	}
	return "", model.ErrSessionTokenNotFound
}

func (s *sessionServer) ExtractSessionFromRequest(r *http.Request) (*model.Session, error) {
	if value := r.Header.Get(SessionTokenHeader); value != "" {
		token := model.Token(value)
		session, err := s.service.Retrieve(token)
		if err != nil {
			return nil, err
		}
		return session, nil
	}
	return nil, model.ErrSessionTokenNotFound
}

func (s *sessionServer) ExtractSession(token model.Token) (*model.Session, error) {
	if token != "" {
		session, err := s.service.Retrieve(token)
		if err != nil {
			return nil, err
		}
		return session, nil
	}
	return nil, model.ErrSessionTokenNotFound
}

func (s *sessionServer) AttachContext(token model.Token, context map[string]string) (*model.Session, error) {
	return s.service.AttachContext(token, context)
}

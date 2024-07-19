package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

// //////////////////////////////////////////////////
// websocket user

type User interface {
	model.User

	HasGameId() bool
	GameId() model.GameId
	SetGameId(gameId model.GameId)
	UnsetGameId()

	PlayerId() model.PlayerId

	IsInactive() bool
	IsActive() bool
	IsNotPlaying() bool
	IsPlaying() bool

	ConnectSocket(w http.ResponseWriter, r *http.Request) error

	Activate()
	Deactivate()

	Send(bytes []byte)
	Close()
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewUser(
	logger *zap.Logger,
	cookie *model.Cookie,
	onMessage func(id model.UserId, message []byte),
	onUpdate func(id model.UserId),
	onClose func(id model.UserId),
) User {
	if cookie.Id == "" {
		panic(model.ErrInvalidCookie)
	}
	return &user{
		User:        model.NewUserFromCookie(cookie),
		logger:      logger.With(zap.String("user", string(cookie.Id))),
		active:      false,
		gameId:      "",
		onMessage:   onMessage,
		onUpdate:    onUpdate,
		onClose:     onClose,
		readClosed:  true,
		writeClosed: true,
		closing:     false,
		closed:      true,
	}
}

type user struct {
	sync.RWMutex
	model.User
	logger           *zap.Logger
	active           bool
	gameId           model.GameId
	send             chan []byte
	closeMessageSent chan struct{}
	pingTicker       *time.Ticker
	onMessage        func(id model.UserId, message []byte)
	onUpdate         func(id model.UserId)
	onClose          func(id model.UserId)
	conn             *ws.Conn
	readClosed       bool
	writeClosed      bool
	closing          bool
	closed           bool
}

func (p *user) HasGameId() bool {
	unlock := p.rlock("HasGameId")
	defer unlock()

	return p.gameId != ""
}

func (p *user) GameId() model.GameId {
	unlock := p.rlock("GameId")
	defer unlock()

	return p.gameId
}

func (p *user) SetGameId(gameId model.GameId) {
	unlock := p.lock("SetGameId")
	defer unlock()

	p.gameId = gameId
}

func (p *user) UnsetGameId() {
	unlock := p.lock("UnsetGameId")
	defer unlock()

	p.gameId = ""
}

func (p *user) PlayerId() model.PlayerId {
	unlock := p.rlock("PlayerId")
	defer unlock()

	return model.NewPlayerId(p.gameId, p.Id())
}

func (p *user) IsInactive() bool {
	unlock := p.rlock("IsInactive")
	defer unlock()

	return !p.active
}

func (p *user) IsActive() bool {
	unlock := p.rlock("IsActive")
	defer unlock()

	return p.active
}

func (p *user) IsNotPlaying() bool {
	unlock := p.rlock("IsNotPlaying")
	defer unlock()

	return p.active && p.gameId == ""

}
func (p *user) IsPlaying() bool {
	unlock := p.rlock("IsPlaying")
	defer unlock()

	return p.active && p.gameId != ""
}

func (p *user) ConnectSocket(w http.ResponseWriter, r *http.Request) error {
	logger := p.logger.With(zap.String("routine", "connect-socket"))
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info(fmt.Sprintf("[ws] user %v → connect :: ERROR %q", p.Id(), err.Error()), zap.Error(err))
		return err
	}

	logger.Info(fmt.Sprintf("[ws] user %v → open...", p.Id()))
	p.Open(conn)
	go p.WriteSocket()
	go p.ReadSocket()
	return nil
}

func (p *user) ReadSocket() {
	logger := p.logger.With(zap.String("routine", "read-socket"))

	defer func() {
		r := recover()
		if r != nil {
			if err, ok := r.(error); ok {
				logger.Info(fmt.Sprintf("[ws] user %v → read CLOSED: ERROR %q → Close", p.Id(), err.Error()), zap.Error(err))
			} else {
				logger.Info(fmt.Sprintf("[ws] user %v → read CLOSED: PANIC → Close", p.Id()), zap.Any("panic", r))
			}
		} else {
			logger.Info(fmt.Sprintf("[ws] user %v → read CLOSED → Close", p.Id()))
		}

		unlock := p.lock("ReadSocket")
		p.readClosed = true
		unlock()

		p.Close()
	}()
	logger.Info(fmt.Sprintf("[ws] user %v → read OPEN", p.Id()))

	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(msg string) error {
		if DebugPing {
			logger.Info(fmt.Sprintf("[ws] user %v ← pong", p.Id()), zap.Any("msg", msg))
		}
		p.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			logger.Warn(fmt.Sprintf("[ws] user %v ← receive ERROR %q → BREAK", p.Id(), err.Error()), zap.Error(err))
			break
		}
		if len(message) == 0 {
			logger.Info(fmt.Sprintf("[ws] user %v ← receive EMPTY message → SKIP", p.Id()))
			continue
		}
		if DebugMessage {
			logger.Info(fmt.Sprintf("[ws] user %v ← receive message ← %s", p.Id(), message))
		}
		if p.onMessage != nil {
			p.onMessage(p.Id(), message)
		}
	}
}

func (p *user) WriteSocket() {
	logger := p.logger.With(zap.String("routine", "write-socket"))

	defer func() {
		r := recover()
		if r != nil {
			if err, ok := r.(error); ok {
				logger.Info(fmt.Sprintf("[ws] user %v → write CLOSED: ERROR %q", p.Id(), err.Error()), zap.Error(err))
			} else {
				logger.Info(fmt.Sprintf("[ws] user %v → write CLOSED: PANIC", p.Id()), zap.Any("panic", r))
			}
		} else {
			logger.Info(fmt.Sprintf("[ws] user %v → write CLOSED", p.Id()))
		}

		unlock := p.lock("WriteSocket")
		p.writeClosed = true
		unlock()

		p.Close()
	}()
	logger.Info(fmt.Sprintf("[ws] user %v → write OPEN", p.Id()))

	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.conn.WriteMessage(ws.CloseMessage, []byte{})
				logger.Info(fmt.Sprintf("[ws] user %v → send channel CLOSED → CLOSE message sent → BREAK", p.Id()))
				p.closeMessageSent <- struct{}{}
				return
			}

			w, err := p.conn.NextWriter(ws.TextMessage)
			if err != nil {
				logger.Info(fmt.Sprintf("[ws] user %v → send message: ERROR %q → BREAK", p.Id(), err.Error()))
				return
			}
			if DebugMessage {
				logger.Info(fmt.Sprintf("[ws] user %v → send message → %s", p.Id(), message))
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				logger.Info(fmt.Sprintf("[ws] user %v → close writer: ERROR %q → BREAK", p.Id(), err.Error()))
				return
			}
		case <-p.pingTicker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				logger.Info(fmt.Sprintf("[ws] user %v → ping: ERROR %q → BREAK", p.Id(), err.Error()))
				return
			}
			if DebugPing {
				logger.Info(fmt.Sprintf("[ws] user %v → ping", p.Id()))
			}
		}
	}
}

func (p *user) Send(bytes []byte) {
	unlock := p.lock("Send")
	defer unlock()

	if p.active && p.send != nil {
		p.send <- bytes
	}
}

func (p *user) IsReadClosed() bool {
	unlock := p.rlock("IsReadClosed")
	defer unlock()

	return p.readClosed
}

func (p *user) IsWriteClosed() bool {
	unlock := p.rlock("IsWriteClosed")
	defer unlock()

	return p.writeClosed
}

func (p *user) IsClosing() bool {
	unlock := p.rlock("IsClosing")
	defer unlock()

	return p.closing
}

func (p *user) IsClosed() bool {
	unlock := p.rlock("IsClosed")
	defer unlock()

	return p.closed
}

func (p *user) Open(conn *ws.Conn) {
	logger := p.logger.With(zap.String("action", "open"))
	if !p.IsClosed() {
		logger.Info(fmt.Sprintf("[ws] user %v → NOT closed → Close", p.Id()))
		p.Close()
	}

	logger.Info(fmt.Sprintf("[ws] user %v → OPEN -> Prepare", p.Id()))

	unlock := p.lock("Open")
	p.conn = conn
	p.send = make(chan []byte, 256)
	p.closeMessageSent = make(chan struct{})
	p.pingTicker = time.NewTicker(pingPeriod)
	p.readClosed = false
	p.writeClosed = false
	p.closing = false
	p.closed = false
	unlock()

	logger.Info(fmt.Sprintf("[ws] user %v → OPEN -> Activate", p.Id()))
	p.Activate()
}

func (p *user) Close() {
	logger := p.logger.With(zap.String("action", "close"))
	if p.IsClosed() {
		logger.Info(fmt.Sprintf("[ws] user %v → ALREADY closed", p.Id()))
		return
	}
	if p.IsClosing() {
		logger.Info(fmt.Sprintf("[ws] user %v → ALREADY closing", p.Id()))
		return
	}

	unlock := p.lock("Closing")
	p.closing = true
	unlock()

	if p.onClose != nil {
		p.onClose(p.Id())
	}

	logger.Info(fmt.Sprintf("[ws] user %v → stop ping ticker", p.Id()))
	p.pingTicker.Stop()

	if p.IsWriteClosed() {
		logger.Info(fmt.Sprintf("[ws] user %v → stop send channel", p.Id()))
		close(p.send)
		close(p.closeMessageSent)
	} else {
		logger.Info(fmt.Sprintf("[ws] user %v → stop send channel", p.Id()))
		close(p.send)
		logger.Info(fmt.Sprintf("[ws] user %v → stop send channel → waiting close message to be sent...", p.Id()))
		<-p.closeMessageSent
		close(p.closeMessageSent)
	}

	logger.Info(fmt.Sprintf("[ws] user %v → closing connection", p.Id()))
	p.conn.Close()

	unlock = p.lock("Close")
	p.pingTicker = nil
	p.send = nil
	p.closeMessageSent = nil
	p.conn = nil
	p.closing = false
	p.closed = true
	unlock()

	logger.Info(fmt.Sprintf("[ws] user %v → CLOSED → DEACTIVATE", p.Id()))
	p.Deactivate()
}

func (p *user) Activate() {
	logger := p.logger.With(zap.String("action", "activate"))
	if p.IsActive() {
		return
	}

	unlock := p.lock("Activate")
	p.active = true
	unlock()

	if p.onUpdate != nil {
		logger.Info(fmt.Sprintf("[ws] user %v → ACTIVE → callback", p.Id()))
		p.onUpdate(p.Id())
	} else {
		logger.Info(fmt.Sprintf("[ws] user %v → ACTIVE", p.Id()))
	}
}

func (p *user) Deactivate() {
	logger := p.logger.With(zap.String("action", "deactivate"))
	if !p.IsActive() {
		return
	}

	unlock := p.lock("Deactivate")
	p.active = false
	unlock()

	if p.onUpdate != nil {
		logger.Info(fmt.Sprintf("[ws] user %v → INACTIVE → callback", p.Id()))
		p.onUpdate(p.Id())
	} else {
		logger.Info(fmt.Sprintf("[ws] user %v → INACTIVE", p.Id()))
	}
}

// //////////////////////////////////////////////////
// lock

func (p *user) rlock(requester string) func() {
	logger := p.logger.WithOptions(zap.AddCallerSkip(1))
	if DebugLock {
		logger.Info(fmt.Sprintf(" >>> R-LOCK >>> user %s >>> %s ", p.Id(), requester))
	}
	p.RLock()
	return func() {
		p.RUnlock()
		if DebugLock {
			logger.Info(fmt.Sprintf(" <<< R-LOCK <<< user %s <<< %s ", p.Id(), requester))
		}
	}
}

func (p *user) lock(requester string) func() {
	logger := p.logger.WithOptions(zap.AddCallerSkip(1))
	if DebugLock {
		logger.Info(fmt.Sprintf(" >>> W-LOCK >>> user %s >>> %s ", p.Id(), requester))
	}
	p.Lock()
	return func() {
		p.Unlock()
		if DebugLock {
			logger.Info(fmt.Sprintf(" <<< W-LOCK <<< user %s <<< %s ", p.Id(), requester))
		}
	}
}

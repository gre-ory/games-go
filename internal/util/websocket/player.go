package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Player[IdT comparable, GameIdT comparable] interface {
	HasId() bool
	Id() IdT

	HasGameId() bool
	GameId() GameIdT
	SetGameId(gameId GameIdT)
	UnsetGameId()
	CanJoin() bool

	ConnectSocket(w http.ResponseWriter, r *http.Request) error
	// ReadSocket()
	// WriteSocket()

	IsActive() bool
	Activate(logger *zap.Logger)
	Deactivate(logger *zap.Logger)

	Send(bytes []byte)
	Close(logger *zap.Logger)
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

func NewPlayer[IdT comparable, GameIdT comparable](
	logger *zap.Logger,
	id IdT,
	onMessage func(id IdT, message []byte),
	onUpdate func(id IdT),
	onClose func(id IdT),
) Player[IdT, GameIdT] {
	return &player[IdT, GameIdT]{
		id:          id,
		logger:      logger,
		onMessage:   onMessage,
		onUpdate:    onUpdate,
		onClose:     onClose,
		readClosed:  true,
		writeClosed: true,
		closing:     false,
		closed:      true,
	}
}

type player[IdT comparable, GameIdT comparable] struct {
	sync.RWMutex
	id               IdT
	gameId           GameIdT
	logger           *zap.Logger
	active           bool
	send             chan []byte
	closeMessageSent chan struct{}
	pingTicker       *time.Ticker
	onMessage        func(id IdT, message []byte)
	onUpdate         func(id IdT)
	onClose          func(id IdT)
	conn             *ws.Conn
	readClosed       bool
	writeClosed      bool
	closing          bool
	closed           bool
}

func (p *player[IdT, GameIdT]) HasId() bool {
	var empty IdT
	return p.id != empty
}

func (p *player[IdT, GameIdT]) Id() IdT {
	return p.id
}

func (p *player[IdT, GameIdT]) HasGameId() bool {
	var empty GameIdT
	return p.gameId != empty
}

func (p *player[IdT, GameIdT]) GameId() GameIdT {
	return p.gameId
}

func (p *player[IdT, GameIdT]) SetGameId(gameId GameIdT) {
	p.gameId = gameId
}

func (p *player[IdT, GameIdT]) UnsetGameId() {
	var empty GameIdT
	p.gameId = empty
}

func (p *player[IdT, GameIdT]) CanJoin() bool {
	p.RLock()
	defer p.RUnlock()
	return p.active && p.HasId() && !p.HasGameId()
}

func (p *player[IdT, GameIdT]) ConnectSocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	logger := p.logger.With(zap.Any("player", p.Id()))

	p.Open(logger, conn)
	go p.WriteSocket(logger)
	go p.ReadSocket(logger)
	return nil
}

func (p *player[IdT, GameIdT]) ReadSocket(logger *zap.Logger) {
	logger = logger.With(zap.String("thread", "read-socket"))

	defer func() {
		r := recover()
		if r != nil {
			if err, ok := r.(error); ok {
				logger.Info(fmt.Sprintf("[ws] player %v → read CLOSED: ERROR %q → Close", p.Id(), err.Error()), zap.Error(err))
			} else {
				logger.Info(fmt.Sprintf("[ws] player %v → read CLOSED: PANIC → Close", p.Id()), zap.Any("panic", r))
			}
		} else {
			logger.Info(fmt.Sprintf("[ws] player %v → read CLOSED → Close", p.Id()))
		}
		p.Lock()
		p.readClosed = true
		p.Unlock()
		p.Close(logger)
	}()
	logger.Info(fmt.Sprintf("[ws] player %v → read OPEN", p.Id()))

	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(msg string) error {
		if DebugPing {
			logger.Info(fmt.Sprintf("[ws] player %v ← pong", p.Id()), zap.Any("msg", msg))
		}
		p.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			logger.Warn(fmt.Sprintf("[ws] player %v ← receive ERROR %q → BREAK", p.Id(), err.Error()), zap.Error(err))
			break
		}
		if len(message) == 0 {
			logger.Info(fmt.Sprintf("[ws] player %v ← receive EMPTY message → SKIP", p.Id()))
			continue
		}
		if DebugMessage {
			logger.Info(fmt.Sprintf("[ws] player %v ← receive message ← %s", p.Id(), message))
		}
		if p.onMessage != nil {
			p.onMessage(p.id, message)
		}
	}
}

func (p *player[IdT, GameIdT]) WriteSocket(logger *zap.Logger) {
	logger = logger.With(zap.String("thread", "write-socket"))

	defer func() {
		r := recover()
		if r != nil {
			if err, ok := r.(error); ok {
				logger.Info(fmt.Sprintf("[ws] player %v → write CLOSED: ERROR %q", p.Id(), err.Error()), zap.Error(err))
			} else {
				logger.Info(fmt.Sprintf("[ws] player %v → write CLOSED: PANIC", p.Id()), zap.Any("panic", r))
			}
		} else {
			logger.Info(fmt.Sprintf("[ws] player %v → write CLOSED", p.Id()))
		}
		p.Lock()
		p.writeClosed = true
		p.Unlock()
		p.Close(logger)
	}()
	logger.Info(fmt.Sprintf("[ws] player %v → write OPEN", p.Id()))

	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				logger.Info(fmt.Sprintf("[ws] player %v → send channel CLOSED → CLOSE message sent → BREAK", p.Id()))
				p.closeMessageSent <- struct{}{}
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Info(fmt.Sprintf("[ws] player %v → send message: ERROR %q → BREAK", p.Id(), err.Error()))
				return
			}
			if DebugMessage {
				logger.Info(fmt.Sprintf("[ws] player %v → send message → %s", p.Id(), message))
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				logger.Info(fmt.Sprintf("[ws] player %v → close writer: ERROR %q → BREAK", p.Id(), err.Error()))
				return
			}
		case <-p.pingTicker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Info(fmt.Sprintf("[ws] player %v → ping: ERROR %q → BREAK", p.Id(), err.Error()))
				return
			}
			if DebugPing {
				logger.Info(fmt.Sprintf("[ws] player %v → ping", p.Id()))
			}
		}
	}
}

func (p *player[IdT, GameIdT]) Send(bytes []byte) {
	p.Lock()
	defer p.Unlock()

	if p.active && p.send != nil {
		p.send <- bytes
	}
}

func (p *player[IdT, GameIdT]) IsReadClosed() bool {
	p.RLock()
	defer p.RUnlock()
	return p.readClosed
}

func (p *player[IdT, GameIdT]) IsWriteClosed() bool {
	p.RLock()
	defer p.RUnlock()
	return p.writeClosed
}

func (p *player[IdT, GameIdT]) IsClosing() bool {
	p.RLock()
	defer p.RUnlock()
	return p.closing
}

func (p *player[IdT, GameIdT]) IsClosed() bool {
	p.RLock()
	defer p.RUnlock()
	return p.closed
}

func (p *player[IdT, GameIdT]) Open(logger *zap.Logger, conn *ws.Conn) {
	if !p.IsClosed() {
		logger.Info(fmt.Sprintf("[ws] player %v → NOT closed → Close", p.Id()))
		p.Close(logger)
	}

	p.Lock()
	p.conn = conn
	p.send = make(chan []byte, 256)
	p.closeMessageSent = make(chan struct{})
	p.pingTicker = time.NewTicker(pingPeriod)
	p.readClosed = false
	p.writeClosed = false
	p.closing = false
	p.closed = false
	p.Unlock()

	logger.Info(fmt.Sprintf("[ws] player %v → OPEN -> ACTIVATE", p.Id()))
	p.Activate(logger)
}

func (p *player[IdT, GameIdT]) Close(logger *zap.Logger) {
	if p.IsClosed() {
		logger.Info(fmt.Sprintf("[ws] player %v → ALREADY closed", p.Id()))
		return
	}
	if p.IsClosing() {
		logger.Info(fmt.Sprintf("[ws] player %v → ALREADY closing", p.Id()))
		return
	}

	p.Lock()
	p.closing = true
	p.Unlock()

	if p.onClose != nil {
		p.onClose(p.id)
	}

	logger.Info(fmt.Sprintf("[ws] player %v → stop ping ticker", p.Id()))
	p.pingTicker.Stop()

	if p.IsWriteClosed() {
		logger.Info(fmt.Sprintf("[ws] player %v → stop send channel", p.Id()))
		close(p.send)
		close(p.closeMessageSent)
	} else {
		logger.Info(fmt.Sprintf("[ws] player %v → stop send channel", p.Id()))
		close(p.send)
		logger.Info(fmt.Sprintf("[ws] player %v → stop send channel → waiting close message to be sent...", p.Id()))
		<-p.closeMessageSent
		close(p.closeMessageSent)
	}

	logger.Info(fmt.Sprintf("[ws] player %v → closing connection", p.Id()))
	p.conn.Close()

	p.Lock()
	p.pingTicker = nil
	p.send = nil
	p.closeMessageSent = nil
	p.conn = nil
	p.closing = false
	p.closed = true
	p.Unlock()

	logger.Info(fmt.Sprintf("[ws] player %v → CLOSED → DEACTIVATE", p.Id()))
	p.Deactivate(logger)
}

func (p *player[IdT, GameIdT]) IsActive() bool {
	p.RLock()
	defer p.RUnlock()

	return p.active
}

func (p *player[IdT, GameIdT]) Activate(logger *zap.Logger) {
	if p.IsActive() {
		return
	}

	p.Lock()
	p.active = true
	p.Unlock()

	if p.onUpdate != nil {
		logger.Info(fmt.Sprintf("[ws] player %v → ACTIVE → callback", p.Id()))
		p.onUpdate(p.id)
	} else {
		logger.Info(fmt.Sprintf("[ws] player %v → ACTIVE", p.Id()))
	}
}

func (p *player[IdT, GameIdT]) Deactivate(logger *zap.Logger) {
	if !p.IsActive() {
		return
	}

	p.Lock()
	p.active = false
	p.Unlock()

	if p.onUpdate != nil {
		logger.Info(fmt.Sprintf("[ws] player %v → INACTIVE → callback", p.Id()))
		p.onUpdate(p.id)
	} else {
		logger.Info(fmt.Sprintf("[ws] player %v → INACTIVE", p.Id()))
	}
}

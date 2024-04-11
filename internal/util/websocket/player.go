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
	GetId() IdT

	HasGameId() bool
	GetGameId() GameIdT
	SetGameId(gameId GameIdT)
	UnsetGameId()
	CanJoin() bool

	ConnectSocket(w http.ResponseWriter, r *http.Request) error
	// ReadSocket()
	// WriteSocket()

	Active() bool
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

func NewPlayer[IdT comparable, GameIdT comparable](
	logger *zap.Logger,
	id IdT,
	onMessage func(id IdT, message []byte),
	onUpdate func(id IdT),
	onClose func(id IdT),
) Player[IdT, GameIdT] {
	return &player[IdT, GameIdT]{
		Id:        id,
		logger:    logger,
		onMessage: onMessage,
		onUpdate:  onUpdate,
		onClose:   onClose,
	}
}

type player[IdT comparable, GameIdT comparable] struct {
	sync.Mutex
	Id         IdT
	gameId     GameIdT
	logger     *zap.Logger
	active     bool
	send       chan []byte
	pingTicker *time.Ticker
	onMessage  func(id IdT, message []byte)
	onUpdate   func(id IdT)
	onClose    func(id IdT)
	conn       *ws.Conn
}

func (p *player[IdT, GameIdT]) HasId() bool {
	var empty IdT
	return p.Id != empty
}

func (p *player[IdT, GameIdT]) GetId() IdT {
	return p.Id
}

func (p *player[IdT, GameIdT]) HasGameId() bool {
	var empty GameIdT
	return p.gameId != empty
}

func (p *player[IdT, GameIdT]) GetGameId() GameIdT {
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
	return p.active && p.HasId() && !p.HasGameId()
}

func (p *player[IdT, GameIdT]) ConnectSocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	p.Open(conn)
	go p.WriteSocket()
	go p.ReadSocket()
	p.Activate()
	return nil
}

func (p *player[IdT, GameIdT]) ReadSocket() {
	p.logger.Info(fmt.Sprintf("[ws] player %v → read OPEN", p.Id))
	p.Activate()
	defer func() {
		p.logger.Info(fmt.Sprintf("[ws] player %v → read CLOSED", p.Id))
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				p.logger.Info(fmt.Sprintf("[ws] player %v → PANIC → ERROR", p.Id), zap.Error(err))
			} else {
				p.logger.Info(fmt.Sprintf("[ws] player %v → PANIC", p.Id), zap.Any("panic", r))
			}
		}
		p.Deactivate()
		p.Close()
	}()
	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(msg string) error {
		// p.logger.Info(fmt.Sprintf("[ws] player %v ← PONG", p.id), zap.Any("msg", msg))
		p.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			p.logger.Error(fmt.Sprintf("[ws] error: %s", err.Error()), zap.Error(err))
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}
		p.logger.Info(fmt.Sprintf("[ws] player %v ← %s", p.Id, message))
		if p.onMessage != nil {
			p.onMessage(p.Id, message)
		}
	}
}

func (p *player[IdT, GameIdT]) WriteSocket() {
	p.logger.Info(fmt.Sprintf("[ws] player %v → write OPEN", p.Id))
	defer func() {
		p.logger.Info(fmt.Sprintf("[ws] player %v → write CLOSED", p.Id))
		if r := recover(); r != nil {
			p.logger.Info(fmt.Sprintf("[ws] player %v → PANIC", p.Id), zap.Any("panic", r))
		}
		p.Deactivate()
		p.Close()
	}()
	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				p.logger.Info(fmt.Sprintf("[ws] player %v → CLOSE message", p.Id))
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			p.logger.Info(fmt.Sprintf("[ws] player %v → %s", p.Id, message))
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-p.pingTicker.C:
			// p.logger.Info(fmt.Sprintf("[ws] player %v → PING", p.id))
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
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

func (p *player[IdT, GameIdT]) Open(conn *ws.Conn) {
	p.Lock()
	defer p.Unlock()

	if p.conn != nil {
		p.close()
	}

	p.conn = conn
	p.send = make(chan []byte, 256)
	p.pingTicker = time.NewTicker(pingPeriod)
	p.logger.Info(fmt.Sprintf("[ws] player %v → OPEN", p.Id))
}

func (p *player[IdT, GameIdT]) Active() bool {
	p.Lock()
	defer p.Unlock()

	return p.active
}

func (p *player[IdT, GameIdT]) Activate() {
	p.Lock()
	defer p.Unlock()

	if !p.active {
		p.active = true
		if p.onUpdate != nil {
			p.onUpdate(p.Id)
		}
		p.logger.Info(fmt.Sprintf("[ws] player %v → ACTIVE", p.Id))
	}
}

func (p *player[IdT, GameIdT]) Close() {
	p.Lock()
	defer p.Unlock()

	if p.conn != nil {
		p.close()
	}
}

func (p *player[IdT, GameIdT]) close() {
	if p.onClose != nil {
		p.onClose(p.Id)
	}

	p.pingTicker.Stop()
	p.pingTicker = nil

	p.logger.Info(fmt.Sprintf("[ws] player %v → stop send channel", p.Id))
	close(p.send)
	p.send = nil

	p.logger.Info(fmt.Sprintf("[ws] player %v → close connection", p.Id))
	p.conn.Close()
	p.conn = nil

	p.logger.Info(fmt.Sprintf("[ws] player %v → CLOSED", p.Id))
}

func (p *player[IdT, GameIdT]) Deactivate() {
	p.Lock()
	defer p.Unlock()

	if p.active {
		p.active = false
		if p.onUpdate != nil {
			p.onUpdate(p.Id)
		}
		p.logger.Info(fmt.Sprintf("[ws] player %v → INACTIVE", p.Id))
	}
}

package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Player[IdT comparable, GameIdT comparable] interface {
	GetId() IdT
	GetGameId() GameIdT
	SetGameId(gameId GameIdT)
	UnsetGameId()
	ConnectSocket(w http.ResponseWriter, r *http.Request) error
	ReadSocket()
	WriteSocket()
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

func NewPlayer[IdT comparable, GameIdT comparable](logger *zap.Logger, id IdT, onMessage func(id IdT, message []byte), onClose func(id IdT)) Player[IdT, GameIdT] {
	return &player[IdT, GameIdT]{
		id:        id,
		logger:    logger,
		active:    true,
		send:      make(chan []byte, 256),
		onMessage: onMessage,
		onClose:   onClose,
	}
}

type player[IdT comparable, GameIdT comparable] struct {
	id        IdT
	gameId    GameIdT
	logger    *zap.Logger
	active    bool
	send      chan []byte
	conn      *ws.Conn
	onMessage func(id IdT, message []byte)
	onClose   func(id IdT)
}

func (p *player[IdT, GameIdT]) GetId() IdT {
	return p.id
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

func (p *player[IdT, GameIdT]) ConnectSocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	p.conn = conn
	go p.WriteSocket()
	go p.ReadSocket()
	return nil
}

func (p *player[IdT, GameIdT]) ReadSocket() {
	p.logger.Info(fmt.Sprintf("[ws] player %v :: read :: OPEN", p.id))
	p.Activate()
	defer func() {
		p.conn.Close()
		// unregister
		// p.onClose(p.id)
		p.Deactivate()
		p.logger.Info(fmt.Sprintf("[ws] player %v :: read :: CLOSED", p.id))
	}()
	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			p.logger.Error(fmt.Sprintf("[ws] error: %s", err.Error()), zap.Error(err))
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}
		p.logger.Info(fmt.Sprintf("[ws] player %v ← %s", p.id, message))
		if p.onMessage != nil {
			p.onMessage(p.id, message)
		}
	}
}

func (p *player[IdT, GameIdT]) WriteSocket() {
	p.logger.Info(fmt.Sprintf("[ws] player %v :: write :: OPEN", p.id))
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.conn.Close()
		p.logger.Info(fmt.Sprintf("[ws] player %v :: write :: CLOSED", p.id))
	}()
	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			p.logger.Info(fmt.Sprintf("[ws] player %v → %s", p.id, message))
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (p *player[IdT, GameIdT]) Activate() {
	p.active = true
}

func (p *player[IdT, GameIdT]) Deactivate() {
	p.active = false
}

func (p *player[IdT, GameIdT]) Send(bytes []byte) {
	p.send <- bytes
}

func (p *player[IdT, GameIdT]) Close() {
	close(p.send)
}

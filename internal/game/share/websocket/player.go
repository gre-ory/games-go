package websocket

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/gre-ory/games-go/internal/util/loc"

	"github.com/gre-ory/games-go/internal/game/share/model"
)

// //////////////////////////////////////////////////
// websocket player

type Player interface {
	model.Player

	CanJoin() bool

	ConnectSocket(w http.ResponseWriter, r *http.Request) error
	// ReadSocket()
	// WriteSocket()

	IsActive() bool
	Activate()
	Deactivate()

	Send(bytes []byte)
	Close()

	LabelSlice() []string
	Labels() string
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

func NewPlayer(
	logger *zap.Logger,
	id model.PlayerId,
	onMessage func(id model.PlayerId, message []byte),
	onUpdate func(id model.PlayerId),
	onClose func(id model.PlayerId),
) Player {
	return &player{
		Player:      model.NewPlayer(id),
		logger:      logger.With(zap.String("player", string(id))),
		onMessage:   onMessage,
		onUpdate:    onUpdate,
		onClose:     onClose,
		readClosed:  true,
		writeClosed: true,
		closing:     false,
		closed:      true,
	}
}

func NewPlayerFromCookie(
	logger *zap.Logger,
	cookie *model.Cookie,
	onMessage func(id model.PlayerId, message []byte),
	onUpdate func(id model.PlayerId),
	onClose func(id model.PlayerId),
) Player {
	return &player{
		Player:      model.NewPlayerFromCookie(cookie),
		logger:      logger.With(zap.String("player", string(cookie.Id))),
		onMessage:   onMessage,
		onUpdate:    onUpdate,
		onClose:     onClose,
		readClosed:  true,
		writeClosed: true,
		closing:     false,
		closed:      true,
	}
}

type player struct {
	model.Player
	sync.RWMutex
	logger           *zap.Logger
	active           bool
	send             chan []byte
	closeMessageSent chan struct{}
	pingTicker       *time.Ticker
	onMessage        func(id model.PlayerId, message []byte)
	onUpdate         func(id model.PlayerId)
	onClose          func(id model.PlayerId)
	conn             *ws.Conn
	readClosed       bool
	writeClosed      bool
	closing          bool
	closed           bool
}

func (p *player) CanJoin() bool {
	p.RLock()
	defer p.RUnlock()
	return p.active && p.HasId() && !p.HasGameId()
}

func (p *player) ConnectSocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	p.Open(conn)
	go p.WriteSocket()
	go p.ReadSocket()
	return nil
}

func (p *player) ReadSocket() {
	logger := p.logger.With(zap.String("routine", "read-socket"))

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
		p.Close()
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
			p.onMessage(p.Id(), message)
		}
	}
}

func (p *player) WriteSocket() {
	logger := p.logger.With(zap.String("routine", "write-socket"))

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
		p.Close()
	}()
	logger.Info(fmt.Sprintf("[ws] player %v → write OPEN", p.Id()))

	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.conn.WriteMessage(ws.CloseMessage, []byte{})
				logger.Info(fmt.Sprintf("[ws] player %v → send channel CLOSED → CLOSE message sent → BREAK", p.Id()))
				p.closeMessageSent <- struct{}{}
				return
			}

			w, err := p.conn.NextWriter(ws.TextMessage)
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
			if err := p.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				logger.Info(fmt.Sprintf("[ws] player %v → ping: ERROR %q → BREAK", p.Id(), err.Error()))
				return
			}
			if DebugPing {
				logger.Info(fmt.Sprintf("[ws] player %v → ping", p.Id()))
			}
		}
	}
}

func (p *player) Send(bytes []byte) {
	p.Lock()
	defer p.Unlock()

	if p.active && p.send != nil {
		p.send <- bytes
	}
}

func (p *player) IsReadClosed() bool {
	p.RLock()
	defer p.RUnlock()
	return p.readClosed
}

func (p *player) IsWriteClosed() bool {
	p.RLock()
	defer p.RUnlock()
	return p.writeClosed
}

func (p *player) IsClosing() bool {
	p.RLock()
	defer p.RUnlock()
	return p.closing
}

func (p *player) IsClosed() bool {
	p.RLock()
	defer p.RUnlock()
	return p.closed
}

func (p *player) Open(conn *ws.Conn) {
	logger := p.logger.With(zap.String("action", "open"))
	if !p.IsClosed() {
		logger.Info(fmt.Sprintf("[ws] player %v → NOT closed → Close", p.Id()))
		p.Close()
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
	p.Activate()
}

func (p *player) Close() {
	logger := p.logger.With(zap.String("action", "close"))
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
		p.onClose(p.Id())
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
	p.Deactivate()
}

func (p *player) IsActive() bool {
	p.RLock()
	defer p.RUnlock()

	return p.active
}

func (p *player) Activate() {
	logger := p.logger.With(zap.String("action", "activate"))
	if p.IsActive() {
		return
	}

	p.Lock()
	p.active = true
	p.Unlock()

	if p.onUpdate != nil {
		logger.Info(fmt.Sprintf("[ws] player %v → ACTIVE → callback", p.Id()))
		p.onUpdate(p.Id())
	} else {
		logger.Info(fmt.Sprintf("[ws] player %v → ACTIVE", p.Id()))
	}
}

func (p *player) Deactivate() {
	logger := p.logger.With(zap.String("action", "deactivate"))
	if !p.IsActive() {
		return
	}

	p.Lock()
	p.active = false
	p.Unlock()

	if p.onUpdate != nil {
		logger.Info(fmt.Sprintf("[ws] player %v → INACTIVE → callback", p.Id()))
		p.onUpdate(p.Id())
	} else {
		logger.Info(fmt.Sprintf("[ws] player %v → INACTIVE", p.Id()))
	}
}

func (p *player) YourMessage(localizer loc.Localizer) template.HTML {
	if p.IsActive() {
		return p.Player.YourMessage(localizer)
	} else {
		return localizer.Loc("YouDisconnected")
	}
}

func (p *player) Message(localizer loc.Localizer) template.HTML {
	if p.IsActive() {
		return p.Player.Message(localizer)
	} else {
		return localizer.Loc("PlayerDisconnected")
	}
}

func (p *player) LabelSlice() []string {
	labels := make([]string, 0)
	labels = append(labels, "player")
	if p.IsActive() {
		labels = append(labels, p.Status().LabelSlice()...)
	} else {
		labels = append(labels, "disconnected")
	}
	return labels
}

func (p *player) Labels() string {
	return strings.Join(p.LabelSlice(), " ")
}

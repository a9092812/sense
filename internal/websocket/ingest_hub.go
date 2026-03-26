package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"github.com/Kartik30R/sense/internal/models"
	"github.com/Kartik30R/sense/internal/redis"
	services "github.com/Kartik30R/sense/internal/service"
	"github.com/Kartik30R/sense/pkg/logger"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Command struct {
	Type     string `json:"type"`
	DeviceID string `json:"deviceId"`
}

type IngestHub struct {
	clients map[*IngestClient]bool

	register   chan *IngestClient
	unregister chan *IngestClient

	mu sync.RWMutex

	activeStreams map[string]string

	redisPub      *redis.Publisher
	packetService *services.PacketService
	mobileService *services.MobileDeviceService
}

func NewIngestHub(redisPub *redis.Publisher, packetService *services.PacketService, mobileService *services.MobileDeviceService) *IngestHub {
	return &IngestHub{
		clients:       make(map[*IngestClient]bool),
		register:      make(chan *IngestClient),
		unregister:    make(chan *IngestClient),
		activeStreams: make(map[string]string),
		redisPub:      redisPub,
		packetService: packetService,
		mobileService: mobileService,
	}
}

func (h *IngestHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			// 📱 Register mobile in DB so dashboard can discover it
			if h.mobileService != nil {
				if err := h.mobileService.EnsureMobile(client.clientID); err != nil {
					logger.Error("Failed to register mobile gateway", zap.Error(err))
				}
			}

			logger.Info("Mobile gateway connected",
				zap.String("clientId", client.clientID),
			)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				// remove ownership if this client was streaming
				for deviceID, ownerID := range h.activeStreams {
					if ownerID == client.clientID {
						delete(h.activeStreams, deviceID)
					}
				}
			}
			h.mu.Unlock()

			logger.Info("Mobile gateway disconnected",
				zap.String("clientId", client.clientID),
			)
		}
	}
}

func (h *IngestHub) SendCommand(mobileID string, deviceID string, cmd Command) error {

	h.mu.Lock()
	defer h.mu.Unlock()

	data, _ := json.Marshal(cmd)

	// START
	if cmd.Type == "start_stream" {

		if _, exists := h.activeStreams[deviceID]; exists {
			return nil
		}

		for client := range h.clients {
			if client.clientID == mobileID {
				h.activeStreams[deviceID] = client.clientID

				client.conn.WriteMessage(websocket.TextMessage, data)

				logger.Info("Assigned stream",
					zap.String("deviceId", deviceID),
					zap.String("clientId", client.clientID),
				)

				return nil
			}
		}
	}

	// STOP
	if cmd.Type == "stop_stream" {

		if ownerID, ok := h.activeStreams[deviceID]; ok {

			for client := range h.clients {
				if client.clientID == ownerID {
					client.conn.WriteMessage(websocket.TextMessage, data)
					break
				}
			}

			delete(h.activeStreams, deviceID)
		}
	}

	return nil
}

type IngestClient struct {
	hub      *IngestHub
	conn     *websocket.Conn
	deviceID string
	clientID string
	ctx      context.Context
}

func (c *IngestClient) readPump() {

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {

		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var packet models.SensorPacketMessage

		if err := json.Unmarshal(raw, &packet); err != nil {

			logger.Warn("Invalid packet", zap.Error(err))
			continue
		}

		if packet.ParsedType == "" {
			packet.ParsedType = (&packet).String("type")
		}

		if c.deviceID == "" && packet.DeviceID != "" {
			c.deviceID = packet.DeviceID
			logger.Info("Associated gateway with device", zap.String("deviceId", c.deviceID))
		}

		// 📍 Discovery: Update association and last_seen for THIS mobile
		// We do this BEFORE the gate so that all mobiles in range get credit for seeing the sensor.
		_ = c.hub.packetService.RegisterDiscovery(
			c.clientID,
			packet.DeviceID,
			packet.DeviceAddress,
			packet.ParsedType,
		)

		// 🔥 GATE: allow only assigned mobile to send this sensor data (for live broadcast)
		c.hub.mu.RLock()
		ownerID := c.hub.activeStreams[packet.DeviceID]
		c.hub.mu.RUnlock()

		if ownerID != c.clientID {
			// ❌ Not the selected mobile for live streaming/broadcast
			c.writeACK(true, "discovery_ok")
			continue
		}

		// 1️⃣ Fast Path for Dashboard (Redis)
		if err := c.hub.redisPub.Publish(redis.SensorDataChannel, packet); err != nil {
			logger.Warn("Redis publish failed", zap.Error(err))
		}

		logger.Debug("WS packet received (live only)",
			zap.String("deviceId", packet.DeviceID),
			zap.String("type", packet.ParsedType),
		)

		c.writeACK(true, "ok")
	}
}

func (c *IngestClient) writePing() {

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {

		case <-c.ctx.Done():
			return

		case <-ticker.C:

			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *IngestClient) writeACK(success bool, msg string) {

	type ack struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	data, _ := json.Marshal(ack{
		Success: success,
		Message: msg,
	})

	c.conn.WriteMessage(websocket.TextMessage, data)
}

func ServeIngest(hub *IngestHub, ctx context.Context, w http.ResponseWriter, r *http.Request) {

	clientID := r.URL.Query().Get("clientId")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &IngestClient{
		hub:      hub,
		conn:     conn,
		clientID: clientID,
		ctx:      ctx,
	}

	hub.register <- client

	go client.writePing()
	go client.readPump()
}
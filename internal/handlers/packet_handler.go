package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Kartik30R/sense/internal/dto"
	services "github.com/Kartik30R/sense/internal/service"
	"github.com/Kartik30R/sense/pkg/logger"
)

type PacketHandler struct {
	packetService *services.PacketService
}

func NewPacketHandler(packetService *services.PacketService) *PacketHandler {

	return &PacketHandler{
		packetService: packetService,
	}
}

func (h *PacketHandler) UploadPacket(c *gin.Context) {

	var req dto.PacketUploadRequest

	// Bind + validate JSON
	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "invalid payload: " + err.Error(),
			Code:    http.StatusBadRequest,
		})

		return
	}

	// Send packet to Kafka via service
	ctx := c.Request.Context()

	err := h.packetService.Ingest(ctx, req)

	if err != nil {

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "failed to enqueue packet",
			Code:    http.StatusInternalServerError,
		})

		return
	}

	// Return 202 Accepted
	c.JSON(http.StatusAccepted, dto.PacketUploadResponse{
		Success: true,
		Message: "packet buffered for processing",
	})
}


func (h *PacketHandler) UploadBatch(c *gin.Context) {
	var reqs []dto.PacketUploadRequest

	if err := c.ShouldBindJSON(&reqs); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "invalid batch payload: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if len(reqs) == 0 {
		c.JSON(http.StatusOK, dto.PacketUploadResponse{
			Success: true,
			Message: "empty batch, nothing to do",
		})
		return
	}

	// 📦 Batch size log
	logger.Info("Batch received",
		zap.Int("size", len(reqs)),
	)

	// 📊 Device distribution
	for _, p := range reqs {
		logger.Debug("Packet in batch",
			zap.String("deviceId", p.DeviceID),
			zap.String("type", p.ParsedType),
		)
	}

	ctx := c.Request.Context()
	err := h.packetService.IngestBatch(ctx, reqs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "failed to enqueue batch",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusAccepted, dto.PacketUploadResponse{
		Success: true,
		Message: "batch buffered for processing",
	})
}
package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Kartik30R/sense/internal/dto"
	"github.com/Kartik30R/sense/internal/models"
	services "github.com/Kartik30R/sense/internal/service"
)

type HistoryHandler struct {
	historyService *services.HistoryService
}

func NewHistoryHandler(hs *services.HistoryService) *HistoryHandler {
	return &HistoryHandler{
		historyService: hs,
	}
}

func (h *HistoryHandler) GetHistory(c *gin.Context) {
	deviceID := strings.TrimSpace(c.Param("id"))
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "device id is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	sensorStr := strings.TrimSpace(c.Query("sensor"))
	if sensorStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "sensor type query parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	sensor := models.SensorType(sensorStr)

	var startTime time.Time
	var endTime time.Time

	startStr := c.Query("start_time")
	if startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err == nil {
			startTime = t
		}
	}

	endStr := c.Query("end_time")
	if endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err == nil {
			endTime = t
		}
	}

	limit := 10000
	limitStr := c.Query("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = l
		}
	}

	data, err := h.historyService.GetDeviceHistory(deviceID, sensor, startTime, endTime, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Kartik30R/sense/internal/repository"
	services "github.com/Kartik30R/sense/internal/service"
)

type DeviceHandler struct {
	service *services.DeviceService
}

func NewDeviceHandler(s *services.DeviceService) *DeviceHandler {
	return &DeviceHandler{service: s}
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
	mobileID := c.Query("mobileId")
	var devices []repository.DeviceRow
	var err error

	if mobileID != "" {
		devices, err = h.service.ListSensorsByMobile(mobileID)
	} else {
		devices, err = h.service.ListDevices()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(devices),
		"data":    devices,
	})
}
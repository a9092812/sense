package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Kartik30R/sense/internal/dto"
	"github.com/Kartik30R/sense/internal/models"
	services "github.com/Kartik30R/sense/internal/service"
)

type MobileDeviceHandler struct {
	mobileService *services.MobileDeviceService
	deviceService *services.DeviceService
}

func NewMobileDeviceHandler(ms *services.MobileDeviceService, ds *services.DeviceService) *MobileDeviceHandler {
	return &MobileDeviceHandler{
		mobileService: ms,
		deviceService: ds,
	}
}

func (h *MobileDeviceHandler) GetMobiles(c *gin.Context) {
	mobiles, err := h.mobileService.ListMobiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to fetch mobile devices",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Calculate online/offline status (e.g. within last 30 seconds)
	threshold := time.Now().Add(-30 * time.Second)
	var response []models.MobileStatus

	for _, m := range mobiles {
		response = append(response, models.MobileStatus{
			MobileDevice: m,
			IsOnline:     m.LastSeen.After(threshold),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *MobileDeviceHandler) GetSensorsForMobile(c *gin.Context) {
	mobileID := c.Param("mobileId")

	sensors, err := h.deviceService.ListSensorsByMobile(mobileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to fetch sensors for mobile",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// The repository already calculates IsLive based on a 15s threshold.
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sensors,
	})
}

package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardCommandRequest struct {
	MobileID string `json:"mobileId"`
	DeviceID string `json:"deviceId"`
	Type     string `json:"type"`
}

func DashboardCommandHandler(hub *IngestHub) gin.HandlerFunc {

	return func(c *gin.Context) {

		var req DashboardCommandRequest

		if err := c.ShouldBindJSON(&req); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid command",
			})

			return
		}

		// Normalize command type for Android
		translatedType := req.Type
		switch req.Type {
		case "START_LIVE":
			translatedType = "start_stream"
		case "STOP_LIVE":
			translatedType = "stop_stream"
		}

		cmd := Command{
			Type:     translatedType,
			DeviceID: req.DeviceID,
		}

		hub.SendCommand(req.MobileID, req.DeviceID, cmd)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}
package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Kartik30R/sense/internal/handlers"
	"github.com/Kartik30R/sense/internal/websocket"
)

func SetupRoutes(
	router *gin.Engine,
	packetHandler *handlers.PacketHandler,
	historyHandler *handlers.HistoryHandler,
	deviceHandler *handlers.DeviceHandler,
	mobileDeviceHandler *handlers.MobileDeviceHandler,
	hub *websocket.Hub,
	ingestHub *websocket.IngestHub,
){

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	api := router.Group("/api")
{
	api.POST("/packet", packetHandler.UploadPacket)
	api.POST("/packets/batch", packetHandler.UploadBatch)

	api.GET("/devices", deviceHandler.ListDevices)   
	api.GET("/devices/:id/history", historyHandler.GetHistory)

	api.GET("/mobiles", mobileDeviceHandler.GetMobiles)
	api.GET("/mobiles/:mobileId/sensors", mobileDeviceHandler.GetSensorsForMobile)

	api.POST("/command", websocket.DashboardCommandHandler(ingestHub))
}

	ws := router.Group("/ws")
	{
		ws.GET("", func(c *gin.Context) {
			websocket.ServeWs(hub, c.Writer, c.Request)
		})

		ws.GET("/ingest", func(c *gin.Context) {
			websocket.ServeIngest(ingestHub, c.Request.Context(), c.Writer, c.Request)
		})
	}
}
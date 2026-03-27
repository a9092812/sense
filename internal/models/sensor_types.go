package models

type SensorType string

const (
	SensorSHT40         SensorType = "SHT40"
	SensorLux           SensorType = "LuxSensor"
	SensorLIS2DH        SensorType = "LIS2DH"
	SensorSoil          SensorType = "SoilSensor"
	SensorSpeedDistance SensorType = "SpeedDistance"
	SensorAmmonia       SensorType = "AmmoniaSensor"
	SensorTempLogger    SensorType = "TempLogger"
	SensorDataLogger    SensorType = "DataLogger"
	SensorSen6x         SensorType = "SEN6x"
)

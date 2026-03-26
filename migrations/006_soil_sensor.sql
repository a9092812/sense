CREATE TABLE IF NOT EXISTS soil_sensor_data (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    nitrogen DOUBLE PRECISION, phosphorus DOUBLE PRECISION, potassium DOUBLE PRECISION,
    moisture DOUBLE PRECISION, temperature DOUBLE PRECISION, ec DOUBLE PRECISION,
    ph DOUBLE PRECISION, salinity DOUBLE PRECISION
);
SELECT create_hypertable('soil_sensor_data', 'time', if_not_exists => TRUE);
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_soil_device_time') THEN
        CREATE INDEX idx_soil_device_time ON soil_sensor_data(device_id, time DESC);
    END IF;
END $$;
CREATE TABLE IF NOT EXISTS lux_sensor_data (

    time TIMESTAMPTZ NOT NULL,

    device_id TEXT NOT NULL,

    lux DOUBLE PRECISION,

    raw_data TEXT

);

SELECT create_hypertable(
    'lux_sensor_data',
    'time',
    if_not_exists => TRUE
);

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_lux_device_time') THEN
        CREATE INDEX idx_lux_device_time ON lux_sensor_data(device_id, time DESC);
    END IF;
END $$;
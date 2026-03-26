CREATE TABLE IF NOT EXISTS temp_logger_data (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    temperature DOUBLE PRECISION, humidity DOUBLE PRECISION,
    raw_temperature INTEGER, raw_humidity INTEGER,
    raw_data TEXT, device_address TEXT
);
SELECT create_hypertable('temp_logger_data', 'time', if_not_exists => TRUE);
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_temp_logger_device_time') THEN
        CREATE INDEX idx_temp_logger_device_time ON temp_logger_data(device_id, time DESC);
    END IF;
END $$;
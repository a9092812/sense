CREATE TABLE IF NOT EXISTS data_logger_data (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    current_packet_id INTEGER, last_packet_id INTEGER,
    accel_data JSONB, raw_data TEXT
);
SELECT create_hypertable('data_logger_data', 'time', if_not_exists => TRUE);
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_datalogger_device_time') THEN
        CREATE INDEX idx_datalogger_device_time ON data_logger_data(device_id, time DESC);
    END IF;
END $$;
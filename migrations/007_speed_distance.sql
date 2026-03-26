CREATE TABLE IF NOT EXISTS speed_distance_data (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    speed DOUBLE PRECISION, distance DOUBLE PRECISION
);
SELECT create_hypertable('speed_distance_data', 'time', if_not_exists => TRUE);
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_speed_device_time') THEN
        CREATE INDEX idx_speed_device_time ON speed_distance_data(device_id, time DESC);
    END IF;
END $$;
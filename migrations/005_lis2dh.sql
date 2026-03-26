CREATE TABLE IF NOT EXISTS lis2dh_data (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    x DOUBLE PRECISION, y DOUBLE PRECISION, z DOUBLE PRECISION
);
SELECT create_hypertable('lis2dh_data', 'time', if_not_exists => TRUE);
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_lis2dh_device_time') THEN
        CREATE INDEX idx_lis2dh_device_time ON lis2dh_data(device_id, time DESC);
    END IF;
END $$;
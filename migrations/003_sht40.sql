CREATE TABLE IF NOT EXISTS sht40_data (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION
);

SELECT create_hypertable(
    'sht40_data',
    'time',
    if_not_exists => TRUE
);

-- Use this block for the index:
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_sht40_device_time') THEN
        CREATE INDEX idx_sht40_device_time ON sht40_data(device_id, time DESC);
    END IF;
END $$;
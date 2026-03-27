CREATE TABLE IF NOT EXISTS sen6x_data (
    time TIMESTAMPTZ NOT NULL,
    device_id TEXT NOT NULL,
    pm1 DOUBLE PRECISION,
    pm25 DOUBLE PRECISION,
    pm4 DOUBLE PRECISION,
    pm10 DOUBLE PRECISION,
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION,
    co2 DOUBLE PRECISION,
    voc DOUBLE PRECISION,
    nox DOUBLE PRECISION
);

-- Convert to hypertable for TimescaleDB optimization
SELECT create_hypertable('sen6x_data', 'time', if_not_exists => TRUE);

-- Add index for device_id
CREATE INDEX IF NOT EXISTS idx_sen6x_device_id ON sen6x_data(device_id, time DESC);

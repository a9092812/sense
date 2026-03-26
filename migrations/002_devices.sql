CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    device_id TEXT UNIQUE NOT NULL,
    device_address TEXT,
    sensor_type sensor_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
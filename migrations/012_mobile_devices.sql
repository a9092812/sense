CREATE TABLE IF NOT EXISTS mobile_devices (
    id SERIAL PRIMARY KEY,
    mobile_id TEXT UNIQUE NOT NULL,
    name TEXT,
    last_seen TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE devices ADD COLUMN IF NOT EXISTS mobile_id TEXT;
ALTER TABLE devices ADD COLUMN IF NOT EXISTS last_seen TIMESTAMPTZ DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_devices_mobile_id ON devices(mobile_id);

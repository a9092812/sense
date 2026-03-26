-- Migration to track mobile-sensor affinity and live status
CREATE TABLE IF NOT EXISTS mobile_sensor_affinity (
    mobile_id TEXT NOT NULL,
    device_id TEXT NOT NULL,
    last_seen TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (mobile_id, device_id)
);

CREATE INDEX IF NOT EXISTS idx_mobile_sensor_affinity_last_seen ON mobile_sensor_affinity(last_seen);

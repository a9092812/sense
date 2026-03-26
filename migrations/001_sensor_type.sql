DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sensor_type') THEN
        CREATE TYPE sensor_type AS ENUM (
            'SHT40',
            'LuxSensor',
            'LIS2DH',
            'SoilSensor',
            'SpeedDistance',
            'AmmoniaSensor',
            'TempLogger',
            'DataLogger'
        );
    END IF;
END $$;
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_database WHERE datname = 'aggregator_db'
    ) THEN
        CREATE DATABASE aggregator_db;
    END IF;
END $$;
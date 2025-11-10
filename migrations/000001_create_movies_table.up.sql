CREATE TABLE IF NOT EXISTS movies (
    id BIGSERIAL PRIMARY KEY,       
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
    title TEXT NOT NULL,
    year INT NOT NULL,
    runtime INT NOT NULL,
    genres TEXT[],
    version INT NOT NULL DEFAULT 1
);
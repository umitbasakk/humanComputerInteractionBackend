CREATE TABLE IF NOT EXISTS requests (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    started_date VARCHAR(50) NOT NULL,
    end_date VARCHAR(50) NOT NULL,
    hash_tag VARCHAR(50) NOT NULL,
    category INT NOT NULL,
    quantity_limit INT NOT NULL,
    request_status INT NOT NULL,
    created_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP
);
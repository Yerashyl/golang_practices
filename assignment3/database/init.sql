CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    age INT,
    is_active BOOLEAN DEFAULT TRUE,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    action VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Optional: Insert some initial data
INSERT INTO users (name, email, age) VALUES 
('Alice Smith', 'alice@example.com', 30),
('Bob Jones', 'bob@example.com', 25)
ON CONFLICT (email) DO NOTHING;

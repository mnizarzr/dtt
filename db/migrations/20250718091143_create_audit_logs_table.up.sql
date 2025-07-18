CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(255) NOT NULL,
    target_resource VARCHAR(100) NOT NULL,
    target_id UUID,
    details JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

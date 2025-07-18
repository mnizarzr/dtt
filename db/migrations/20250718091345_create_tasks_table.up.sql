CREATE TABLE IF NOT EXISTS tasks (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title          VARCHAR(100),
    description    TEXT,
    status         VARCHAR(20) CHECK (status IN ('pending', 'in_progress', 'completed')) DEFAULT 'pending',
    priority       VARCHAR(10) CHECK (priority IN ('low', 'medium', 'high')) DEFAULT 'medium',
    due_date       DATE,
    project_id     UUID REFERENCES projects(id) ON DELETE CASCADE,
    created_by     UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_to    UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP DEFAULT NOW()
);

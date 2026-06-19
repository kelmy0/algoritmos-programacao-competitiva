CREATE TYPE difficulty_level AS ENUM ('beginner', 'intermediate', 'advanced', 'expert');

CREATE TABLE algorithms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(8) UNIQUE NOT NULL,
    slug VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    category VARCHAR(64) NOT NULL,
    difficulty difficulty_level NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_atualizar_updated_at
BEFORE UPDATE ON algorithms
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
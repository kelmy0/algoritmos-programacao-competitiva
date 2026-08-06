-- ENUM AND FUNCTIONS
CREATE TYPE difficulty_level AS ENUM ('beginner', 'intermediate', 'advanced', 'expert');
CREATE TYPE status_text AS ENUM('approved', 'pending', 'rejected', 'deleted');

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
--------------------------------------
-- GENERAL TABLES
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(16) NOT NULL UNIQUE,
    is_employee BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(64) NOT NULL UNIQUE
);

CREATE TABLE role_permissions (
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);
--------------------------------------
-- USERS TABLE
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL CHECK (char_length(name) >= 6),
    username VARCHAR(32) NOT NULL UNIQUE CHECK (char_length(username) >= 6),
    email VARCHAR(128) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    recovery_token_hash CHAR(128),
    recovery_token_expires_at TIMESTAMPTZ,
    enable BOOLEAN NOT NULL DEFAULT TRUE,
    two_factor_authentication BOOLEAN NOT NULL DEFAULT FALSE,
    two_factor_secret VARCHAR(255),
    role_id INT NOT NULL,
    failed_attempts INT NOT NULL DEFAULT 0,
    last_login TIMESTAMPTZ,
    last_logout_all TIMESTAMPTZ DEFAULT NULL,
    blocked_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);

CREATE INDEX idx_users_recovery_token ON users(recovery_token_hash) WHERE recovery_token_hash IS NOT NULL;
CREATE INDEX idx_users_role_id ON users(role_id);

CREATE TRIGGER trigger_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- TABLES THAT DEPENDS USERS
CREATE TABLE algorithms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(8) UNIQUE NOT NULL,
    slug VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    category VARCHAR(64) NOT NULL,
    difficulty difficulty_level NOT NULL,
    content TEXT NOT NULL,
    status status_text DEFAULT 'pending',
    author_id UUID,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_algorithms_author_id ON algorithms(author_id);
CREATE INDEX idx_algorithms_status ON algorithms(status);
CREATE INDEX idx_algorithms_author_status_updated ON algorithms (author_id, status, updated_at DESC);
CREATE INDEX idx_algorithms_status_updated ON algorithms (status, updated_at DESC);

CREATE TRIGGER trigger_algorithms_updated_at
BEFORE UPDATE ON algorithms
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

--------------------------------------

CREATE TABLE refresh_tokens (
    id CHAR(32) PRIMARY KEY,
    user_id UUID NOT NULL,
    family_id CHAR(32) NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_family_id ON refresh_tokens(family_id);

--------------------------------------

CREATE TABLE user_social_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    provider VARCHAR(50) NOT NULL,
    social_user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_provider_social_id UNIQUE (provider, social_user_id)
);

CREATE INDEX idx_user_social_accounts_user_id ON user_social_accounts(user_id);
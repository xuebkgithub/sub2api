-- 053_add_ldap_tables.sql
-- Add LDAP authentication support

-- 1. ldap_configs 表
CREATE TABLE IF NOT EXISTS ldap_configs (
    id                      BIGSERIAL PRIMARY KEY,
    server_url              VARCHAR(255) NOT NULL,
    bind_dn                 VARCHAR(255) NOT NULL,
    bind_password_encrypted TEXT NOT NULL,
    base_dn                 VARCHAR(255) NOT NULL,
    user_filter             VARCHAR(255) NOT NULL DEFAULT '(uid=%s)',
    enabled                 BOOLEAN NOT NULL DEFAULT FALSE,
    tls_enabled             BOOLEAN NOT NULL DEFAULT FALSE,
    tls_skip_verify         BOOLEAN NOT NULL DEFAULT FALSE,
    config_source           VARCHAR(50) NOT NULL DEFAULT 'database',
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ldap_configs_enabled ON ldap_configs(enabled);

COMMENT ON TABLE ldap_configs IS 'LDAP server configurations';
COMMENT ON COLUMN ldap_configs.server_url IS 'LDAP server URL (e.g., ldap://ldap.example.com:389)';
COMMENT ON COLUMN ldap_configs.bind_password_encrypted IS 'Encrypted bind password (AES-256-GCM)';
COMMENT ON COLUMN ldap_configs.config_source IS 'Configuration source: env or database';

-- 2. ldap_users 表
CREATE TABLE IF NOT EXISTS ldap_users (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ldap_username  VARCHAR(255) NOT NULL,
    ldap_dn        VARCHAR(500) NOT NULL,
    last_sync_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id),
    UNIQUE(ldap_username)
);

COMMENT ON TABLE ldap_users IS 'LDAP user associations';
COMMENT ON COLUMN ldap_users.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN ldap_users.ldap_username IS 'LDAP username (uid)';
COMMENT ON COLUMN ldap_users.ldap_dn IS 'LDAP Distinguished Name';

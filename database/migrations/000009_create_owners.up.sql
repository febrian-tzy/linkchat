-- Migration: 000009_create_owners.up.sql
-- Owner: platform administrator (separate from tenant users)

CREATE TABLE IF NOT EXISTS `owners` (
    `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `full_name`        VARCHAR(255)    NOT NULL,
    `email`            VARCHAR(255)    NOT NULL,
    `username`         VARCHAR(100)    NOT NULL,
    `password_hash`    VARCHAR(512)    NOT NULL,
    `is_active`        TINYINT(1)      NOT NULL DEFAULT 1,
    `last_login_at`    DATETIME        NULL DEFAULT NULL,
    `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_owners_username` (`username`),
    UNIQUE KEY `uq_owners_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Default owner account: username=owner, password=owner123
-- bcrypt hash for "owner123" at cost 12
INSERT INTO `owners` (full_name, email, username, password_hash) VALUES (
    'Platform Owner',
    'owner@linkchat.id',
    'owner',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/RK.s5uOeK'
);

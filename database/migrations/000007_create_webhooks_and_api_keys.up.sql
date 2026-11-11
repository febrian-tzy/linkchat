-- Migration: 000007_create_webhooks_and_api_keys.up.sql

CREATE TABLE IF NOT EXISTS `webhooks` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT UNSIGNED NOT NULL,
    `url`         VARCHAR(512)    NOT NULL,
    `secret`      VARCHAR(255)    NULL DEFAULT NULL,
    `events`      JSON            NOT NULL COMMENT 'Array of event names',
    `is_active`   TINYINT(1)      NOT NULL DEFAULT 1,
    `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_webhooks_tenant` (`tenant_id`),
    CONSTRAINT `fk_webhooks_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `api_keys` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT UNSIGNED NOT NULL,
    `name`        VARCHAR(255)    NOT NULL,
    `key_hash`    VARCHAR(255)    NOT NULL COMMENT 'SHA256 hash of the API key',
    `key_prefix`  VARCHAR(10)     NOT NULL COMMENT 'First 8 chars for display',
    `is_active`   TINYINT(1)      NOT NULL DEFAULT 1,
    `last_used_at` DATETIME       NULL DEFAULT NULL,
    `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_api_keys_hash` (`key_hash`),
    INDEX `idx_api_keys_tenant` (`tenant_id`),
    CONSTRAINT `fk_api_keys_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

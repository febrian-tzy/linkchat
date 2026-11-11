-- Migration: 000008_create_automations.up.sql
CREATE TABLE IF NOT EXISTS `automations` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT UNSIGNED NOT NULL,
    `name`        VARCHAR(255)    NOT NULL,
    `description` TEXT            NULL DEFAULT NULL,
    `event`       VARCHAR(100)    NOT NULL COMMENT 'e.g. conversation_created, message_received',
    `conditions`  JSON            NULL DEFAULT NULL,
    `actions`     JSON            NOT NULL,
    `is_active`   TINYINT(1)      NOT NULL DEFAULT 1,
    `run_count`   INT UNSIGNED    NOT NULL DEFAULT 0,
    `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_automations_tenant` (`tenant_id`),
    CONSTRAINT `fk_automations_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

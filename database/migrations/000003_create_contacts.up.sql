-- Migration: 000003_create_contacts.up.sql
CREATE TABLE IF NOT EXISTS `contacts` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT UNSIGNED NOT NULL,
    `name`        VARCHAR(255)    NOT NULL,
    `phone`       VARCHAR(30)     NOT NULL,
    `email`       VARCHAR(255)    NULL DEFAULT NULL,
    `avatar_url`  VARCHAR(512)    NULL DEFAULT NULL,
    `notes`       TEXT            NULL DEFAULT NULL,
    `is_blocked`  TINYINT(1)      NOT NULL DEFAULT 0,
    `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_contacts_tenant_phone` (`tenant_id`, `phone`),
    INDEX `idx_contacts_tenant_id` (`tenant_id`),
    CONSTRAINT `fk_contacts_tenant_id`
        FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Migration: 000002_create_tenant_users.up.sql
-- Creates the tenant_users table

CREATE TABLE IF NOT EXISTS `tenant_users` (
    `id`                  BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    `tenant_id`           BIGINT UNSIGNED  NOT NULL,
    `full_name`           VARCHAR(255)     NOT NULL,
    `email`               VARCHAR(255)     NOT NULL,
    `username`            VARCHAR(100)     NOT NULL,
    `password_hash`       VARCHAR(512)     NOT NULL,
    `role`                ENUM('admin', 'agent', 'supervisor') NOT NULL DEFAULT 'admin',
    `is_active`           TINYINT(1)       NOT NULL DEFAULT 1,
    `email_verified_at`   DATETIME         NULL DEFAULT NULL,
    `last_login_at`       DATETIME         NULL DEFAULT NULL,
    `created_at`          DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_tenant_users_username` (`username`),
    UNIQUE KEY `uq_tenant_users_email` (`email`),
    CONSTRAINT `fk_tenant_users_tenant_id`
        FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
        ON DELETE CASCADE ON UPDATE CASCADE,
    INDEX `idx_tenant_users_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

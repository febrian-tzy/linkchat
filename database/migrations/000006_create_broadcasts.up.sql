-- Migration: 000006_create_broadcasts.up.sql
CREATE TABLE IF NOT EXISTS `broadcasts` (
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tenant_id`    BIGINT UNSIGNED NOT NULL,
    `inbox_id`     BIGINT UNSIGNED NOT NULL,
    `name`         VARCHAR(255)    NOT NULL,
    `message`      TEXT            NOT NULL,
    `status`       ENUM('draft','scheduled','sending','sent','failed') NOT NULL DEFAULT 'draft',
    `scheduled_at` DATETIME        NULL DEFAULT NULL,
    `sent_at`      DATETIME        NULL DEFAULT NULL,
    `total`        INT UNSIGNED    NOT NULL DEFAULT 0,
    `sent_count`   INT UNSIGNED    NOT NULL DEFAULT 0,
    `failed_count` INT UNSIGNED    NOT NULL DEFAULT 0,
    `created_by`   BIGINT UNSIGNED NOT NULL,
    `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_broadcasts_tenant` (`tenant_id`),
    CONSTRAINT `fk_broadcasts_tenant`     FOREIGN KEY (`tenant_id`) REFERENCES `tenants`      (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_broadcasts_inbox`      FOREIGN KEY (`inbox_id`)  REFERENCES `inboxes`      (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_broadcasts_created_by` FOREIGN KEY (`created_by`) REFERENCES `tenant_users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `broadcast_recipients` (
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `broadcast_id` BIGINT UNSIGNED NOT NULL,
    `contact_id`   BIGINT UNSIGNED NOT NULL,
    `status`       ENUM('pending','sent','failed','delivered','read') NOT NULL DEFAULT 'pending',
    `sent_at`      DATETIME NULL DEFAULT NULL,
    `error`        VARCHAR(512) NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_broadcast_recipients_broadcast` (`broadcast_id`),
    CONSTRAINT `fk_broadcast_recipients_broadcast`
        FOREIGN KEY (`broadcast_id`) REFERENCES `broadcasts` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_broadcast_recipients_contact`
        FOREIGN KEY (`contact_id`)   REFERENCES `contacts`   (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

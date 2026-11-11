-- Migration: 000004_create_inboxes.up.sql
CREATE TABLE IF NOT EXISTS `inboxes` (
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tenant_id`    BIGINT UNSIGNED NOT NULL,
    `name`         VARCHAR(255)    NOT NULL,
    `channel_type` ENUM('whatsapp','web','email','api') NOT NULL DEFAULT 'whatsapp',
    `is_enabled`   TINYINT(1)      NOT NULL DEFAULT 1,
    `config`       JSON            NULL DEFAULT NULL,
    `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_inboxes_tenant_id` (`tenant_id`),
    CONSTRAINT `fk_inboxes_tenant_id`
        FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Inbox members: which agents handle which inbox
CREATE TABLE IF NOT EXISTS `inbox_members` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `inbox_id`   BIGINT UNSIGNED NOT NULL,
    `user_id`    BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_inbox_members` (`inbox_id`, `user_id`),
    CONSTRAINT `fk_inbox_members_inbox` FOREIGN KEY (`inbox_id`) REFERENCES `inboxes` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_inbox_members_user`  FOREIGN KEY (`user_id`)  REFERENCES `tenant_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

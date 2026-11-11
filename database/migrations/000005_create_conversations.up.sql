-- Migration: 000005_create_conversations.up.sql
CREATE TABLE IF NOT EXISTS `conversations` (
    `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `tenant_id`        BIGINT UNSIGNED NOT NULL,
    `inbox_id`         BIGINT UNSIGNED NOT NULL,
    `contact_id`       BIGINT UNSIGNED NOT NULL,
    `assigned_user_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `status`           ENUM('open','resolved','pending','snoozed') NOT NULL DEFAULT 'open',
    `unread_count`     INT UNSIGNED    NOT NULL DEFAULT 0,
    `last_message_at`  DATETIME        NULL DEFAULT NULL,
    `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_conversations_tenant`  (`tenant_id`),
    INDEX `idx_conversations_inbox`   (`inbox_id`),
    INDEX `idx_conversations_contact` (`contact_id`),
    INDEX `idx_conversations_status`  (`status`),
    CONSTRAINT `fk_conv_tenant`   FOREIGN KEY (`tenant_id`)        REFERENCES `tenants`      (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_conv_inbox`    FOREIGN KEY (`inbox_id`)         REFERENCES `inboxes`      (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_conv_contact`  FOREIGN KEY (`contact_id`)       REFERENCES `contacts`     (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_conv_assignee` FOREIGN KEY (`assigned_user_id`) REFERENCES `tenant_users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `messages` (
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `conversation_id` BIGINT UNSIGNED NOT NULL,
    `tenant_id`       BIGINT UNSIGNED NOT NULL,
    `sender_type`     ENUM('agent','contact','system','bot') NOT NULL DEFAULT 'agent',
    `sender_id`       BIGINT UNSIGNED NULL DEFAULT NULL,
    `content`         TEXT            NOT NULL,
    `content_type`    ENUM('text','image','document','audio','video','template') NOT NULL DEFAULT 'text',
    `media_url`       VARCHAR(512)    NULL DEFAULT NULL,
    `is_read`         TINYINT(1)      NOT NULL DEFAULT 0,
    `external_id`     VARCHAR(255)    NULL DEFAULT NULL COMMENT 'WhatsApp message ID',
    `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_messages_conversation` (`conversation_id`),
    INDEX `idx_messages_tenant`       (`tenant_id`),
    CONSTRAINT `fk_messages_conversation`
        FOREIGN KEY (`conversation_id`) REFERENCES `conversations` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

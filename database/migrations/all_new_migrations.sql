-- === 000003_create_contacts.up.sql ===
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


-- === 000004_create_inboxes.up.sql ===
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


-- === 000005_create_conversations.up.sql ===
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


-- === 000006_create_broadcasts.up.sql ===
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


-- === 000007_create_webhooks_and_api_keys.up.sql ===
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


-- === 000008_create_automations.up.sql ===
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




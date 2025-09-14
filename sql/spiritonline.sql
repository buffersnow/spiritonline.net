-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: MySQL
-- Generated at: 2025-09-14T13:08:53.696Z

CREATE TABLE `accounts` (
  `user_id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `screenname` varchar(64) NOT NULL,
  `mail` varchar(128) NOT NULL,
  `priv_status` tinyint NOT NULL DEFAULT 0,
  `acs_status` tinyint NOT NULL DEFAULT 0,
  `is_verified` boolean NOT NULL DEFAULT false,
  `last_updated_on` datetime(3) NOT NULL
);

CREATE TABLE `users` (
  `user_id` bigint PRIMARY KEY NOT NULL,
  `avatar_hash` varchar(32),
  `contact_list_rev` bigint NOT NULL DEFAULT 0,
  `last_cl_update` datetime(3),
  `last_login` datetime(3),
  `first_login` datetime(3),
  `signup_date` datetime(3) NOT NULL,
  `last_updated_on` datetime(3) NOT NULL
);

CREATE TABLE `profiles` (
  `user_id` bigint PRIMARY KEY NOT NULL,
  `status_message` text,
  `away_message` text,
  `offline_message` text,
  `last_updated_on` datetime(3) NOT NULL
);

CREATE TABLE `contacts` (
  `sender_id` bigint NOT NULL,
  `friend_id` bigint NOT NULL,
  `reason` mediumtext,
  `is_blocked` boolean NOT NULL DEFAULT false,
  `added_on` datetime(3) NOT NULL,
  `last_updated_on` datetime(3) NOT NULL
);

CREATE TABLE `app_passwords` (
  `user_id` bigint NOT NULL,
  `type` varchar(32) NOT NULL,
  `contents` mediumtext,
  `last_updated_on` datetime(3) NOT NULL
);

CREATE TABLE `wfc_accounts` (
  `wfc_id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `linked_id` bigint,
  `nand_ids` json,
  `console_ids` json,
  `ip_addresses` json,
  `last_updated_on` datetime(3) NOT NULL
);

CREATE TABLE `wfc_suspensions` (
  `audit_id` bigint PRIMARY KEY NOT NULL,
  `wfc_id` bigint NOT NULL,
  `moderator_id` bigint NOT NULL,
  `ban_message` longtext,
  `ban_reason` longtext,
  `last_banned_on` datetime(3) NOT NULL,
  `ban_expires_on` datetime(3)
);

CREATE INDEX `idx_screenname` ON `accounts` (`screenname`);

CREATE INDEX `idx_mail` ON `accounts` (`mail`);

CREATE INDEX `idx_sender_id` ON `contacts` (`sender_id`);

CREATE INDEX `idx_friend_id` ON `contacts` (`friend_id`);

CREATE UNIQUE INDEX `udx_sender_id_friend_id` ON `contacts` (`sender_id`, `friend_id`);

CREATE UNIQUE INDEX `udx_user_id_type` ON `app_passwords` (`user_id`, `type`);

CREATE INDEX `idx_linked_id` ON `wfc_accounts` (`linked_id`);

CREATE INDEX `idx_wfc_id` ON `wfc_suspensions` (`wfc_id`);

CREATE INDEX `idx_moderator_id` ON `wfc_suspensions` (`moderator_id`);

ALTER TABLE `users` ADD FOREIGN KEY (`user_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `profiles` ADD FOREIGN KEY (`user_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `contacts` ADD FOREIGN KEY (`sender_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `contacts` ADD FOREIGN KEY (`friend_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `app_passwords` ADD FOREIGN KEY (`user_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `wfc_accounts` ADD FOREIGN KEY (`linked_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `wfc_suspensions` ADD FOREIGN KEY (`wfc_id`) REFERENCES `wfc_accounts` (`wfc_id`);

ALTER TABLE `wfc_suspensions` ADD FOREIGN KEY (`moderator_id`) REFERENCES `accounts` (`user_id`);

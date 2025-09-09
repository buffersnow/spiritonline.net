-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: MySQL
-- Generated at: 2025-09-09T09:53:48.863Z

CREATE TABLE `accounts` (
  `user_id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `screenname` varchar(64) NOT NULL,
  `mail` varchar(128) NOT NULL,
  `flags` bigint NOT NULL DEFAULT 0,
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

CREATE INDEX `idx_accounts_screenname` ON `accounts` (`screenname`);

CREATE INDEX `idx_accounts_mail` ON `accounts` (`mail`);

CREATE UNIQUE INDEX `idx_sender_id_friend_id` ON `contacts` (`sender_id`, `friend_id`);

CREATE UNIQUE INDEX `idx_user_id_type` ON `app_passwords` (`user_id`, `type`);

ALTER TABLE `users` ADD FOREIGN KEY (`user_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `profiles` ADD FOREIGN KEY (`user_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `contacts` ADD FOREIGN KEY (`sender_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `contacts` ADD FOREIGN KEY (`friend_id`) REFERENCES `accounts` (`user_id`);

ALTER TABLE `app_passwords` ADD FOREIGN KEY (`user_id`) REFERENCES `accounts` (`user_id`);

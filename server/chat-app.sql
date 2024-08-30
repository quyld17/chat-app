CREATE TABLE `users` (
  `id` integer PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `username` varchar(255) UNIQUE NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE `chats` (
  `id` integer PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` varchar(255),
  `is_group` boolean DEFAULT false,
  `created_at` timestamp NOT NULL
);

CREATE TABLE `chat_participants` (
  `id` integer PRIMARY KEY NOT NULL,
  `room_id` integer NOT NULL,
  `user_id` integer NOT NULL
);

CREATE TABLE `messages` (
  `id` integer PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `room_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE `status` (
  `user_id` integer NOT NULL,
  `is_online` boolean DEFAULT false
);

ALTER TABLE `chat_participants` ADD FOREIGN KEY (`room_id`) REFERENCES `chats` (`id`);

ALTER TABLE `chat_participants` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `messages` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `messages` ADD FOREIGN KEY (`room_id`) REFERENCES `chats` (`id`);

ALTER TABLE `status` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

CREATE TABLE `users` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `email`        VARCHAR(255)    NOT NULL,
  `password`     VARCHAR(255)    NOT NULL,
  `first_name`   VARCHAR(255)    NOT NULL,
  `last_name`    VARCHAR(255)    NOT NULL,
  `phone_number` VARCHAR(20)     NOT NULL,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at`   DATETIME        NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_email` (`email`),
  UNIQUE KEY `uni_users_phone_number` (`phone_number`),
  INDEX `idx_users_deleted_at` (`deleted_at`)
);
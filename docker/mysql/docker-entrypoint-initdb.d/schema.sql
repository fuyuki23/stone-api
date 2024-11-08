CREATE TABLE `user`
(
    `id`        binary(16) PRIMARY KEY,
    `email`     varchar(100) UNIQUE NOT NULL,
    `password`  varchar(100)        NOT NULL,
    `name`      varchar(50)         NULL,
    `create_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
    `update_at` timestamp DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE `diary`
(
    `id`        binary(16) PRIMARY KEY,
    `user_id`   binary(16) NOT NULL,
    `title`     varchar(255)  NOT NULL,
    `content`   varchar(1024) NOT NULL,
    `mood`      varchar(20)   NOT NULL,
    `create_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
    `update_at` timestamp DEFAULT (CURRENT_TIMESTAMP)
);

ALTER TABLE `diary`
    ADD CONSTRAINT `fk_diary_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

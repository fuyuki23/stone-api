CREATE TABLE `user`
(
    `id`        binary(16) PRIMARY KEY,
    `email`     varchar(100) UNIQUE NOT NULL,
    `password`  varchar(100)        NOT NULL,
    `name`      varchar(50)         NOT NULL,
    `create_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
    `update_at` timestamp DEFAULT (CURRENT_TIMESTAMP)
);

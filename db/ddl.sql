CREATE DATABASE IF NOT EXISTS chat;
USE chat;
CREATE TABLE IF NOT EXISTS users (
    id int PRIMARY KEY AUTO_INCREMENT,
    user_id varchar(36) NOT NULL UNIQUE,
    name  varchar(36) NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms (
    id int PRIMARY KEY AUTO_INCREMENT,
    name varchar(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id int PRIMARY KEY AUTO_INCREMENT,
    message varchar(140) NOT NULL,
    created_at datetime NOT NULL,
    user_id int NOT NULL,
    room_id int NOT NULL,
    CONSTRAINT
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT
        FOREIGN KEY (room_id)
        REFERENCES rooms (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

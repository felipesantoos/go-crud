CREATE DATABASE `crud`;

USE `crud`;

CREATE TABLE `names` (
    `id` INT(6) NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(30) NOT NULL,
    `email` VARCHAR(30) NOT NULL,
    CONSTRAINT `namesId` PRIMARY KEY (`id`)
);

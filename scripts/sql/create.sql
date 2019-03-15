DROP DATABASE IF EXISTS `familytree`;
CREATE DATABASE `familytree`;
USE `familytree`;
CREATE TABLE people (
    `UUID` VARCHAR(128) NOT NULL,
    `Name` VARCHAR(255),
    `Fnam` VARCHAR(50),
    `Mnam` VARCHAR(255),
    `Snam` VARCHAR(50),
    `Sex` VARCHAR(10),
    `DOB` VARCHAR(10),
    `DOD` VARCHAR(10),
    PRIMARY KEY (`UUID`)
);
CREATE TABLE children (
    `id` INT NOT NULL,
    `parentId` VARCHAR(128) NOT NULL,
    `childId` VARCHAR(128) NOT NULL,
    PRIMARY KEY(`id`)
);
CREATE TABLE partners (
    `id` INT NOT NULL,
    `bitchId` VARCHAR(128) NOT NULL,
    `butchId` VARCHAR(128) NOT NULL,
    PRIMARY KEY(`id`)
);

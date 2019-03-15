CREATE DATABASE family-tree;
USE family-tree;

CREATE TABLE people (
  UUID VARCHAR(128) NOT NULL,
  Name VARCHAR(255),
  Fnam VARCHAR(50),
  Mnam VARCHAR(255),
  Snam VARCHAR(50),
  Sex VARCHAR(10),
  DOB VARCHAR(10),
  DOD VARCHAR(10),
  SpouseID VARCHAR(128),
  PRIMARY KEY (UUID)
);

Children      []string  `yaml:"children"`
Parents       []string  `yaml:"parents"`

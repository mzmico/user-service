CREATE DATABASE IF NOT EXISTS db_user DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE db_user;
CREATE TABLE tb_account (
  app_id VARCHAR(32) NOT NULL,
  uid VARCHAR(32) NOT NULL,
  account VARCHAR(128) NOT NULL ,
  certificate VARCHAR(256) NOT NULL,
  type INT NOT NULL ,
  PRIMARY KEY (app_id,uid,account)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;




CREATE TABLE tb_user (
  app_id VARCHAR(32) NOT NULL,
  uid VARCHAR(32) NOT NULL,
  name VARCHAR(32) NOT NULL DEFAULT '',
  nick VARCHAR(32) NOT NULL DEFAULT '',
  avatar TEXT NOT NULL,
  extend JSON NOT NULL,
  PRIMARY KEY (app_id,uid)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

DROP TABLE tb_account;
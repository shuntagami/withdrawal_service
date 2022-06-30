create database if not exists `withdrawal_service` character set utf8mb4 collate utf8mb4_bin;

use `withdrawal_service`;

drop table if exists users;
create table if not exists users
(
  id int unsigned not null primary key auto_increment,
  name varchar(128) not null
) character set utf8mb4 collate utf8mb4_bin;

drop table if exists histories;
create table if not exists histories
(
  id int unsigned not null primary key auto_increment,
  user_id int unsigned not null,
  amount int not null,
  CONSTRAINT fk_histories_users FOREIGN KEY (user_id) REFERENCES users (id)
) character set utf8mb4 collate utf8mb4_bin;

insert into users (name)
values ('user1'),
       ('user2');

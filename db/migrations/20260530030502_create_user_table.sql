-- migrate:up
create table if not exists users(
  id int auto_increment primary key,
  email varchar(255) not null unique,
  username varchar(255) not null,
  password varchar(500) not null,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp
)

-- migrate:down
drop table if exists users

-- migrate:up
create table if not exists refresh_token(
  id int primary key auto_increment,
  user_id int not null,
  refresh_token text not null,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  constraint fk_user_id_refresh_token foreign key (user_id) references users(id)
)

-- migrate:down
drop table if exists refresh_token

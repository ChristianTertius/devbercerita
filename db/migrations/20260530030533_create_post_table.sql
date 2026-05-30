-- migrate:up
create table if not exists posts(
  id int auto_increment primary key,
  user_id int not null,
  title varchar(255) not null, 
  content longtext not null,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  deleted_at timestamp null,
  constraint fk_user_id_posts foreign key (user_id) references users(id)
)

-- migrate:down
drop table if exists posts

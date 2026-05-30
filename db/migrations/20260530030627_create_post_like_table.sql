-- migrate:up
create table if not exists post_likes(
  id int primary key auto_increment,
  post_id int not null,
  user_id int not null,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  constraint fk_user_id_post_likes foreign key (user_id) references users(id),
  constraint fk_post_id_post_likes foreign key (post_id) references posts(id)
)

-- migrate:down
drop table if exists post_likes

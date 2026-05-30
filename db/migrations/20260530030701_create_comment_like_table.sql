-- migrate:up
create table if not exists comment_likes(
  id int primary key auto_increment,
  comment_id int not null,
  user_id int not null,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  constraint fk_user_id_comment_likes foreign key (user_id) references users(id),
  constraint fk_comment_id_comment_likes foreign key (comment_id) references comments(id)
)

-- migrate:down
drop table if exists comment_likes

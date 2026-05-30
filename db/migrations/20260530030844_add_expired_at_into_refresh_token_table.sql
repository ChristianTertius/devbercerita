-- migrate:up
alter table refresh_token add expired_at timestamp null after refresh_token;

-- migrate:down
alter table refresh_token drop column expired_at;

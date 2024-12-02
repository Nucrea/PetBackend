create table if not exists action_tokens (
    id int generated always as identity,
    user_id int,
    value text,
    target int,
    expiration timestamp,

    primary key(id)
);
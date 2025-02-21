create table if not exists users (
    id integer primary key generated always as identity,
    email varchar(256) unique not null,
    secret varchar(256) not null,
    full_name varchar(256) not null,
    email_verified boolean not null default false,
    active boolean,
    created_at timestamp,
    updated_at timestamp
);

create index if not exists idx_users_email on users(email);

create or replace trigger trg_user_created
    before insert on users
    for each row
    execute function trg_proc_row_created();

create or replace trigger trg_user_updated
    before update on users
    for each row
    when(new is distinct from old)
    execute function trg_proc_row_updated();
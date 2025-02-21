create table if not exists users (
    id integer primary key generated always as identity,
    email varchar(256) unique not null,
    secret varchar(256) not null,
    full_name varchar(256) not null,
    email_verified boolean not null default false,
    created_at timestamp,
    updated_at timestamp
);

create index if not exists users_email_idx on users(email);

create or replace function set_created_at()
returns trigger as $$
begin
    new.created_at = now();
	new.updated_at = now();
    return new;
end;
$$ language plpgsql;

create or replace trigger on_user_created
    before insert on users
    for each row
    execute function set_created_at();

create or replace function set_updated_at()
returns trigger as $$
begin
    if new is distinct from old then
        new.updated_at = now();
    end if;
    return new; 
end;
$$ language plpgsql;

create or replace trigger on_user_updated
    before update on users
    for each row
    when(new is distinct from old)
    execute function set_updated_at();
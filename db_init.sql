create table users (
    id int generated always as identity,
    email text unique not null,
    secret text not null,
    name text not null,

    primary key (id)
);
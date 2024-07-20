create table users (
    id int generated always as identity,
    login text unique not null,
    secret text not null,
    name text not null,

    primary key (id)
);
create table users (
    id text primary key,
    name text not null,
    password text,
    created_at datetime not null,
    updated_at datetime not null
);
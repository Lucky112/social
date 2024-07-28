create schema scl;

create table scl.users (
    id bigserial PRIMARY KEY,
    login varchar(50) UNIQUE NOT NULL,
    password varchar(100) NOT NULL,
    email varchar(300) UNIQUE NOT NULL
);

create table scl.profiles (
    id bigserial PRIMARY KEY,
    user_id varchar(50) NOT NULL,
    name varchar(100),
    surname varchar(100),
    age smallint,
    sex varchar(20),
    address_id bigint
);

create table scl.hobbies (
    id bigserial PRIMARY KEY,
    title varchar(200)
);

create table scl.profile_hobbies (
    profile_id bigint,
    hobby_id bigint
);

create table scl.addresses (
    id bigserial PRIMARY KEY,
    city varchar(100),
    country varchar(100)
);


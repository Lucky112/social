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
    birthdate date,
    sex varchar(20),
    address varchar,
    hobbies varchar
);



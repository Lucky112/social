create table scl.posts (
    id bigserial PRIMARY KEY,
    user_id varchar(50) NOT NULL,
    content varchar,
    created_at timestamp DEFAULT now()
);

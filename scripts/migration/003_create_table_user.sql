-- Write your migrate up statements here
create table avito_banner."user"
(
    user_id         serial primary key,
    user_name  varchar(100) not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

alter table avito_banner."user" owner to postgres;
---- create above / drop below ----
drop table avito_banner."user";
-- Write your migrate up statements here
create table avito_banner."user"
(
    id         serial primary key,
    user_name  varchar(100) not null,
    created_dt timestamp default now() not null,
    updated_dt timestamp default now() not null
);

alter table avito_banner."user" owner to postgres;
---- create above / drop below ----
drop table avito_banner."user";
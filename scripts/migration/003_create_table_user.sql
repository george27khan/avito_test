-- Write your migrate up statements here
create table avito_banner."user"
(
    user_id     serial primary key,
    user_name   varchar(100) not null,
    password    varchar(200) not null,
    is_admin    bool default false not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index idx_user_user_name on avito_banner."user" (user_name);

alter table avito_banner."user" owner to postgres;
---- create above / drop below ----
drop table avito_banner."user";
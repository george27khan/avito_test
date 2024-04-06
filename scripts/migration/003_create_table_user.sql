-- Write your migrate up statements here
create table avito."user"
(
    id         bigint primary key,--serial primary key,
    user_name  varchar(100) not null,
    created_dt timestamp default now() not null
);

alter table avito."user"
    owner to postgres;
---- create above / drop below ----
drop table avito."user";
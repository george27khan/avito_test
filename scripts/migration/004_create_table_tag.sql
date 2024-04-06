-- Write your migrate up statements here
create table avito."tag"
(
    id         bigint primary key,
    created_dt timestamp default now() not null
);

alter table avito."tag"
    owner to postgres;
---- create above / drop below ----
drop table avito."tag";
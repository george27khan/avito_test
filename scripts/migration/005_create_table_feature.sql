-- Write your migrate up statements here
create table avito."feature"
(
    id         bigint primary key,
    created_dt timestamp default now() not null
);

alter table avito."feature"
    owner to postgres;
---- create above / drop below ----
drop table avito."feature";
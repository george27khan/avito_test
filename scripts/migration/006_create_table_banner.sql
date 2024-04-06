-- Write your migrate up statements here
create table avito."banner"
(
    id         bigint primary key,
    json_value json,
    created_dt timestamp default now() not null
);

alter table avito."banner"
    owner to postgres;
---- create above / drop below ----
drop table avito."banner";
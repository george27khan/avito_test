-- Write your migrate up statements here
create table avito_banner."feature"
(
    id         serial primary key,
    created_dt timestamp default now() not null
);

alter table avito_banner."feature"
    owner to postgres;
---- create above / drop below ----
drop table avito_banner."feature";
-- Write your migrate up statements here
create table avito_banner."banner"
(
    id         serial primary key,
    json_value json,
    is_active  boolean not null,
    id_feature bigint not null constraint fk_banner_feature references avito_banner."feature",
    created_dt timestamp default now() not null,
    updated_dt timestamp default now() not null
);

alter table avito_banner."banner"
    owner to postgres;
---- create above / drop below ----
drop table avito_banner."banner";
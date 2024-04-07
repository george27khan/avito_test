-- Write your migrate up statements here
create table avito_banner."banner"
(
    id         serial primary key,
    content    varchar,
    is_active  boolean not null,
    feature_id bigint not null,
    created_dt timestamp default now() not null,
    updated_dt timestamp default now() not null
);
create index idx_fk_banner_feature on avito_banner."tag_feature" (feature_id);
alter table avito_banner."banner"
    owner to postgres;
---- create above / drop below ----
drop table avito_banner."banner";
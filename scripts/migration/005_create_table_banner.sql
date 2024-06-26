-- Write your migrate up statements here
create table avito_banner."banner"
(
    banner_id  serial primary key,
    content    json,
    is_active  boolean not null,
    feature_id bigint not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);
create index idx_fk_banner_feature on avito_banner."tag_feature" (feature_id);
alter table avito_banner."banner"
    owner to postgres;
---- create above / drop below ----
drop table avito_banner."banner";
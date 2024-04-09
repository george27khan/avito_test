-- Write your migrate up statements here
create table avito_banner."tag_feature"
(
    tag_feature_id  serial primary key,
    tag_id          bigint not null,
    feature_id      bigint not null,
    created_at      timestamp default now() not null,
    is_active       boolean default true not null
);
create index idx_tag_feature_tag_id on avito_banner."tag_feature" (tag_id);
create index idx_tag_feature_feature_id on avito_banner."tag_feature" (feature_id);

alter table avito_banner."tag_feature"
    owner to postgres;
---- create above / drop below ----
drop table avito_banner."tag_feature";
-- Write your migrate up statements here
create table avito_banner."tag_feature"
(
    id         serial primary key,
    id_tag     bigint not null constraint fk_tag_feature_1 references avito_banner."tag",
    id_feature bigint not null constraint fk_tag_feature_2 references avito_banner."feature",
    created_dt timestamp default now() not null
);

alter table avito_banner."feature"
    owner to postgres;
---- create above / drop below ----
drop table avito_banner."feature";
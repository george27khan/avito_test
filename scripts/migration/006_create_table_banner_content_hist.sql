-- Write your migrate up statements here
create table avito_banner."banner_content_hist"
(
    id         serial primary key,
    banner_id  bigint not null constraint banner_content_hist_banner_banner_id_fk references avito_banner."banner",
    content    json,
    version    int,
    created_at timestamp default now() not null
);

alter table avito_banner."banner_content_hist"
    owner to postgres;
---- create above / drop below ----
drop table avito_banner."banner_content_hist";
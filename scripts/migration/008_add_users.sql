-- Write your migrate up statements here
INSERT INTO avito_banner."user" (user_name, password, is_admin)
VALUES ('admin1', crypt('admin1', gen_salt('md5')) , true),
       ('admin2', crypt('admin2', gen_salt('md5')) , true),
       ('user1', crypt('user1', gen_salt('md5')) , false),
       ('user2', crypt('user2', gen_salt('md5')) , false),
       ('user3', crypt('user3', gen_salt('md5')) , false),
       ('user4', crypt('user4', gen_salt('md5')) , false);
---- create above / drop below ----
delete from avito_banner."user";
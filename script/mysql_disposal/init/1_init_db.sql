CREATE DATABASE IF NOT EXISTS argus;
USE argus;

CREATE TABLE IF NOT EXESTS articles(
    id int not null primary key,
    title varchar(256) not null,
    create_date date not null,
    update_date date,
    href varchar(512) not null,
    image_href varchar(512) not null,
    private int default 0,
);

CREATE TABLE IF NOT EXISTS categories(
    id int not null primary key,
    title varchar(128) not null
);

CREATE TABLE IF NOT EXISTS article_category(
    article_id int not null,
    category_id int not null,
    constraint fk_article_id_with_category
        foreign key (article_id)
        references articles (id)
        on delete cascade on update cascade,
    constraint fk_category_id_with_article
        foreign key (category_id)
        references categories (id)
        on delete cascade on update cascade
);

CREATE TABLE IF NOT EXISTS authenticate(
    id int not null primary key,
    username varchar(50) not null,
    password varchar(256) not null,
    role varchar(32) not null
);
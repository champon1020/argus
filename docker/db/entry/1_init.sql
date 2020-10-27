CREATE DATABASE IF NOT EXISTS argus;
CREATE DATABASE IF NOT EXISTS argus_test;
USE argus_test;

CREATE TABLE IF NOT EXISTS articles(
    id varchar(128) primary key not null,
    title varchar(256) not null,
    created_date date not null,
    updated_date date not null,
    content longtext character set utf8mb4 not null,
    image_name varchar(512) not null,
    private int default 0,
    index id_index (id)
);

CREATE TABLE IF NOT EXISTS drafts(
    id varchar(128) primary key not null,
    title varchar(256) not null,
    categories varchar(256) not null,
    updated_date date not null,
    content longtext character set utf8mb4 not null,
    image_name varchar(512) not null,
    index id_index (id)
);

CREATE TABLE IF NOT EXISTS categories(
    id varchar(128) primary key not null,
    name varchar(128) not null
);

CREATE TABLE IF NOT EXISTS article_category(
    article_id varchar(128) not null,
    category_id varchar(128) not null,
    constraint fk_article_id_with_category
        foreign key (article_id)
        references articles (id)
        on delete cascade on update cascade,
    constraint fk_category_id_with_article
        foreign key (category_id)
        references categories (id)
        on delete cascade on update cascade
);

USE argus;

CREATE TABLE articles LIKE argus_test.articles;
CREATE TABLE drafts LIKE argus_test.drafts;
CREATE TABLE categories LIKE argus_test.categories;
CREATE TABLE article_category LIKE argus_test.article_category;

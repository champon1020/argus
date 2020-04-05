CREATE DATABASE IF NOT EXISTS argus;
CREATE DATABASE IF NOT EXISTS argus_test;
USE argus_test;

CREATE TABLE IF NOT EXISTS articles(
    id varchar(128) primary key not null,
    sorted_id int not null auto_increment,
    title varchar(256) not null,
    create_date datetime not null,
    update_date datetime not null,
    content longtext character set utf8mb4,
    image_hash varchar(512) not null,
    private int default 0,
    index sorted_index (sorted_id)
);

CREATE TABLE IF NOT EXISTS drafts(
    id varchar(128) primary key not null,
    sorted_id int not null auto_increment,
    title varchar(256) not null,
    categories varchar(256) not null,
    update_date datetime not null,
    content longtext character set utf8mb4,
    image_hash varchar(512) not null,
    index sorted_index (sorted_id)
);

CREATE TABLE IF NOT EXISTS categories(
    id varchar(128) not null primary key,
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
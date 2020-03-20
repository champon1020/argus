USE argus_test;

LOAD DATA LOCAL INFILE
    "/docker/db/csv/articles.csv"
    INTO TABLE articles
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"'
    (id, title, create_date, update_date, content_hash, image_hash, private);

LOAD DATA LOCAL INFILE
    "/docker/db/csv/drafts.csv"
    INTO TABLE drafts
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"'
    (id, title, categories, update_date, content_hash, image_hash);

LOAD DATA LOCAL INFILE
    "/docker/db/csv/categories.csv"
    INTO TABLE categories
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"';

LOAD DATA LOCAL INFILE
    "/docker/db/csv/article_category.csv"
    INTO TABLE article_category
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"';
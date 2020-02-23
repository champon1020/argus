LOAD DATA LOCAL INFILE
    "/docker/db/csv/articles.csv"
    INTO TABLE articles
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"';

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

LOAD DATA LOCAL INFILE
    "/docker/db/csv/authenticate.csv"
    INTO TABLE authenticate
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"';
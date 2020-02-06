LOAD DATA LOCAL INFILE
    "../datas/articles.csv"
    INTO TABLE articles
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"';

LOAD DATA LOCAL INFILE
    "../datas/categories.csv"
    INTO TABLE categories
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"';

LOAD DATA LOCAL INFILE
    "../datas/article_category.csv"
    INTO TABLE article_category
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"';

LOAD DATA LOCAL INFILE
    "../datas/authenticate.csv"
    INTO TABLE authenticate
    CHARACTER SET utf8
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"';
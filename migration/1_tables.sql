CREATE DATABASE IF NOT EXISTS argus;

CREATE TABLE IF NOT EXISTS argus.articles (
  id          VARCHAR(128)  NOT NULL,
  title       VARCHAR(128)  NOT NULL,
  created_at  DATETIME      NOT NULL,
  updated_at  DATETIME      NOT NULL,
  content     LONGTEXT      NOT NULL,
  image_url   VARCHAR(1024) NOT NULL,
  status      INT(1)        NOT NULL  DEFAULT 1,
  PRIMARY KEY (id),
  INDEX id_idx (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS argus.tags (
  article_id  VARCHAR(128)  NOT NULL,
  name        VARCHAR(64)   NOT NULL,
  PRIMARY KEY (article_id, name),
  INDEX article_tag (article_id, name),
  CONSTRAINT fk_article_id_with_tag FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET GLOBAL local_infile=1;

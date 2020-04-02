FROM mysql:5.7

MAINTAINER champon

ADD entry/1_init.sql /docker-entrypoint-initdb.d
ADD entry/3_user.sh /docker-entrypoint-initdb.d
ADD custom.cnf /etc/mysql/conf.d

RUN echo "building ..."
FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD_FILE=/run/secrets/mysql_password

ADD ./create_db.sql /docker-entrypoint-initdb.d
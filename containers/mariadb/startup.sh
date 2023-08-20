#!/bin/bash

# 启动 MariaDB
/usr/local/bin/docker-entrypoint.sh mysqld &

# 等待 MariaDB 启动
until mysqladmin ping -hlocalhost -P3306 --silent; do
    sleep 1
done

# 执行 SQL 脚本
mysql -u root -p"$MYSQL_ROOT_PASSWORD" < /usr/local/bin/startup.sql

# 不要让容器退出，否则 CMD 将无法执行
tail -f /dev/null

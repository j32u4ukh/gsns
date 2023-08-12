# 指定要從 mariadb:latest 映像檔開始。
FROM mariadb:10.6

# 設定環境變量 TZ 的值為 UTC。這個環境變量用於指定時區，這里將其設置為協調世界時（UTC）。
ENV TZ=UTC

# 將位於本地 ./containers/mariadb/my.cnf 路徑下的檔案覆制到鏡像中的 /etc/my.cnf 路徑。這個指令用於將自定義的 mariadb 配置文件覆制到鏡像中，以覆蓋默認的 mariadb 配置。
COPY ./containers/mariadb/my.cnf /etc/mysql/my.cnf

# # 將初始化SQL腳本添加到容器中，這些腳本將在創建數據庫時執行（可選）
# # 如果您有一個名為init.sql的SQL腳本位於同一目錄下，它將被複製到容器內
# COPY ./containers/mariadb/init.sql /docker-entrypoint-initdb.d/

# COPY ./containers/mariadb/startup.sql /usr/local/bin/startup.sql

# # 使用自定義的入口腳本啟動容器
# COPY ./containers/mariadb/startup.sh /usr/local/bin/startup.sh
# RUN chmod +x /usr/local/bin/startup.sh
# CMD ["/usr/local/bin/startup.sh"]

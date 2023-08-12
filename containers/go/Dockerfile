# 指定要從 golang:1.18 映像檔開始。
FROM golang:1.18

# 指定容器的作業目錄為 /app。
WORKDIR /usr/src/app

# 將 go.mod 和 go.sum 檔案複製到容器的 /app 目錄。
COPY ./app /usr/src/app
# COPY ./app/go.sum /usr/src/app

# 安裝 cmd 相依套件
RUN cd /usr/src/app/cmd && go mod download

# 安裝 internal 相依套件
RUN cd /usr/src/app/internal && go mod download
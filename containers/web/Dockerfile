# 使用基於 Node.js 18.4 版本的 Alpine Linux 作為基本映像
FROM node:18.4-alpine

# 設置工作目錄
WORKDIR /usr/src/web

# 將檔案複製到容器的 /usr/src/app 目錄。
COPY ./web /usr/src/web

# 安裝全局的 Vite 工具
RUN yarn global add vite

# 更新和升級軟件包庫
RUN apk update && apk upgrade

# 安裝所需的工具
RUN apk add --no-cache bash git

# 指定內部應用程序使用的端口
EXPOSE 3000



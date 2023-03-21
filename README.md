# gsns
OpenSNS base on Golang, hope everyone can build a SNS and connect to each other simply by this project.

**應用多款獨立開發的套件，開發此`社群網路服務`，開發的同時完善所用到的套件**

* github.com/j32u4ukh/glog (public)
* github.com/j32u4ukh/gos (private)
* github.com/j32u4ukh/gosql (private)

## [glog](https://github.com/j32u4ukh/glog)

參考 Python logging 套件，將 log 區分等級，release 版本直接提升 log 等級，讓開發時的 log 不要印出。

可選擇 log 資訊包含"呼叫的函式名稱"、"所在檔案位置"、"所在行數"等資訊，方便追蹤問題所在。

log 輸出成檔案時，可根據時間或檔案大小，進行分檔。

## gos

伺服器框架，目前支援 tcp 與 http 兩種協定，之後也將支援 WebSocket。

基本上都使用 Golang 原生套件，之後開源後也歡迎發 Pull Request 給我，支援更多種協定。

基於之前開發所需，目前採用單線程，之後也將開發多線程版本，讓使用者根據需求做選擇。

## gosql

使用 Protobuf 來生成 SQL 指令的 ORM 工具，包含讀取 .proto 檔來建立表格。

透過 Where 結構來定義 WHERE 部分的語法，無需自行撰寫語法。

.proto 檔可放在不同伺服器，協定更新或表格更新都很方便。

另有 [j32u4ukh/gorm](https://github.com/j32u4ukh/gorm) 提供類似功能，差別在於使用 Golang 本身的結構，而非 Protobuf。

若要在不同伺服器上使用 j32u4ukh/gorm，則需要再定義一次該結構。

# 其他套件

* [j32u4ukh/cntr](https://github.com/j32u4ukh/cntr): 一些常用容器
* [j32u4ukh/gogettable](https://github.com/j32u4ukh/gogettable): 使自己的 Golang 套件可被別人透過 go get 下載的教學。

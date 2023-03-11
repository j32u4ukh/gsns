# gsns
OpenSNS base on Golang, hope everyone can build a SNS and connect to each other simply by this project.

## Git Hook
---
透過 Git Hook 讓我們在 Chink-in 程式碼（ git push）前、或是在 Local 提交變動（git commit）前，先執行自定義的測試腳本，以確保程式碼是沒有錯誤的。

測試腳本放置於 .git\hook 資料夾下。

# 伺服器位置

範例中，將所有伺服器放在同一台機器，因此利用 port 來區分不同位置：

* Main: 1023
* Dba: 1022
* Account: 1021
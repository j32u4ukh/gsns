## 新增表格流程

1. 定義 proto 檔以及欄位設置
2. dba/define.go 中新增 Table id
3. dba/database.go 的 initTables 中新增新的表格名稱與 Table id
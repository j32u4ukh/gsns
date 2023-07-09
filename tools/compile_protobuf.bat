@ECHO ON
:: 指令格式為: PATH_TO_PROTOCAL_EXECUTE_FILE PATH_TO_PROTOBUF_FILE --xxx_out=PATH_TO_OUTPUT_FOLDER
:: PATH_TO_PROTOCAL_EXECUTE_FILE: 若已將 protoc 編譯程式的路徑加入環境變數，則可直接使用 protoc
CD ..
SET ROOT=%CD%
SET PROTO=%ROOT%/cmd/pb
SET GO_OUT=%ROOT%/internal/pbgo

:: 若不在相同的上層資料夾當中，需透過 -I 來指示 proto 檔案的來源位置
protoc -I %PROTO% %PROTO%/*.proto --go_out=%GO_OUT%
PAUSE
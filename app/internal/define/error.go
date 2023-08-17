package define

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var Error *ErrorConfig

func init() {
	Error = &ErrorConfig{}
	err := LoadConfig("../data/error.yaml", Error)
	if err != nil {
		fmt.Printf("Failed to load error.yaml, err: %+v", err)
		Error.None = 0
	}
}

type ErrorConfig struct {
	None int32
	// 缺少參數
	MissingParameters int32
	// 根據參數找不到用戶
	NotFoundUser int32
	// 無法送出數據
	CannotSendMessage int32
	// 錯誤的作用對象
	InvalidTarget int32
	// 重複的作用實體
	DuplicateEntity int32
	// 無效的 Body 數據
	InvalidBodyData int32
	// 連線身分錯誤
	WrongConnectionIdentity int32
	// 無效寫入數據
	InvalidInsertData int32
	// 無效讀取數據
	InvalidSelectData int32
	// 無效更新數據
	InvalidUpdateData int32
	// 無效刪除數據
	InvalidDeleteData int32
	// 寫入 DB 數據錯誤
	FailedInsertDb int32
	// 讀取 DB 數據錯誤
	FailedSelectDb int32
	// 更新 DB 數據錯誤
	FailedUpdateDb int32
	// 刪除 DB 數據錯誤
	FailedDeleteDb int32
}

func LoadConfig(path string, obj any) error {
	var b []byte
	var err error
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "Failed to read file %s.", path)
	}
	err = yaml.Unmarshal(b, obj)
	if err != nil {
		return errors.Wrapf(err, "讀取 Config 時發生錯誤(path: %s)", path)
	}
	return nil
}

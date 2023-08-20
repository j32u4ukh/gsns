package define

import (
	"fmt"
	"io/ioutil"

	"github.com/j32u4ukh/gos/base/ghttp"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var Error ErrorConfig

func init() {
	Error = ErrorConfig{}
	err := LoadConfig("../data/error.yaml", &Error)
	if err != nil {
		fmt.Printf("Failed to load error.yaml, err: %+v\n", err)
		Error.None = 0
	}
}

type ErrorConfig struct {
	None int32 `yaml:"None"`

	// 400 Bad Request
	BadRequest int32 `yaml:"BadRequest"`
	// 缺少參數
	MissingParameters int32 `yaml:"MissingParameters"`
	// 無效寫入數據
	InvalidInsertData int32 `yaml:"InvalidInsertData"`
	// 無效讀取數據
	InvalidSelectData int32 `yaml:"InvalidSelectData"`
	// 無效更新數據
	InvalidUpdateData int32 `yaml:"InvalidUpdateData"`
	// 無效刪除數據
	InvalidDeleteData int32 `yaml:"InvalidDeleteData"`
	// 無效的 Body 數據
	InvalidBodyData int32 `yaml:"InvalidBodyData"`
	// 無效的作用對象
	InvalidTarget int32 `yaml:"InvalidTarget"`
	// 錯誤的參數
	WrongParameter int32 `yaml:"WrongParameter"`
	// 過多的參數
	TooManyParameter int32 `yaml:"TooManyParameter"`

	// 401 Unauthorized
	Unauthorized int32 `yaml:"Unauthorized"`
	// 連線身分錯誤
	WrongConnectionIdentity int32 `yaml:"WrongConnectionIdentity"`

	// 404 Not Found
	NotFound int32 `yaml:"NotFound"`
	// 根據參數找不到用戶
	NotFoundUser int32 `yaml:"NotFoundUser"`

	// 409 Conflict
	Conflict int32 `yaml:"Conflict"`
	// 重複的實體
	DuplicateEntity int32 `yaml:"DuplicateEntity"`

	// 500 Internal Server Error
	InternalServerError int32 `yaml:"InternalServerError"`
	// 無法送出數據
	CannotSendMessage int32 `yaml:"CannotSendMessage"`
	// 寫入 DB 數據錯誤
	FailedInsertDb int32 `yaml:"FailedInsertDb"`
	// 讀取 DB 數據錯誤
	FailedSelectDb int32 `yaml:"FailedSelectDb"`
	// 更新 DB 數據錯誤
	FailedUpdateDb int32 `yaml:"FailedUpdateDb"`
	// 刪除 DB 數據錯誤
	FailedDeleteDb int32 `yaml:"FailedDeleteDb"`
}

func LoadConfig(path string, obj any) error {
	var bs []byte
	var err error
	bs, err = ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "Failed to read file %s.", path)
	}
	err = yaml.Unmarshal(bs, obj)
	if err != nil {
		return errors.Wrapf(err, "讀取 Config 時發生錯誤(path: %s)", path)
	}
	return nil
}

func GetStatus(code int32) int32 {
	var status int32 = 0
	switch code {
	// 400 Bad Request
	case Error.BadRequest, Error.MissingParameters, Error.InvalidInsertData, Error.InvalidSelectData, Error.InvalidUpdateData, Error.InvalidDeleteData,
		Error.InvalidBodyData, Error.InvalidTarget, Error.WrongParameter, Error.TooManyParameter:
		status = ghttp.StatusBadRequest

	// 401 Unauthorized
	case Error.Unauthorized, Error.WrongConnectionIdentity:
		status = ghttp.StatusUnauthorized

	// 404 Not Found
	case Error.NotFound, Error.NotFoundUser:
		status = ghttp.StatusNotFound

	// 409 Conflict
	case Error.Conflict, Error.DuplicateEntity:
		status = ghttp.StatusConflict

	// 500 Internal Server Error
	case Error.InternalServerError, Error.CannotSendMessage,
		Error.FailedInsertDb, Error.FailedSelectDb, Error.FailedUpdateDb, Error.FailedDeleteDb:
		status = ghttp.StatusInternalServerError
	}
	return status
}

func ErrorMessage(code int32, contents ...any) (int32, int32, string) {
	status := GetStatus(code)
	var em string = ""
	defer func() {
		if err := recover(); err != nil {
			status = ghttp.StatusInternalServerError
			em = "Failed to generate error message"
			fmt.Printf("%s, err: %+v", em, err)
			return
		}
	}()
	switch code {
	// 400 Bad Request
	case Error.BadRequest:
		em = "400 Bad Request"
		if len(contents) > 0 {
			em = fmt.Sprintf("%s | %s", em, contents[0])
		}
	// 缺少參數
	case Error.MissingParameters:
		em = fmt.Sprintf("Not found parameter %+v", contents[0])
	// 無效寫入數據
	case Error.InvalidInsertData:
		em = fmt.Sprintf("Invalid insert data: %+v", contents[0])
	// 無效讀取數據
	case Error.InvalidSelectData:
		em = fmt.Sprintf("Invalid query data: %+v", contents[0])
	// 無效更新數據
	case Error.InvalidUpdateData:
		em = fmt.Sprintf("Invalid update data: %+v", contents[0])
	// 無效刪除數據
	case Error.InvalidDeleteData:
		em = fmt.Sprintf("Invalid delete data: %+v", contents[0])
	// 無效的 Body 數據
	case Error.InvalidBodyData:
		em = "Invalid body data"
	// 無效的作用對象
	case Error.InvalidTarget:
		em = fmt.Sprintf("Target %+v is not valid to %+v", contents[1], contents[0])
	// 錯誤的參數
	case Error.WrongParameter:
		em = fmt.Sprintf("Wrong parameter: %s(%+v)", contents[0].(string), contents[1])
	// 過多的參數
	case Error.TooManyParameter:
		em = fmt.Sprintf("#Parameter except: %d, actually: %d", contents[0].(int), contents[1].(int))

	// 401 Unauthorized
	case Error.Unauthorized:
		em = "401 Unauthorized"
		if len(contents) > 0 {
			em = fmt.Sprintf("%s | %s", em, contents[0])
		}
	// 連線身分錯誤
	case Error.WrongConnectionIdentity:
		em = fmt.Sprintf("Wrong connection identity | Cipher: %s, Identity: %d", contents[0].(string), contents[1].(int32))

	// 404 Not Found
	case Error.NotFound:
		em = "404 Not Found"
		if len(contents) > 0 {
			em = fmt.Sprintf("%s | %s", em, contents[0])
		}
	// 根據參數找不到用戶
	case Error.NotFoundUser:
		nContent := len(contents)
		if nContent == 2 {
			em = fmt.Sprintf("Not found user refer to %s: %+v", contents[0].(string), contents[1])
		} else if nContent == 4 {
			em = fmt.Sprintf("Not found user refer to (%s, %s): (%+v, %+v)", contents[0].(string), contents[2].(string), contents[1], contents[3])
		} else {
			em = fmt.Sprintf("Not found user refer to %+v", contents[0])
		}

	// 409 Conflict
	case Error.Conflict:
		em = "409 Conflict"
		if len(contents) > 0 {
			em = fmt.Sprintf("%s | %s", em, contents[0])
		}
	// 重複的作用實體
	case Error.DuplicateEntity:
		em = fmt.Sprintf("DuplicateEntity | %s", contents[0].(string))

	// 500 Internal Server Error
	case Error.InternalServerError:
		em = "500 Internal Server Error"
		if len(contents) > 0 {
			em = fmt.Sprintf("%s | %s", em, contents[0].(string))
		}
	// 無法送出數據
	case Error.CannotSendMessage:
		em = fmt.Sprintf("Failed to send %s", contents[0].(string))
	// 寫入 DB 數據錯誤
	case Error.FailedInsertDb:
		em = fmt.Sprintf("Failed to insert %s", contents[0].(string))
	// 讀取 DB 數據錯誤
	case Error.FailedSelectDb:
		em = fmt.Sprintf("Failed to query %s", contents[0].(string))
	// 更新 DB 數據錯誤
	case Error.FailedUpdateDb:
		em = fmt.Sprintf("Failed to update %s", contents[0].(string))
	// 刪除 DB 數據錯誤
	case Error.FailedDeleteDb:
		em = fmt.Sprintf("Failed to delete %s", contents[0].(string))
	}
	return status, code, em
}

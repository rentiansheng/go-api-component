package _default

import (
	"github.com/rentiansheng/go-api-component/middleware/errors/code"
	"github.com/rentiansheng/go-api-component/middleware/errors/register"
)

func init() {
	register.Register(langName, asstCodes)
}

var asstCodes = map[int32]string{

	code.JSONDecodeErrCode:   "request body decode error. err: %s",
	code.MapperActionErrCode: "mapper error. action: %s, err: %s",
	code.FileNotFoundErrCode: "file not found. file name: %s",

	code.RawErrWrapErrCode: "raw error wrap: %v",
}

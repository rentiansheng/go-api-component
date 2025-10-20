package context

import (
	"fmt"
	"strings"
	"testing"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/7/7
    @desc:

***************************/

func TestErr_Error(t *testing.T) {
	outputs := []string{
		"go-api-component/middleware/errors/error.go:",
		"go-api-component/middleware/errors/error.go:",
		"go-api-component/middleware/context/err.go:",
		"go-api-component/middleware/context/err_test.go:",
	}
	appErr := (&err{}).Error(-1, fmt.Errorf("app error"))
	callers := appErr.Caller()
	for idx, caller := range callers {
		if idx >= len(outputs) {
			break
		}
		if !strings.HasPrefix(caller, outputs[idx]) {
			t.Errorf("index: %d, atcual: %s, expect: %s", idx, caller, outputs[idx])
		}
	}

}

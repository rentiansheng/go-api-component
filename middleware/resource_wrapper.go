package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	coreContext "github.com/rentiansheng/go-api-component/middleware/context"
	"github.com/rentiansheng/go-api-component/middleware/errors"
)

var (
	apiChecker   sync.Once
	loginChecker CheckLogin
)

const (
	responseHTTHeaderRequestID = "trace-Id"
	RequestId                  = responseHTTHeaderRequestID
)

type Handler func(ctx coreContext.Contexts) errors.Error
type HttpJsonResponse struct {
	Retcode int         `json:"retcode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CheckLogin func(ctx coreContext.Context) errors.Error

func OkResponseExtra(message string, data interface{}, extraRespData map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	for key, val := range extraRespData {
		result[key] = val
	}
	result["retcode"] = 0
	result["message"] = message
	result["data"] = data

	return result

}

func OkResponse(message string, data interface{}) HttpJsonResponse {
	return HttpJsonResponse{
		Retcode: 0,
		Message: message,
		Data:    data,
	}
}

func FailResponse(retcode int, message string, data interface{}) HttpJsonResponse {
	return HttpJsonResponse{
		Retcode: retcode,
		Message: message,
		Data:    data,
	}
}

func Wrapper(h Handler) func(g *gin.Context) {
	return wrapperOptions(func(ctx coreContext.Contexts) errors.Error {
		return h(ctx)
	}, DefaultOption())
}

func wrapperOptions(h Handler, o Option) func(g *gin.Context) {
	return func(g *gin.Context) {
		ctx := coreContext.NewContext(g)

		requestID := ctx.GetRequestID()
		g.Writer.Header().Add(responseHTTHeaderRequestID, requestID)

		// 从panic中恢复
		defer func() {
			if e := recover(); e != nil {
				ctx.Log().Panicf("panic. err: %#v", e)
				// gin  返回 json
				g.JSON(500, e)
				return
			}
		}()

		// 记录请求body
		requestRecords(ctx)
		var err errors.Error
		var data interface{}

		if !o.IsNoLogin() && loginChecker != nil {
			err = loginChecker(ctx)
		}
		// 没有前置错误
		if err == nil {
			err = h(ctx)
			data = ctx.GetData()
			responseRecords(ctx, data, err)
		}
		if err != nil {
			if eerr, ok := err.(errors.Error); ok {
				g.JSON(200, FailResponse(int(eerr.Code()), eerr.Message(), data))
			} else {
				g.JSON(500, FailResponse(-1, err.Error(), data))
			}
		} else {

			if fileName, fileContent, exists := ctx.GetResponseFile(); exists {
				// 返回文件下载
				w := g.Writer
				w.Header().Add("Content-Disposition", "attachment; filename="+fileName)
				w.Header().Add("Content-Type", "application/octet-stream")
				w.Header().Add("Content-Length", fmt.Sprintf("%d", fileContent.Len()))
				_, _ = w.Write(fileContent.Bytes())

			} else if typ, body, exists := ctx.GetRawResponse(); exists {
				w := g.Writer
				if typ != "" {
					w.Header().Add("Content-Type", typ)
				}
				w.Header().Add("Content-Length", fmt.Sprintf("%d", len(body)))
				_, _ = w.Write(body)
			} else {

				extraRespData := ctx.GetExtraResponse()
				if len(extraRespData) > 0 {
					g.JSON(200, OkResponseExtra("", data, extraRespData))
				} else {
					g.JSON(200, OkResponse("", data))
				}
			}
		}
	}
}

func responseRecords(ctx coreContext.Contexts, data interface{}, err error) {

	log := ctx.Log()
	defer func() {
		if r := recover(); r != nil {
			log.PanicJSON("response record. panic info: %s", r)
		}
	}()
	if err != nil {
		if eerr, ok := err.(errors.Error); ok {
			log.ErrorJSON("response record, response error. err code: %d, err message: %s, raw msg: %s, caller: %s", eerr.Code(), eerr.Message(), eerr.RawErrorString(), eerr.Caller())
		} else if eerr, ok := err.(coreContext.BaseErrI); ok {
			log.ErrorJSON("response record, response error. err code: %d, err message: %s", eerr.Code(), eerr.Message())
		} else {
			log.ErrorJSON("response record, response error. err: %s", err)
		}
	} else {
		log.InfoJSON("response record. data: %s", data)
	}
	if ctx.Response() != nil {
		log.InfoJSON("response record http header. header: %s", ctx.Response().Header())
	}
}

// requestRecords 记录请求body < 1m且content_type=application/josn 的http 请求的body
func requestRecords(c coreContext.Context) {
	log := c.Log()

	parentReqId := c.Request().Header.Get(RequestId)

	if c.Request().ContentLength < 1024*1024*1 {

		// Ignore requests smaller than 1MB. This helps prevent delaying
		ct := c.Request().Header.Get("Content-Type")
		if strings.ToLower(ct) == "application/json" || strings.HasPrefix(strings.ToLower(ct), "application/json;") {
			bodyBytes, err := ioutil.ReadAll(c.Request().Body)
			if err != nil {
				log.Errorf("io read request.Body fail, %+v", err)
			}
			// 去除换行符，避免日志换行
			strBody := strings.Replace(string(bodyBytes), "\n", "", -1)
			log.Infof("middleware: record body. method: %s, uri: %s, parent request id: %s, body: %s",
				c.Request().Method, c.Request().RequestURI, parentReqId, strBody)
			c.Request().Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		} else {
			log.Infof("middleware: record body. method: %s, uri: %s, parent request id: %s, body: not support Content-Type=%v",
				c.Request().Method, c.Request().RequestURI, parentReqId, ct)
		}
	} else {
		// body 超过1m 不记录
		log.Infof("middleware: record body. method: %s, uri: %s, parent request id: %s, request body more than 1MB",
			c.Request().Method, c.Request().RequestURI, parentReqId)
	}
}

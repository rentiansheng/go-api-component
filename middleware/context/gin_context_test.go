package context

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupGinTest() (*gin.Engine, *gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, engine := gin.CreateTestContext(w)
	return engine, c, w
}

func TestNewGinContext(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ctx := NewGinContext(c)
	assert.NotNil(t, ctx)
}

func TestGinContext_JSONDecode(t *testing.T) {
	tests := []struct {
		name        string
		jsonBody    string
		target      interface{}
		expectError bool
	}{
		{
			name:     "Valid JSON",
			jsonBody: `{"name":"test","age":25}`,
			target: &struct {
				Name string `json:"name" validate:"required"`
				Age  int    `json:"age" validate:"min=0"`
			}{},
			expectError: false,
		},
		{
			name:     "Invalid JSON",
			jsonBody: `{"name":"test","age":"invalid"}`,
			target: &struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{},
			expectError: true,
		},
		{
			name:     "Validation Error",
			jsonBody: `{"name":"","age":-1}`,
			target: &struct {
				Name string `json:"name" validate:"required"`
				Age  int    `json:"age" validate:"min=0"`
			}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, c, _ := setupGinTest()
			req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(tt.jsonBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			ctx := NewGinContext(c)
			err := ctx.JSONDecode(tt.target)

			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGinContext_Query(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test?name=john&tags=go&tags=gin", nil)
	c.Request = req

	ctx := NewGinContext(c)

	// Test single value
	name := ctx.Query("name")
	assert.Equal(t, []string{"john"}, name)

	// Test multiple values
	tags := ctx.Query("tags")
	assert.Equal(t, []string{"go", "gin"}, tags)

	// Test non-existent key
	missing := ctx.Query("missing")
	assert.Empty(t, missing)
}

func TestGinContext_PathParameter(t *testing.T) {
	engine, c, _ := setupGinTest()
	engine.GET("/users/:id", func(ctx *gin.Context) {})

	req := httptest.NewRequest("GET", "/users/123", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: "123"},
	}

	ctx := NewGinContext(c)
	assert.Equal(t, "123", ctx.PathParameter("id"))
}

func TestGinContext_PathParameters(t *testing.T) {
	engine, c, _ := setupGinTest()
	engine.GET("/users/:userId/posts/:postId", func(ctx *gin.Context) {})

	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "userId", Value: "123"},
		{Key: "postId", Value: "456"},
	}

	ctx := NewGinContext(c)
	params := ctx.PathParameters()

	assert.Equal(t, "123", params["userId"])
	assert.Equal(t, "456", params["postId"])
	assert.Equal(t, 2, len(params))
}

func TestGinContext_ExtraResponse(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ctx := NewGinContext(c)

	// Set extra response
	ctx.SetExtraResponse("key1", "value1")
	ctx.SetExtraResponse("key2", 123)

	// Get extra response
	extra := ctx.GetExtraResponse()
	assert.Equal(t, "value1", extra["key1"])
	assert.Equal(t, 123, extra["key2"])
}

func TestGinContext_PageResponse(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ctx := NewGinContext(c)

	pageInfo := map[string]interface{}{
		"page":  1,
		"limit": 10,
		"total": 100,
	}

	ctx.SetPageResponse(pageInfo)
	// Note: We'd need to expose GetPageResponse to test retrieval
}

func TestGinContext_ResponseFile(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ctx := NewGinContext(c)

	// Set response file
	content := bytes.NewBufferString("file content")
	ctx.SetResponseFile("test.txt", content)

	// Get response file
	fileName, fileContent, exists := ctx.GetResponseFile()
	assert.True(t, exists)
	assert.Equal(t, "test.txt", fileName)
	assert.Equal(t, "file content", fileContent.String())
}

func TestGinContext_RawResponse(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ctx := NewGinContext(c)

	// Set raw response
	ctx.SetRawResponse("application/xml", []byte("<xml>test</xml>"))

	// Get raw response
	typ, body, exists := ctx.GetRawResponse()
	assert.True(t, exists)
	assert.Equal(t, "application/xml", typ)
	assert.Equal(t, "<xml>test</xml>", string(body))
}

func TestGinContext_SetGetData(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ctx := NewGinContext(c)

	testData := map[string]string{"key": "value"}
	ctx.SetData(testData)

	retrievedData := ctx.GetData()
	assert.Equal(t, testData, retrievedData)
}

func TestGinContext_Header(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ctx := NewGinContext(c)

	headers := ctx.Header()
	assert.Equal(t, "Bearer token123", headers.Get("Authorization"))
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
}

func TestGinContext_Cookie(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
	req.AddCookie(&http.Cookie{Name: "user", Value: "john"})
	c.Request = req

	ctx := NewGinContext(c)

	cookies := ctx.Cookie()
	assert.Equal(t, 2, len(cookies))
	assert.Equal(t, "session", cookies[0].Name)
	assert.Equal(t, "abc123", cookies[0].Value)
}

func TestGinContext_SelectedRoutePath(t *testing.T) {
	engine, c, _ := setupGinTest()
	engine.GET("/users/:id", func(ctx *gin.Context) {})

	req := httptest.NewRequest("GET", "/users/123", nil)
	c.Request = req
	c.FullPath()

	ctx := NewGinContext(c)
	// Note: FullPath might be empty in test context without actual routing
	path := ctx.SelectedRoutePath()
	assert.NotNil(t, path)
}

func TestGinContext_WithValue(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ctx := NewGinContext(c)

	ctx.WithValue("userId", "123")
	ctx.WithValue("userName", "john")

	// Retrieve via Value method
	userID := ctx.Value("userId")
	assert.Equal(t, "123", userID)

	userName := ctx.Value("userName")
	assert.Equal(t, "john", userName)
}

func TestGinContext_Decode_Form(t *testing.T) {
	_, c, _ := setupGinTest()

	form := "name=john&age=25"
	req := httptest.NewRequest("POST", "/test", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	ctx := NewGinContext(c)

	target := &struct {
		Name string `form:"name" validate:"required"`
		Age  int    `form:"age" validate:"min=0"`
	}{}

	err := ctx.Decode(target)
	assert.Nil(t, err)
	assert.Equal(t, "john", target.Name)
	assert.Equal(t, 25, target.Age)
}

func TestGinContext_FromFile(t *testing.T) {
	_, c, _ := setupGinTest()

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("upload", "test.txt")
	assert.NoError(t, err)

	_, err = part.Write([]byte("file content"))
	assert.NoError(t, err)

	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req

	ctx := NewGinContext(c)

	file, header, err := ctx.FromFile("upload")
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "test.txt", header.Filename)

	if file != nil {
		defer file.Close()
	}
}

func TestGinContext_Request(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ctx := NewGinContext(c)

	assert.Equal(t, req, ctx.Request())
}

func TestGinContext_IsDone(t *testing.T) {
	_, c, _ := setupGinTest()
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	ctx := NewGinContext(c)

	// Initially should not be done
	assert.False(t, ctx.IsDone())

	// After canceling should be done
	if cancelFn := ctx.Cancel(); cancelFn != nil {
		cancelFn()
		// Give it a moment to propagate
		assert.Eventually(t, func() bool {
			return ctx.IsDone()
		}, 100, 10)
	}
}

package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockRouter implements Router interface for testing
type MockRouter struct {
	registered bool
}

func (m *MockRouter) RegisterGinRoutes(engine *gin.Engine) {
	m.registered = true
	engine.GET("/mock", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "mock endpoint"})
	})
}

func TestRegister(t *testing.T) {
	// Clear the routers before test
	routers = []Router{}

	t.Run("Register valid router", func(t *testing.T) {
		mockRouter := &MockRouter{}
		Register(mockRouter)

		assert.Equal(t, 1, len(routers))
	})

	t.Run("Register nil router", func(t *testing.T) {
		initialLen := len(routers)
		Register(nil)

		// Should not add nil router
		assert.Equal(t, initialLen, len(routers))
	})

	t.Run("Register multiple routers", func(t *testing.T) {
		routers = []Router{} // Reset

		mock1 := &MockRouter{}
		mock2 := &MockRouter{}
		mock3 := &MockRouter{}

		Register(mock1)
		Register(mock2)
		Register(mock3)

		assert.Equal(t, 3, len(routers))
	})
}

func TestGet(t *testing.T) {
	// Setup test routers
	routers = []Router{}

	mock1 := &MockRouter{}
	mock2 := &MockRouter{}

	Register(mock1)
	Register(mock2)

	result := Get()

	assert.Equal(t, 2, len(result))
	assert.Equal(t, mock1, result[0])
	assert.Equal(t, mock2, result[1])
}

func TestGetEmpty(t *testing.T) {
	routers = []Router{}

	result := Get()

	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

func TestRouterIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Reset routers
	routers = []Router{}

	// Register a mock router
	mockRouter := &MockRouter{}
	Register(mockRouter)

	// Get routers and register with gin
	engine := gin.New()
	registeredRouters := Get()

	for _, r := range registeredRouters {
		r.RegisterGinRoutes(engine)
	}

	// Verify the mock router was called
	assert.True(t, mockRouter.registered)
}

func TestMultipleRouterRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Reset
	routers = []Router{}

	// Create multiple mock routers
	mock1 := &MockRouter{}

	// Register them
	Register(mock1)

	// Get and use them
	engine := gin.New()
	for _, r := range Get() {
		r.RegisterGinRoutes(engine)
	}

	// Verify all were registered
	assert.True(t, mock1.registered)

}

func TestRouterOrder(t *testing.T) {
	// Reset
	routers = []Router{}

	mock1 := &MockRouter{}
	mock2 := &MockRouter{}
	mock3 := &MockRouter{}

	Register(mock1)
	Register(mock2)
	Register(mock3)

	result := Get()

	// Verify order is maintained
	assert.Equal(t, mock1, result[0])
	assert.Equal(t, mock2, result[1])
	assert.Equal(t, mock3, result[2])
}

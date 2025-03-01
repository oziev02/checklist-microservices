package handlers

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)
	r.POST("/2fa/setup", setup2FAHandler)
	r.POST("/2fa/verify", verify2FAHandler)

	auth := r.Group("/", authMiddleware)
	{
		// Tasks
		auth.POST("/create", createTaskHandler)
		auth.GET("/list", listTasksHandler)
		auth.DELETE("/delete", deleteTaskHandler)
		auth.PUT("/done", markTaskDoneHandler)

		// Profile
		auth.GET("/profile", getProfileHandler)
		auth.PUT("/profile", updateProfileHandler)
	}
}

// authMiddleware проверяет JWT-токен (пока это заглушка)
func authMiddleware(c *gin.Context) {
	// Пока просто пропускаем запрос
	c.Next()
}

func registerHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Register endpoint placeholder"})
}

func loginHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Login endpoint placeholder"})
}

// настраивает 2FA для пользователя
func setup2FAHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "2FA setup endpoint placeholder"})
}

// проверяет 2FA-код
func verify2FAHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "2FA verify endpoint placeholder"})
}

package handlers

import "github.com/gin-gonic/gin"

func getProfileHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get profile endpoint placeholder"})
}

func updateProfileHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update profile endpoint placeholder"})
}

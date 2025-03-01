package handlers

import "github.com/gin-gonic/gin"

func createTaskHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Create task endpoint placeholder"})
}

func listTasksHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "List tasks endpoint placeholder"})
}

func deleteTaskHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Task ID is required"})
		return
	}
	c.JSON(204, gin.H{"message": "Delete task endpoint placeholder", "id": id})
}

func markTaskDoneHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Task ID is required"})
		return
	}
	c.JSON(200, gin.H{"message": "Mark task done endpoint placeholder", "id": id})
}

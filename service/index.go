// @Author TrandLiu
// @Date 2024/12/12 19:34:00
// @Desc
package service

import "github.com/gin-gonic/gin"

// GetIndex
// @Tags        首页
// @Success      200  {string}   welcome
// @Router       /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome!!",
	})
}

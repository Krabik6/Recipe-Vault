package handler

import "github.com/gin-gonic/gin"

func (h *Handler) signUp(c *gin.Context) {
	c.JSON(200, gin.H{"handler": "signUp"})
}

func (h *Handler) signIn(c *gin.Context) {
	c.JSON(200, gin.H{"handler": "signIn"})
}

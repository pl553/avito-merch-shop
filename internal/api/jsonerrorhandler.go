package api

import (
	"github.com/gin-gonic/gin"
)

func JSONErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 && !c.Writer.Written() {
		errMessages := ""
		for i, err := range c.Errors {
			errMessages += err.Error()
			if i != len(c.Errors)-1 {
				errMessages += "; "
			}
		}
		status := c.Writer.Status()

		c.JSON(status, ErrorResponse{
			Errors: &errMessages,
		})
	}
}

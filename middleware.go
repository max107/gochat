package main

import (
	"github.com/gin-gonic/gin"
)

func DummyMiddleware() gin.HandlerFunc {
	// Do some initialization logic here
	// Foo()
	return func(c *gin.Context) {
		c.Next()
	}
}

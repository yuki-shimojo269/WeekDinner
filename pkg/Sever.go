package week_menue

import(
  "net/http"

  "github.com/gin-gonic/gin"
)

func Server() {
    engine:= gin.Default()
    engine.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "hello world",
        })
    })
    engine.Run(":3000")
}

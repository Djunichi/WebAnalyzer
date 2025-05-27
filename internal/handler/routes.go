package handler

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

const routePrefix = "api/v1"

func (h *httpHandler) addRoutes() {
	h.router.Use(UseRecoverMiddleware())
	h.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	h.router.GET("ping", h.healthCheck())
	h.router.GET("version", h.version())

	h.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	webPageRouter := h.router.Group(pathWithAction("web-pages"))
	webPageRouter.POST("analyze", h.analyzePage())

	analysisRouter := h.router.Group(pathWithAction("analyses"))
	analysisRouter.GET("/all", h.getAll())
	analysisRouter.GET("/by-id", h.getById())
}

func UseRecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {

					stackInfo := debug.Stack()
					log.Println("Recovered from panic", err, string(stackInfo))

					// Respond with an internal server error
					c.AbortWithStatus(http.StatusInternalServerError)
				} else {
					// If the recovered value is not an error, log it as a message
					log.Println("Recovered from panic with non-error value", r)

					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()

		// Continue with the next middleware or route
		c.Next()
	}
}

func pathWithAction(action string) string {
	return fmt.Sprintf("%s/%s", routePrefix, action)
}

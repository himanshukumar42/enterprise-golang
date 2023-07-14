package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	uuid "github.com/google/uuid"
	"github.com/himanshukumar42/enterprise/db"
	"github.com/himanshukumar42/enterprise/forms"
	"github.com/joho/godotenv"
	"storj.io/storj/multinode/console/controllers"
)

// CORSMiddleware ...
// CORS (Cross-Origin Resource Sharing)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// RequestIDMiddleware
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}

func main() {
	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load .env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	// start the default gin server
	r := gin.Default()

	// custom form validator
	binding.Validator = new(forms.DefaultValidator)

	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	//Start PostgreSQL database
	// Example: db.GetDB()
	db.Init()

	// Start Redis on database 1 - It's used to store the JWT but you can use it for anything else
	// Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()
	db.InitRedis(1)

	v1 := r.Group("/v1")
	{
		/*** START USER **/
		user := new(controllers.UserController)
	}
}

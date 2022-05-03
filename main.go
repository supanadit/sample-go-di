package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func RegisterEnvironmentVariable() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func NewGin(lc fx.Lifecycle) *gin.Engine {
	port := os.Getenv("APP_PORT")
	g := gin.Default()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        g,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			fmt.Printf("Application started at port %s\n", port)
			return s.ListenAndServe()
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
	return g
}

func RegisterDummyRouter(g *gin.Engine) {
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func main() {
	// Create Default Options
	options := []fx.Option{
		fx.Provide(
			NewGin,
		),
		fx.Invoke(RegisterEnvironmentVariable, RegisterDummyRouter),
	}

	// Create Inversion of Control
	app := fx.New(options...)

	// Start Context
	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

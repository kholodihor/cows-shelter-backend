package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/handler"
)

func main(){
	log.Println("Starting server...")

	router := gin.Default()

	config.Connect()

	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Allow your frontend URL
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	handler.RouterHandler(&handler.Config{
		R:router,
	})

	router.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")

	if port == ""{
		port = "8080"
	}

	srv := &http.Server{
		Addr: ":" + port,
		Handler: router,
	}

	go func(){
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port: %v\n", srv.Addr)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down the server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n",err)
	}
}

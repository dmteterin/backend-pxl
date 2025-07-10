package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app-server/internal/handlers"
	"app-server/internal/middleware"
	"app-server/internal/websocket"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	hub    *websocket.Hub
}

func New() *Server {
	router := gin.New()
	hub := websocket.NewHub()

	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	router.GET("/", handlers.Home)
	router.GET("/health", handlers.Health)

	wsHandler := handlers.NewWebSocketHandler(hub)
	router.GET("/ws", wsHandler.HandleWebSocket)
	router.GET("/ws/stats", wsHandler.GetStats)

	go hub.Run()

	return &Server{
		router: router,
		hub:    hub,
	}
}

func (s *Server) Run(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
	return nil
}

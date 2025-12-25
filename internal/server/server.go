package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"site/internal/db"
	"site/internal/models"

	"github.com/gorilla/websocket"
)

const (
	devUserID    = "dev:ephemeral"
	devUserName  = "Test User"
	devUserEmail = "matt@localhost"
)

type Server struct {
	config    Config
	db        *db.DB
	wsClients map[*websocket.Conn]bool
	wsLock    sync.Mutex
	upgrader  websocket.Upgrader
	devUser   *models.User
}

func Run(cfg Config) error {
	s := &Server{
		config:    cfg,
		wsClients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	database, err := db.New("data/sqlite.db")
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	s.db = database
	defer s.db.Close()

	if cfg.DevMode {
		// Create ephemeral dev user
		s.devUser = &models.User{
			ID:        devUserID,
			Email:     devUserEmail,
			Name:      devUserName,
			AvatarURL: "https://avatar.vercel.sh/dev-user.svg?text=MB",
			CreatedAt: time.Now(),
		}
		if err := s.db.CreateOrUpdateUser(s.devUser); err != nil {
			log.Printf("Failed to create dev user: %v", err)
		} else {
			log.Println("Ephemeral dev user created (comments will be cleaned up on shutdown)")
		}

		log.Println("Starting initial build...")
		start := time.Now()
		if err := s.rebuild(); err != nil {
			log.Printf("Initial build failed: %v", err)
		} else {
			log.Printf("Initial build complete in %v", time.Since(start))
		}
	}

	mux := s.setupRoutes()

	if cfg.DevMode {
		go s.watchFiles()
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		// Clean up ephemeral dev user data
		if cfg.DevMode && s.devUser != nil {
			log.Println("Cleaning up ephemeral dev user data...")
			if err := s.db.CleanupUserData(devUserID); err != nil {
				log.Printf("Failed to cleanup dev user data: %v", err)
			} else {
				log.Println("Dev user data cleaned up successfully")
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Printf("Server starting on http://localhost:%d", cfg.Port)
	if cfg.DevMode {
		log.Println("Development mode enabled with hot reload")
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not listen on %d: %w", cfg.Port, err)
	}

	<-done
	log.Println("Server stopped")
	return nil
}

func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/reactions", s.handleReactions)
	mux.HandleFunc("/api/reactions/user", s.handleUserReactions)
	mux.HandleFunc("/api/me", s.handleMe)
	mux.HandleFunc("/api/search", s.handleSearch)
	mux.HandleFunc("/api/comments", s.handleComments)
	mux.HandleFunc("/api/comments/", s.handleComment)

	mux.HandleFunc("/auth/google", s.handleGoogleAuth)
	mux.HandleFunc("/auth/google/callback", s.handleGoogleCallback)
	mux.HandleFunc("/auth/github", s.handleGitHubAuth)
	mux.HandleFunc("/auth/github/callback", s.handleGitHubCallback)
	mux.HandleFunc("/auth/logout", s.handleLogout)
	mux.HandleFunc("/api/auth/providers", s.handleAuthProviders)

	mux.HandleFunc("/github", s.handleGitHubRedirect)
	mux.HandleFunc("/linkedin", s.handleLinkedInRedirect)
	mux.HandleFunc("/email", s.handleEmailRedirect)

	if s.config.DevMode {
		mux.HandleFunc("/ws", s.handleWebSocket)
	}

	mux.HandleFunc("/", s.handleStatic)

	return mux
}

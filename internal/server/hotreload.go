package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"site/internal/build"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	s.wsLock.Lock()
	s.wsClients[conn] = true
	s.wsLock.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			s.wsLock.Lock()
			delete(s.wsClients, conn)
			s.wsLock.Unlock()
			conn.Close()
			break
		}
	}
}

func (s *Server) notifyClients() {
	s.wsLock.Lock()
	defer s.wsLock.Unlock()

	for conn := range s.wsClients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
			conn.Close()
			delete(s.wsClients, conn)
		}
	}
}

func (s *Server) watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create file watcher: %v", err)
		return
	}
	defer watcher.Close()

	dirs := []string{s.config.ContentDir, s.config.TemplateDir, s.config.StaticDir}
	for _, dir := range dirs {
		if err := addDirRecursive(watcher, dir); err != nil {
			log.Printf("Failed to watch %s: %v", dir, err)
		}
	}

	var debounceTimer *time.Timer
	debounceDelay := 100 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(debounceDelay, func() {
					log.Printf("File changed: %s", event.Name)
					start := time.Now()
					if err := s.rebuild(); err != nil {
						log.Printf("Rebuild failed: %v", err)
					} else {
						elapsed := time.Since(start)
						log.Printf("Rebuild complete in %v", elapsed)
						s.notifyClients()
					}
				})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func addDirRecursive(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}

func (s *Server) rebuild() error {
	cfg := build.Config{
		ContentDir:  s.config.ContentDir,
		OutputDir:   s.config.OutputDir,
		StaticDir:   s.config.StaticDir,
		TemplateDir: s.config.TemplateDir,
		BaseURL:     s.config.BaseURL,
		DevMode:     s.config.DevMode,
		DB:          s.db,
	}
	return build.Build(cfg)
}

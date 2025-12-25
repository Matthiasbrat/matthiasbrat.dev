package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	filePath := filepath.Join(s.config.OutputDir, path)

	if !strings.Contains(filepath.Base(path), ".") {
		indexPath := filepath.Join(filePath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			filePath = indexPath
		}
	}

	info, err := os.Stat(filePath)
	if err != nil {
		htmlPath := filePath + ".html"
		if _, err := os.Stat(htmlPath); err == nil {
			filePath = htmlPath
		} else {
			http.NotFound(w, r)
			return
		}
	}

	if info != nil && info.IsDir() {
		filePath = filepath.Join(filePath, "index.html")
	}

	setContentType(w, filePath)
	http.ServeFile(w, r, filePath)
}

func setContentType(w http.ResponseWriter, filePath string) {
	ext := filepath.Ext(filePath)
	contentTypes := map[string]string{
		".html": "text/html; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".js":   "application/javascript; charset=utf-8",
		".json": "application/json; charset=utf-8",
		".xml":  "application/xml; charset=utf-8",
		".pdf":  "application/pdf",
	}
	if ct, ok := contentTypes[ext]; ok {
		w.Header().Set("Content-Type", ct)
	}
}

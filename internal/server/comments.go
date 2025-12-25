package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// commentMarkdown is a simple markdown renderer for comments
// It uses GFM extensions but does NOT allow unsafe HTML
var commentMarkdown = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
		// NOTE: We do NOT use html.WithUnsafe() to prevent XSS
	),
)

// renderCommentMarkdown renders markdown content to HTML
func renderCommentMarkdown(content string) string {
	var buf bytes.Buffer
	if err := commentMarkdown.Convert([]byte(content), &buf); err != nil {
		return content // Return original content on error
	}
	return buf.String()
}

// handleComments handles GET and POST for comments
func (s *Server) handleComments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getComments(w, r)
	case http.MethodPost:
		s.createComment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleComment handles PUT and DELETE for a single comment
func (s *Server) handleComment(w http.ResponseWriter, r *http.Request) {
	// Extract comment ID from path: /api/comments/123
	path := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		s.updateComment(w, r, id)
	case http.MethodDelete:
		s.deleteComment(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type commentResponse struct {
	ID         int64  `json:"id"`
	Content    string `json:"content"`
	ContentHTML string `json:"contentHtml"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	UserID     string `json:"userId"`
	UserName   string `json:"userName"`
	UserAvatar string `json:"userAvatar"`
}

// getComments returns all comments for a post
func (s *Server) getComments(w http.ResponseWriter, r *http.Request) {
	postSlug := r.URL.Query().Get("post")
	if postSlug == "" {
		http.Error(w, "Missing post parameter", http.StatusBadRequest)
		return
	}

	comments, err := s.db.GetComments(postSlug)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	response := make([]commentResponse, 0, len(comments))
	for _, c := range comments {
		response = append(response, commentResponse{
			ID:          c.ID,
			Content:     c.Content,
			ContentHTML: renderCommentMarkdown(c.Content),
			CreatedAt:   c.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
			UserID:      c.UserID,
			UserName:    c.UserName,
			UserAvatar:  c.UserAvatar,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// createComment creates a new comment
func (s *Server) createComment(w http.ResponseWriter, r *http.Request) {
	user := s.getSessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Post    string `json:"post"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Trim and validate content
	req.Content = strings.TrimSpace(req.Content)
	if req.Post == "" || req.Content == "" {
		http.Error(w, "Missing post or content", http.StatusBadRequest)
		return
	}

	// Limit content length (10KB max)
	if len(req.Content) > 10240 {
		http.Error(w, "Content too long", http.StatusBadRequest)
		return
	}

	comment, err := s.db.CreateComment(user.ID, req.Post, req.Content)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(commentResponse{
		ID:          comment.ID,
		Content:     comment.Content,
		ContentHTML: renderCommentMarkdown(comment.Content),
		CreatedAt:   comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   comment.UpdatedAt.Format(time.RFC3339),
		UserID:      user.ID,
		UserName:    user.Name,
		UserAvatar:  user.AvatarURL,
	})
}

// updateComment updates an existing comment
func (s *Server) updateComment(w http.ResponseWriter, r *http.Request, id int64) {
	user := s.getSessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	if len(req.Content) > 10240 {
		http.Error(w, "Content too long", http.StatusBadRequest)
		return
	}

	err := s.db.UpdateComment(id, user.ID, req.Content)
	if err == sql.ErrNoRows {
		http.Error(w, "Comment not found or not owned by user", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	// Fetch the updated comment to return
	comment, err := s.db.GetComment(id)
	if err != nil || comment == nil {
		http.Error(w, "Failed to get updated comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commentResponse{
		ID:          comment.ID,
		Content:     comment.Content,
		ContentHTML: renderCommentMarkdown(comment.Content),
		CreatedAt:   comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   comment.UpdatedAt.Format(time.RFC3339),
		UserID:      user.ID,
		UserName:    user.Name,
		UserAvatar:  user.AvatarURL,
	})
}

// deleteComment deletes a comment
func (s *Server) deleteComment(w http.ResponseWriter, r *http.Request, id int64) {
	user := s.getSessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := s.db.DeleteComment(id, user.ID)
	if err == sql.ErrNoRows {
		http.Error(w, "Comment not found or not owned by user", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package server

import (
	"encoding/json"
	"net/http"

	"site/internal/models"
)

// handleReactions handles GET and POST for reactions
func (s *Server) handleReactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getReactions(w, r)
	case http.MethodPost:
		s.postReaction(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getReactions returns reaction counts for a post
func (s *Server) getReactions(w http.ResponseWriter, r *http.Request) {
	postSlug := r.URL.Query().Get("post")
	if postSlug == "" {
		http.Error(w, "Missing post parameter", http.StatusBadRequest)
		return
	}

	counts, err := s.db.GetReactionCounts(postSlug)
	if err != nil {
		http.Error(w, "Failed to get reactions", http.StatusInternalServerError)
		return
	}

	if counts == nil {
		counts = []models.ReactionCount{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}

// postReaction adds or removes a reaction
func (s *Server) postReaction(w http.ResponseWriter, r *http.Request) {
	user := s.getSessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Post  string `json:"post"`
		Emoji string `json:"emoji"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Post == "" || req.Emoji == "" {
		http.Error(w, "Missing post or emoji", http.StatusBadRequest)
		return
	}

	if !models.IsValidEmoji(req.Emoji) {
		http.Error(w, "Invalid emoji", http.StatusBadRequest)
		return
	}

	added, err := s.db.AddReaction(user.ID, req.Post, req.Emoji)
	if err != nil {
		http.Error(w, "Failed to toggle reaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"added": added})
}

// handleUserReactions returns the current user's reactions for a post
func (s *Server) handleUserReactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := s.getSessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postSlug := r.URL.Query().Get("post")
	if postSlug == "" {
		http.Error(w, "Missing post parameter", http.StatusBadRequest)
		return
	}

	emojis, err := s.db.GetUserReactions(user.ID, postSlug)
	if err != nil {
		http.Error(w, "Failed to get user reactions", http.StatusInternalServerError)
		return
	}

	if emojis == nil {
		emojis = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emojis)
}

// handleMe returns the current user's info
func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := s.getSessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":     user.ID,
		"email":  user.Email,
		"name":   user.Name,
		"avatar": user.AvatarURL,
	})
}

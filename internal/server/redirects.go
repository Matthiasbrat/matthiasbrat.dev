package server

import "net/http"

func (s *Server) handleGitHubRedirect(w http.ResponseWriter, r *http.Request) {
	if s.config.Profile.GitHub == "" {
		http.Error(w, "GitHub profile not configured", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, s.config.Profile.GitHub, http.StatusMovedPermanently)
}

func (s *Server) handleLinkedInRedirect(w http.ResponseWriter, r *http.Request) {
	if s.config.Profile.LinkedIn == "" {
		http.Error(w, "LinkedIn profile not configured", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, s.config.Profile.LinkedIn, http.StatusMovedPermanently)
}

func (s *Server) handleEmailRedirect(w http.ResponseWriter, r *http.Request) {
	if s.config.Profile.Email == "" {
		http.Error(w, "Email not configured", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "mailto:"+s.config.Profile.Email, http.StatusMovedPermanently)
}

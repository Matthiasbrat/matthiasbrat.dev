package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"site/internal/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

const (
	sessionCookieName = "session"
	sessionDuration   = 30 * 24 * time.Hour // 30 days
)

// getOAuthConfig returns the Google OAuth2 config
func (s *Server) getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  s.config.BaseURL + "/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// handleGoogleAuth initiates the Google OAuth flow
func (s *Server) handleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	cfg := s.getOAuthConfig()
	if cfg.ClientID == "" {
		http.Error(w, "OAuth not configured", http.StatusServiceUnavailable)
		return
	}

	// Store redirect URL in state (simple approach)
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/"
	}

	state := base64.URLEncoding.EncodeToString([]byte(redirect))
	url := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback handles the OAuth callback
func (s *Server) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	cfg := s.getOAuthConfig()

	// Get the redirect URL from state
	state := r.URL.Query().Get("state")
	redirectBytes, _ := base64.URLEncoding.DecodeString(state)
	redirect := string(redirectBytes)
	if redirect == "" {
		redirect = "/"
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := cfg.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info
	client := cfg.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Create or update user
	user := &models.User{
		ID:        "google:" + userInfo.ID,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: userInfo.Picture,
		CreatedAt: time.Now(),
	}
	if err := s.db.CreateOrUpdateUser(user); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Create session
	sessionToken, err := generateToken(32)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(sessionDuration)
	if err := s.db.CreateSession(sessionToken, user.ID, expiresAt); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

// handleLogout clears the session
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		s.db.DeleteSession(cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/"
	}
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

// getSessionUser returns the user from the session cookie
func (s *Server) getSessionUser(r *http.Request) *models.User {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		// In dev mode, return the ephemeral dev user if no session
		if s.config.DevMode && s.devUser != nil {
			return s.devUser
		}
		return nil
	}

	userID, err := s.db.GetSession(cookie.Value)
	if err != nil || userID == "" {
		// In dev mode, return the ephemeral dev user if session invalid
		if s.config.DevMode && s.devUser != nil {
			return s.devUser
		}
		return nil
	}

	user, err := s.db.GetUser(userID)
	if err != nil {
		return nil
	}

	return user
}

// generateToken creates a random token
func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// getGitHubOAuthConfig returns the GitHub OAuth2 config
func (s *Server) getGitHubOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  s.config.BaseURL + "/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

// handleGitHubAuth initiates the GitHub OAuth flow
func (s *Server) handleGitHubAuth(w http.ResponseWriter, r *http.Request) {
	cfg := s.getGitHubOAuthConfig()
	if cfg.ClientID == "" {
		http.Error(w, "GitHub OAuth not configured", http.StatusServiceUnavailable)
		return
	}

	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/"
	}

	state := base64.URLEncoding.EncodeToString([]byte(redirect))
	url := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGitHubCallback handles the GitHub OAuth callback
func (s *Server) handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	cfg := s.getGitHubOAuthConfig()

	state := r.URL.Query().Get("state")
	redirectBytes, _ := base64.URLEncoding.DecodeString(state)
	redirect := string(redirectBytes)
	if redirect == "" {
		redirect = "/"
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	token, err := cfg.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := cfg.Client(r.Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// GitHub may not return email in main response, fetch from emails endpoint
	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if json.NewDecoder(emailResp.Body).Decode(&emails) == nil {
				for _, e := range emails {
					if e.Primary {
						userInfo.Email = e.Email
						break
					}
				}
			}
		}
	}

	// Use login as name if name is empty
	displayName := userInfo.Name
	if displayName == "" {
		displayName = userInfo.Login
	}

	user := &models.User{
		ID:        "github:" + strconv.FormatInt(userInfo.ID, 10),
		Email:     userInfo.Email,
		Name:      displayName,
		AvatarURL: userInfo.AvatarURL,
		CreatedAt: time.Now(),
	}
	if err := s.db.CreateOrUpdateUser(user); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	sessionToken, err := generateToken(32)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(sessionDuration)
	if err := s.db.CreateSession(sessionToken, user.ID, expiresAt); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

// handleAuthProviders returns the list of configured OAuth providers
func (s *Server) handleAuthProviders(w http.ResponseWriter, r *http.Request) {
	type provider struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var providers []provider

	if os.Getenv("GOOGLE_CLIENT_ID") != "" {
		providers = append(providers, provider{ID: "google", Name: "Google"})
	}
	if os.Getenv("GITHUB_CLIENT_ID") != "" {
		providers = append(providers, provider{ID: "github", Name: "GitHub"})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]provider{"providers": providers})
}

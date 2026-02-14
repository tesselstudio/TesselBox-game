package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// Provider represents an OAuth provider
type Provider string

const (
	ProviderGoogle Provider = "google"
	ProviderGitHub Provider = "github"
)

// User represents an authenticated user
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Token    string   `json:"token"`
	Provider Provider `json:"provider"`
}

// Session represents a user session
type Session struct {
	User      *User
	CreatedAt time.Time
	ExpiresAt time.Time
}

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	GoogleClientID     string        `json:"google_client_id"`
	GoogleClientSecret string        `json:"google_client_secret"`
	GoogleCallbackURL  string        `json:"google_callback_url"`

	GitHubClientID     string        `json:"github_client_id"`
	GitHubClientSecret string        `json:"github_client_secret"`
	GitHubCallbackURL  string        `json:"github_callback_url"`

	SessionDuration time.Duration `json:"session_duration"`
}

// AuthServiceConfig allows optional providers and graceful fallback
// If both providers are nil/unconfigured, service operates in guest mode.
type AuthServiceConfig struct {
	OAuth *OAuthConfig
}

// AuthService manages OAuth authentication
type AuthService struct {
	config      *OAuthConfig
	googleConf  *oauth2.Config
	githubConf  *oauth2.Config
	sessions    map[string]*Session
	currentUser *User
	server      *http.Server
	authChan    chan *User
}

// NewAuthService creates a new OAuth authentication service (legacy, panics on error). Deprecated: use NewAuthServiceSafe.
func NewAuthService(config *OAuthConfig) *AuthService {
	auth, err := NewAuthServiceSafe(&AuthServiceConfig{OAuth: config})
	if err != nil {
		log.Fatal(err)
	}
	return auth
}

// NewAuthServiceSafe creates a new authentication service and returns errors instead of exiting.
// Supports optional providers and guest mode when none are configured.
func NewAuthServiceSafe(cfg *AuthServiceConfig) (*AuthService, error) {
	if cfg == nil || cfg.OAuth == nil {
		// Guest mode: no providers configured
		auth := &AuthService{
			config:   &OAuthConfig{SessionDuration: 24 * time.Hour},
			sessions: make(map[string]*Session),
			authChan: make(chan *User, 1),
		}
		return auth, nil
	}

	config := cfg.OAuth

	auth := &AuthService{
		config:   config,
		sessions: make(map[string]*Session),
		authChan: make(chan *User, 1),
	}

	// Configure Google OAuth if fully specified
	if config.GoogleClientID != "" && config.GoogleClientSecret != "" && config.GoogleCallbackURL != "" {
		auth.googleConf = &oauth2.Config{
			ClientID:     config.GoogleClientID,
			ClientSecret: config.GoogleClientSecret,
			RedirectURL:  config.GoogleCallbackURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	} else if config.GoogleClientID != "" || config.GoogleClientSecret != "" || config.GoogleCallbackURL != "" {
		return nil, errors.New("incomplete Google OAuth configuration: require client_id, client_secret, and callback_url or none")
	}

	// Configure GitHub OAuth if fully specified
	if config.GitHubClientID != "" && config.GitHubClientSecret != "" && config.GitHubCallbackURL != "" {
		auth.githubConf = &oauth2.Config{
			ClientID:     config.GitHubClientID,
			ClientSecret: config.GitHubClientSecret,
			RedirectURL:  config.GitHubCallbackURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		}
	} else if config.GitHubClientID != "" || config.GitHubClientSecret != "" || config.GitHubCallbackURL != "" {
		return nil, errors.New("incomplete GitHub OAuth configuration: require client_id, client_secret, and callback_url or none")
	}

	if auth.googleConf != nil {
		log.Printf("Google callback URL: %s", auth.googleConf.RedirectURL)
	}
	if auth.githubConf != nil {
		log.Printf("GitHub callback URL: %s", auth.githubConf.RedirectURL)
	}

	return auth, nil
}

// LoadConfigFromEnv loads OAuth configuration from environment variables (preferred)
func LoadConfigFromEnv() *OAuthConfig {
	config := &OAuthConfig{
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleCallbackURL:  os.Getenv("GOOGLE_CALLBACK_URL"),

		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GitHubCallbackURL:  os.Getenv("GITHUB_CALLBACK_URL"),

		SessionDuration: 24 * time.Hour,
	}

	// You can override session duration via env if needed
	if durStr := os.Getenv("SESSION_DURATION"); durStr != "" {
		if dur, err := time.ParseDuration(durStr); err == nil {
			config.SessionDuration = dur
		}
	}

	return config
}

// LoadConfigFromFile loads OAuth configuration from a JSON file (fallback)
func LoadConfigFromFile(filePath string) (*OAuthConfig, error) {
	type RawOAuthConfig struct {
		GoogleClientID     string `json:"google_client_id"`
		GoogleClientSecret string `json:"google_client_secret"`
		GoogleCallbackURL  string `json:"google_callback_url"`
		GitHubClientID     string `json:"github_client_id"`
		GitHubClientSecret string `json:"github_client_secret"`
		GitHubCallbackURL  string `json:"github_callback_url"`
		SessionDuration    string `json:"session_duration"`
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var rawConfig RawOAuthConfig
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Parse session duration
	var sessionDuration time.Duration = 24 * time.Hour
	if rawConfig.SessionDuration != "" {
		sessionDuration, err = time.ParseDuration(rawConfig.SessionDuration)
		if err != nil {
			return nil, fmt.Errorf("invalid session duration: %w", err)
		}
	}

	config := &OAuthConfig{
		GoogleClientID:     rawConfig.GoogleClientID,
		GoogleClientSecret: rawConfig.GoogleClientSecret,
		GoogleCallbackURL:  rawConfig.GoogleCallbackURL,
		GitHubClientID:     rawConfig.GitHubClientID,
		GitHubClientSecret: rawConfig.GitHubClientSecret,
		GitHubCallbackURL:  rawConfig.GitHubCallbackURL,
		SessionDuration:    sessionDuration,
	}

	return config, nil
}

// StartServer starts the local OAuth callback server
func (auth *AuthService) StartServer() error {
	if auth.server != nil {
		return fmt.Errorf("OAuth server already running")
	}

	mux := http.NewServeMux()
	if auth.googleConf != nil {
		mux.HandleFunc("/auth/google/callback", auth.handleGoogleCallback)
	}
	if auth.githubConf != nil {
		mux.HandleFunc("/auth/github/callback", auth.handleGitHubCallback)
	}

	auth.server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("OAuth callback server listening on :8080")
		if err := auth.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("OAuth server error: %v", err)
		}
	}()

	return nil
}

// StopServer stops the OAuth callback server
func (auth *AuthService) StopServer() error {
	if auth.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return auth.server.Shutdown(ctx)
	}
	return nil
}

// GetAuthURL returns the OAuth authorization URL for the specified provider
func (auth *AuthService) GetAuthURL(provider Provider) string {
	var url string
	switch provider {
	case ProviderGoogle:
		if auth.googleConf == nil {
			return ""
		}
		url = auth.googleConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	case ProviderGitHub:
		if auth.githubConf == nil {
			return ""
		}
		url = auth.githubConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	default:
		return ""
	}

	log.Printf("Generated auth URL for %s: %s", provider, url)

	return url
}

// handleGoogleCallback handles the Google OAuth callback
func (auth *AuthService) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := auth.googleConf.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	// Get user info
	client := auth.googleConf.Client(context.TODO(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response: %v", err), http.StatusInternalServerError)
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(userData, &userInfo); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse user info: %v", err), http.StatusInternalServerError)
		return
	}

	user := &User{
		ID:       fmt.Sprintf("%v", userInfo["id"]),
		Username: fmt.Sprintf("%v", userInfo["name"]),
		Email:    fmt.Sprintf("%v", userInfo["email"]),
		Token:    token.AccessToken,
		Provider: ProviderGoogle,
	}

	// Store session
	sessionID := fmt.Sprintf("%s_%s", ProviderGoogle, user.ID)
	auth.sessions[sessionID] = &Session{
		User:      user,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(auth.config.SessionDuration),
	}

	auth.currentUser = user
	auth.authChan <- user

	// Show success page
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<html>
		<head><title>Login Successful</title></head>
		<body>
			<h1>Login Successful!</h1>
			<p>You can now close this window and return to the application.</p>
			<p>Logged in as: %s (%s)</p>
		</body>
		</html>
	`, user.Username, user.Email)
}

// handleGitHubCallback handles the GitHub OAuth callback (similar improvements)
func (auth *AuthService) handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := auth.githubConf.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	client := auth.githubConf.Client(context.TODO(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response: %v", err), http.StatusInternalServerError)
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(userData, &userInfo); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse user info: %v", err), http.StatusInternalServerError)
		return
	}

	user := &User{
		ID:       fmt.Sprintf("%v", userInfo["id"]),
		Username: fmt.Sprintf("%v", userInfo["login"]),
		Email:    "",
		Token:    token.AccessToken,
		Provider: ProviderGitHub,
	}

	// Try to get primary email if available
	if email, ok := userInfo["email"].(string); ok && email != "" {
		user.Email = email
	} else {
		// Fallback to /user/emails
		resp2, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer resp2.Body.Close()
			emailData, _ := io.ReadAll(resp2.Body)
			var emails []map[string]interface{}
			if json.Unmarshal(emailData, &emails) == nil {
				for _, e := range emails {
					if primary, ok := e["primary"].(bool); ok && primary {
						if em, ok := e["email"].(string); ok {
							user.Email = em
							break
						}
					}
				}
			}
		}
	}

	sessionID := fmt.Sprintf("%s_%s", ProviderGitHub, user.ID)
	auth.sessions[sessionID] = &Session{
		User:      user,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(auth.config.SessionDuration),
	}

	auth.currentUser = user
	auth.authChan <- user

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<html>
		<head><title>Login Successful</title></head>
		<body>
			<h1>Login Successful!</h1>
			<p>You can now close this window and return to the application.</p>
			<p>Logged in as: %s</p>
		</body>
		</html>
	`, user.Username)
}

// GetCurrentUser returns the currently authenticated user
func (auth *AuthService) GetCurrentUser() *User {
	return auth.currentUser
}

// IsAuthenticated returns true if a user is currently authenticated
func (auth *AuthService) IsAuthenticated() bool {
	return auth.currentUser != nil
}

// Logout logs out the current user
func (auth *AuthService) Logout() {
	auth.currentUser = nil
}

// WaitForAuth waits for an OAuth authentication to complete
func (auth *AuthService) WaitForAuth(timeout time.Duration) (*User, error) {
	select {
	case user := <-auth.authChan:
		return user, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("authentication timeout after %v", timeout)
	}
}

// GetAuthChannel returns the authentication channel
func (auth *AuthService) GetAuthChannel() <-chan *User {
	return auth.authChan
}

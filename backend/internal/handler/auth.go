package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"MCPark/internal/config"
	"MCPark/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

type AuthHandler struct {
	db          *sql.DB
	cfg         *config.Config
	oauthConfig *oauth2.Config
}

func NewAuthHandler(db *sql.DB, cfg *config.Config) *AuthHandler {
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	return &AuthHandler{db: db, cfg: cfg, oauthConfig: oauthCfg}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
		return
	}

	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := h.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	user, err := h.findOrCreateUser(userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	jwtToken, err := h.generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, h.cfg.FrontendURL+"/auth/callback?token="+jwtToken)
}

func (h *AuthHandler) GoogleTokenLogin(c *gin.Context) {
	var req struct {
		Credential string `json:"credential" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	payload, err := idtoken.Validate(context.Background(), req.Credential, h.cfg.GoogleClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google token"})
		return
	}

	googleID, _ := payload.Claims["sub"].(string)
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	if googleID == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	user, err := h.findOrCreateUser(googleID, email, name, picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user"})
		return
	}

	jwtToken, err := h.generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{Token: jwtToken, User: *user})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user model.User
	err := h.db.QueryRow("SELECT id, google_id, email, name, avatar, is_admin, created_at FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.Avatar, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) findOrCreateUser(googleID, email, name, avatar string) (*model.User, error) {
	var user model.User
	err := h.db.QueryRow("SELECT id, google_id, email, name, avatar, is_admin, created_at FROM users WHERE google_id = ?", googleID).
		Scan(&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.Avatar, &user.IsAdmin, &user.CreatedAt)

	if err == sql.ErrNoRows {
		isAdmin := email == h.cfg.AdminEmail
		isAdminInt := 0
		if isAdmin {
			isAdminInt = 1
		}
		result, err := h.db.Exec("INSERT INTO users (google_id, email, name, avatar, is_admin) VALUES (?, ?, ?, ?, ?)",
			googleID, email, name, avatar, isAdminInt)
		if err != nil {
			return nil, err
		}
		user.ID, _ = result.LastInsertId()
		user.GoogleID = googleID
		user.Email = email
		user.Name = name
		user.Avatar = avatar
		user.IsAdmin = isAdmin
		user.CreatedAt = time.Now()
	} else if err != nil {
		return nil, err
	} else {
		// Update is_admin status based on current config
		isAdmin := email == h.cfg.AdminEmail
		if isAdmin != user.IsAdmin {
			isAdminInt := 0
			if isAdmin {
				isAdminInt = 1
			}
			h.db.Exec("UPDATE users SET is_admin = ? WHERE id = ?", isAdminInt, user.ID)
			user.IsAdmin = isAdmin
		}
	}

	return &user, nil
}

func (h *AuthHandler) generateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"name":     user.Name,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

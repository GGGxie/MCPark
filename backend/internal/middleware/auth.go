package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		claims, err := validateToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		setClaims(c, claims)
		c.Next()
	}
}

func AuthOptional(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			if claims, err := validateToken(token, jwtSecret); err == nil {
				setClaims(c, claims)
			}
		}
		c.Next()
	}
}

func AdminRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		claims, err := validateToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		isAdmin, _ := claims["is_admin"].(bool)
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		setClaims(c, claims)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearer := c.GetHeader("Authorization")
	if strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimPrefix(bearer, "Bearer ")
	}
	return ""
}

// validateToken verifies an HS256 JWT signed with secret. Tokens must carry
// an expiry, and must include a numeric user_id claim.
func validateToken(tokenStr, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if _, ok := claims["user_id"].(float64); !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func setClaims(c *gin.Context, claims jwt.MapClaims) {
	if id, ok := claims["user_id"].(float64); ok {
		c.Set("user_id", int64(id))
	}
	if email, ok := claims["email"].(string); ok {
		c.Set("email", email)
	}
	if name, ok := claims["name"].(string); ok {
		c.Set("name", name)
	}
	if isAdmin, ok := claims["is_admin"].(bool); ok {
		c.Set("is_admin", isAdmin)
	}
}

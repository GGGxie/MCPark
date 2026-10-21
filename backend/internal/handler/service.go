package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"MCPark/internal/model"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	db    *sql.DB
	dedup *actionDedup
}

func NewServiceHandler(db *sql.DB, dedup *actionDedup) *ServiceHandler {
	return &ServiceHandler{db: db, dedup: dedup}
}

func NewDedup() *actionDedup {
	return newActionDedup()
}

func (h *ServiceHandler) ListServices(c *gin.Context) {
	page, pageSize := parsePaging(c, 12, 50)
	sort := c.DefaultQuery("sort", "hottest")
	category := c.Query("category")
	search := c.Query("search")

	offset := (page - 1) * pageSize

	query := `SELECT id, name, title, title_zh, description, desc_zh, category, icon, image, author_id, author_name, author_avatar,
		tags, status, featured, likes, views, install_count, created_at, updated_at FROM services WHERE status = 'approved'`
	countQuery := `SELECT COUNT(*) FROM services WHERE status = 'approved'`
	args := []interface{}{}

	if category != "" {
		query += " AND category = ?"
		countQuery += " AND category = ?"
		args = append(args, category)
	}

	if search != "" {
		query += " AND (title LIKE ? OR description LIKE ? OR tags LIKE ? OR title_zh LIKE ? OR desc_zh LIKE ?)"
		countQuery += " AND (title LIKE ? OR description LIKE ? OR tags LIKE ? OR title_zh LIKE ? OR desc_zh LIKE ?)"
		s := "%" + search + "%"
		args = append(args, s, s, s, s, s)
	}

	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	h.db.QueryRow(countQuery, countArgs...).Scan(&total)

	switch sort {
	case "latest":
		query += " ORDER BY created_at DESC"
	default:
		query += " ORDER BY likes DESC, views DESC"
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}
	defer rows.Close()

	services := []model.Service{}
	for rows.Next() {
		var s model.Service
		rows.Scan(&s.ID, &s.Name, &s.Title, &s.TitleZH, &s.Description, &s.DescZH, &s.Category, &s.Icon, &s.Image,
			&s.AuthorID, &s.AuthorName, &s.AuthorAvatar, &s.Tags, &s.Status, &s.Featured,
			&s.Likes, &s.Views, &s.InstallCount, &s.CreatedAt, &s.UpdatedAt)
		services = append(services, s)
	}

	c.JSON(http.StatusOK, model.PaginatedResponse{
		Data:       services,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

func (h *ServiceHandler) GetFeaturedServices(c *gin.Context) {
	rows, err := h.db.Query(`SELECT id, name, title, title_zh, description, desc_zh, category, icon, image, author_id, author_name, author_avatar,
		tags, status, featured, likes, views, install_count, created_at, updated_at
		FROM services WHERE featured = 1 AND status = 'approved' ORDER BY likes DESC LIMIT 6`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch featured services"})
		return
	}
	defer rows.Close()

	services := []model.Service{}
	for rows.Next() {
		var s model.Service
		rows.Scan(&s.ID, &s.Name, &s.Title, &s.TitleZH, &s.Description, &s.DescZH, &s.Category, &s.Icon, &s.Image,
			&s.AuthorID, &s.AuthorName, &s.AuthorAvatar, &s.Tags, &s.Status, &s.Featured,
			&s.Likes, &s.Views, &s.InstallCount, &s.CreatedAt, &s.UpdatedAt)
		services = append(services, s)
	}

	c.JSON(http.StatusOK, services)
}

func (h *ServiceHandler) GetCategories(c *gin.Context) {
	rows, err := h.db.Query("SELECT DISTINCT category FROM services WHERE status = 'approved' ORDER BY category")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	categories := []string{}
	for rows.Next() {
		var cat string
		rows.Scan(&cat)
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}

func (h *ServiceHandler) GetService(c *gin.Context) {
	id := c.Param("id")

	ok, err := h.approvedServiceExists(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	key := fmt.Sprintf("view:service:%s:%s", c.ClientIP(), id)
	if h.dedup.Allow(key, 1*time.Hour) {
		if _, err := h.db.Exec("UPDATE services SET views = views + 1 WHERE id = ?", id); err != nil {
			log.Printf("update views failed: %v", err)
		}
	}

	var s model.Service
	err = h.db.QueryRow(`SELECT id, name, title, title_zh, description, desc_zh, category, icon, image, config_json,
		content, content_zh, author_id, author_name, author_avatar, tags, status, featured, likes, views, install_count, created_at, updated_at
		FROM services WHERE id = ? AND status = 'approved'`, id).
		Scan(&s.ID, &s.Name, &s.Title, &s.TitleZH, &s.Description, &s.DescZH, &s.Category, &s.Icon, &s.Image, &s.ConfigJSON,
			&s.Content, &s.ContentZH, &s.AuthorID, &s.AuthorName, &s.AuthorAvatar, &s.Tags, &s.Status, &s.Featured,
			&s.Likes, &s.Views, &s.InstallCount, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	var req model.CreateServiceRequest
	if !bindJSON(c, &req) {
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("name")

	var avatar string
	h.db.QueryRow("SELECT avatar FROM users WHERE id = ?", userID).Scan(&avatar)

	result, err := h.db.Exec(`INSERT INTO services (name, title, title_zh, description, desc_zh, category, icon, image, config_json, content, content_zh, tags, author_id, author_name, author_avatar, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')`,
		req.Name, req.Title, req.TitleZH, req.Description, req.DescZH, req.Category,
		req.Icon, req.Image, req.ConfigJSON, req.Content, req.ContentZH, req.Tags,
		userID, userName, avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Service submitted successfully. Awaiting review."})
}

func (h *ServiceHandler) LikeService(c *gin.Context) {
	serviceID := c.Param("id")
	userID, _ := c.Get("user_id")

	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like"})
		return
	}
	defer tx.Rollback() // no-op after a successful Commit

	var one int
	err = tx.QueryRow("SELECT 1 FROM services WHERE id = ? AND status = 'approved'", serviceID).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like"})
		return
	}

	res, err := tx.Exec("DELETE FROM likes WHERE user_id = ? AND service_id = ?", userID, serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like"})
		return
	}
	removed, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like"})
		return
	}

	liked := removed == 0
	if liked {
		if _, err = tx.Exec("INSERT INTO likes (user_id, service_id) VALUES (?, ?)", userID, serviceID); err == nil {
			_, err = tx.Exec("UPDATE services SET likes = likes + 1 WHERE id = ?", serviceID)
		}
	} else {
		_, err = tx.Exec("UPDATE services SET likes = MAX(likes - 1, 0) WHERE id = ?", serviceID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like"})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"liked": liked})
}

func (h *ServiceHandler) InstallService(c *gin.Context) {
	serviceID := c.Param("id")
	ok, err := h.approvedServiceExists(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update install count"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	key := fmt.Sprintf("install:%s:%s", c.ClientIP(), serviceID)
	if !h.dedup.Allow(key, 24*time.Hour) {
		c.JSON(http.StatusOK, gin.H{"message": "Already counted", "counted": false})
		return
	}
	if _, err := h.db.Exec("UPDATE services SET install_count = install_count + 1 WHERE id = ?", serviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update install count"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Install count updated", "counted": true})
}

// approvedServiceExists reports whether id refers to an approved service.
func (h *ServiceHandler) approvedServiceExists(id string) (bool, error) {
	var one int
	err := h.db.QueryRow("SELECT 1 FROM services WHERE id = ? AND status = 'approved'", id).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (h *ServiceHandler) GetComments(c *gin.Context) {
	serviceID := c.Param("id")
	rows, err := h.db.Query("SELECT id, service_id, user_id, user_name, avatar, content, created_at FROM comments WHERE service_id = ? ORDER BY created_at DESC", serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	comments := []model.Comment{}
	for rows.Next() {
		var cm model.Comment
		rows.Scan(&cm.ID, &cm.ServiceID, &cm.UserID, &cm.UserName, &cm.Avatar, &cm.Content, &cm.CreatedAt)
		comments = append(comments, cm)
	}

	c.JSON(http.StatusOK, comments)
}

func (h *ServiceHandler) CreateComment(c *gin.Context) {
	serviceID := c.Param("id")
	var req model.CreateCommentRequest
	if !bindJSON(c, &req) {
		return
	}
	ok, err := h.approvedServiceExists(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("name")

	var avatar string
	h.db.QueryRow("SELECT avatar FROM users WHERE id = ?", userID).Scan(&avatar)

	result, err := h.db.Exec("INSERT INTO comments (service_id, user_id, user_name, avatar, content) VALUES (?, ?, ?, ?, ?)",
		serviceID, userID, userName, avatar, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// maxPage bounds the page query parameter so offset arithmetic cannot overflow.
const maxPage = 100000

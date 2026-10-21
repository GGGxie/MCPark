package handler

import (
	"database/sql"
	"log"
	"math"
	"net/http"

	"MCPark/internal/model"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) ListAllServices(c *gin.Context) {
	page, pageSize := parsePaging(c, 20, 100)
	status := c.DefaultQuery("status", "")

	offset := (page - 1) * pageSize

	query := `SELECT id, name, title, title_zh, description, desc_zh, category, icon, image, author_id, author_name, author_avatar,
		tags, status, featured, likes, views, install_count, created_at, updated_at FROM services WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM services WHERE 1=1`
	args := []interface{}{}

	if status != "" {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, status)
	}

	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	h.db.QueryRow(countQuery, countArgs...).Scan(&total)

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
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

func (h *AdminHandler) ApproveService(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec("UPDATE services SET status = 'approved', updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve service"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service approved"})
}

func (h *AdminHandler) RejectService(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec("UPDATE services SET status = 'rejected', updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject service"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service rejected"})
}

func (h *AdminHandler) DeleteService(c *gin.Context) {
	id := c.Param("id")
	if err := h.deleteServiceTx(id); err != nil {
		log.Printf("delete service %s failed: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted"})
}

// deleteServiceTx removes a service together with its likes and comments in a
// single transaction so a failure cannot leave partial data behind.
func (h *AdminHandler) deleteServiceTx(id string) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op after a successful Commit
	for _, q := range []string{
		"DELETE FROM likes WHERE service_id = ?",
		"DELETE FROM comments WHERE service_id = ?",
		"DELETE FROM services WHERE id = ?",
	} {
		if _, err := tx.Exec(q, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"MCPark/internal/model"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	db    *sql.DB
	dedup *actionDedup
}

func NewArticleHandler(db *sql.DB, dedup *actionDedup) *ArticleHandler {
	return &ArticleHandler{db: db, dedup: dedup}
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	category := c.Query("category")

	query := `SELECT id, title, title_zh, description, desc_zh, category, content, content_zh, author_id, author_name, views, created_at FROM articles`
	args := []interface{}{}

	if category != "" {
		query += " WHERE category = ?"
		args = append(args, category)
	}
	query += " ORDER BY created_at DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}
	defer rows.Close()

	articles := []model.Article{}
	for rows.Next() {
		var a model.Article
		rows.Scan(&a.ID, &a.Title, &a.TitleZH, &a.Description, &a.DescZH, &a.Category, &a.Content, &a.ContentZH,
			&a.AuthorID, &a.AuthorName, &a.Views, &a.CreatedAt)
		articles = append(articles, a)
	}

	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id := c.Param("id")

	key := fmt.Sprintf("view:article:%s:%s", c.ClientIP(), id)
	if h.dedup.Allow(key, 1*time.Hour) {
		h.db.Exec("UPDATE articles SET views = views + 1 WHERE id = ?", id)
	}

	var a model.Article
	err := h.db.QueryRow(`SELECT id, title, title_zh, description, desc_zh, category, content, content_zh, author_id, author_name, views, created_at FROM articles WHERE id = ?`, id).
		Scan(&a.ID, &a.Title, &a.TitleZH, &a.Description, &a.DescZH, &a.Category, &a.Content, &a.ContentZH,
			&a.AuthorID, &a.AuthorName, &a.Views, &a.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, a)
}

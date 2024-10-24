package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
)
func GetAllNews(c *gin.Context) {
	news := []models.News{}
	config.DB.Find(&news)
	c.JSON(http.StatusOK, &news)
}

func GetNews(c *gin.Context) {
	var news []models.News
	var total int64

	// Default values for pagination
	limit := 10  // Default limit of 10 items per page
	page := 1    // Default to the first page

	// Parse limit and page from query parameters (if provided)
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}

	// Calculate the offset (skip items)
	offset := (page - 1) * limit

	// Count the total number of records
	config.DB.Model(&models.News{}).Count(&total)

	// Fetch the paginated results
	if err := config.DB.Limit(limit).Offset(offset).Find(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching news"})
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"data":       news,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

func GetNewsByID(c *gin.Context) {
	var news models.News
	id := c.Param("id")
	if err := config.DB.Where("id = ?", id).First(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}
	c.JSON(http.StatusOK, &news)
}

func CreateNews(c *gin.Context) {
	var news models.News

	// Bind JSON to the news model, which includes image URL and other details
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate required fields (optional, can be handled by a validation library)
	if news.TitleEn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title in English is required"})
		return
	}

	// Set CreatedAt to the current time
	news.CreatedAt = time.Now()

	// Save the news entry to the database
	if err := config.DB.Create(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news entry"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          news.ID,
		"title_en":   news.TitleEn,
		"title_ua":   news.TitleUa,
		"subtitle_en": news.SubtitleEn,
		"subtitle_ua": news.SubtitleUa,
		"content_en":  news.ContentEn,
		"content_ua":  news.ContentUa,
		"image_url":   news.ImageUrl,
		"created_at":  news.CreatedAt,
	})
}


func UpdateNews(c *gin.Context) {
	var news models.News
	if err := config.DB.Where("id = ?", c.Param("id")).First(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	if err := c.BindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Save(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news"})
		return
	}

	c.JSON(http.StatusOK, &news)
}

func DeleteNews(c *gin.Context) {
	var news models.News
	if err := config.DB.Where("id = ?", c.Param("id")).First(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	if err := config.DB.Delete(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news"})
		return
	}

	// Optionally, delete the file from the server
	if err := os.Remove("." + news.ImageUrl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News deleted"})
}
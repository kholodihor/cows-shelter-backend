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

func GetAllGalleries(c *gin.Context) {
	galleries := []models.Gallery{}
	config.DB.Find(&galleries)
	c.JSON(http.StatusOK, &galleries)
}

func GetGalleries(c *gin.Context) {
	var galleries []models.Gallery
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
	config.DB.Model(&models.Gallery{}).Count(&total)

	// Fetch the paginated results
	if err := config.DB.Limit(limit).Offset(offset).Find(&galleries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching galleries"})
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"data":       galleries,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

func GetGalleryByID(c *gin.Context) {
	var gallery models.Gallery
	id := c.Param("id")
	if err := config.DB.Where("id = ?", id).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}
	c.JSON(http.StatusOK, &gallery)
}

func CreateGallery(c *gin.Context) {
	var gallery models.Gallery

	// Bind JSON to the gallery model, which includes the image URL
	if err := c.ShouldBindJSON(&gallery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Save the gallery entry to the database
	gallery.CreatedAt = time.Now()
	if err := config.DB.Create(&gallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gallery entry"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        gallery.ID,
		"image_url": gallery.ImageUrl,
	})
}

func UpdateGallery(c *gin.Context) {
	var gallery models.Gallery
	if err := config.DB.Where("id = ?", c.Param("id")).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}

	if err := c.BindJSON(&gallery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Save(&gallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gallery"})
		return
	}

	c.JSON(http.StatusOK, &gallery)
}

func DeleteGallery(c *gin.Context) {
	var gallery models.Gallery
	if err := config.DB.Where("id = ?", c.Param("id")).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}

	if err := config.DB.Delete(&gallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete gallery"})
		return
	}

	// Optionally, delete the file from the server
	if err := os.Remove("." + gallery.ImageUrl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gallery deleted"})
}

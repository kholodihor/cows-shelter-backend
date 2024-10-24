package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
)

func GetAllExcursions(c *gin.Context) {
	excursions := []models.Excursion{}
	config.DB.Find(&excursions)
	c.JSON(http.StatusOK, &excursions)
}

func GetExcursions(c *gin.Context) {
	var excursions []models.Excursion
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
	if err := config.DB.Model(&models.Excursion{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting excursions"})
		return
	}

	// Fetch the paginated results
	if err := config.DB.Limit(limit).Offset(offset).Find(&excursions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching excursions"})
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"excursions":   excursions,
		"totalLength":  total, // Total number of excursions in the database
		"page":         page,
		"limit":        limit,
		"totalPages":   (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}

func GetExcursionByID(c *gin.Context) {
	var excursion models.Excursion
	id := c.Param("id")
	if err := config.DB.Where("id = ?", id).First(&excursion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Excursion not found"})
		return
	}
	c.JSON(http.StatusOK, &excursion)
}

func CreateExcursion(c *gin.Context) {
	var excursion models.Excursion
	if err := c.ShouldBindJSON(&excursion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Save the excursion to the database
	if err := config.DB.Create(&excursion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create excursion"})
		return
	}

	// Return the created excursion details
	c.JSON(http.StatusCreated, gin.H{
		"id":               excursion.ID,
		"title_en":        excursion.TitleEn,
		"title_ua":        excursion.TitleUa,
		"description_en":   excursion.DescriptionEn,
		"description_ua":   excursion.DescriptionUa,
		"time_to":         excursion.TimeTo,
		"time_from":       excursion.TimeFrom,
		"amount_of_persons": excursion.AmountOfPersons,
		"image_url":       excursion.ImageUrl,
		"created_at":      excursion.CreatedAt,
	})
}

func UpdateExcursion(c *gin.Context) {
	var excursion models.Excursion
	if err := config.DB.Where("id = ?", c.Param("id")).First(&excursion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Excursion not found"})
		return
	}

	if err := c.BindJSON(&excursion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Save(&excursion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update excursion"})
		return
	}

	c.JSON(http.StatusOK, &excursion)
}

func DeleteExcursion(c *gin.Context) {
	var excursion models.Excursion
	if err := config.DB.Where("id = ?", c.Param("id")).First(&excursion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Excursion not found"})
		return
	}

	if err := config.DB.Delete(&excursion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete excursion"})
		return
	}

	// Optionally, delete the file from the server
	if err := os.Remove("." + excursion.ImageUrl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Excursion deleted"})
}
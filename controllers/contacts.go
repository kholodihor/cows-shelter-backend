package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/config"
	"github.com/kholodihor/cows-shelter-backend/models"
)

// GetContacts - Retrieve all contacts
func GetContacts(c *gin.Context) {
	var contacts []models.Contact
	config.DB.Find(&contacts)
	c.JSON(http.StatusOK, &contacts)
}

// GetContactByID - Retrieve a specific contact by ID
func GetContactByID(c *gin.Context) {
	var contact models.Contact
	id := c.Param("id")

	if err := config.DB.Where("id = ?", id).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}
	c.JSON(http.StatusOK, &contact)
}

// CreateContact - Create a new contact entry
func CreateContact(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone" binding:"required"`
	}

	// Validate the request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Create the new contact entry
	contact := models.Contact{
		Email: requestBody.Email,
		Phone: requestBody.Phone,
	}

	if err := config.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating contact"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    contact.Id,
		"email": contact.Email,
		"phone": contact.Phone,
	})
}

// UpdateContact - Update an existing contact
func UpdateContact(c *gin.Context) {
	id := c.Param("id")
	var contact models.Contact

	if err := config.DB.Where("id = ?", id).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	var requestBody struct {
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	contact.Email = requestBody.Email
	contact.Phone = requestBody.Phone

	config.DB.Save(&contact)
	c.JSON(http.StatusOK, &contact)
}

// DeleteContact - Delete a specific contact by ID
func DeleteContact(c *gin.Context) {
	id := c.Param("id")
	var contact models.Contact

	// Check if contact exists before deleting
	if err := config.DB.Where("id = ?", id).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	// Delete the contact
	config.DB.Delete(&contact)
	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

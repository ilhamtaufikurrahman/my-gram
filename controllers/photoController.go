package controllers

import (
	"my-gram/database"
	"my-gram/helpers"
	"my-gram/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreatePhoto(c *gin.Context) {
	var (
		db          = database.GetDB()
		userData    = c.MustGet("userData").(jwt.MapClaims)
		userId      = uint(userData["id"].(float64))
		contentType = helpers.GetContentType(c)
		Photo       = models.Photo{}
		err         error
	)

	if contentType == appJSON {
		err = c.ShouldBindJSON(&Photo)
	} else {
		err = c.ShouldBind(&Photo)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	if Photo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your photo title is required",
		})
		return
	}

	if Photo.PhotoUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your photo url is required",
		})
		return
	}

	Photo.UserId = uint(userId)

	err = db.Debug().Create(&Photo).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Photo.Id,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserId,
		"created_at": Photo.CreatedAt,
	})
}

func GetPhotos(c *gin.Context) {
	var (
		db             = database.GetDB()
		Photos         = []models.Photo{}
		PhotosResponse = []models.PhotosResponse{}
		err            error
	)

	err = db.Preload("User").Find(&Photos).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	for _, photo := range Photos {
		PhotosResponse = append(PhotosResponse, models.PhotosResponse{
			Id:        photo.Id,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoUrl:  photo.PhotoUrl,
			UserId:    photo.UserId,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User: models.UserPhoto{
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		})
	}

	c.JSON(http.StatusOK, PhotosResponse)
}

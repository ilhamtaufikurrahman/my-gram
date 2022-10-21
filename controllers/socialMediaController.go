package controllers

import (
	"my-gram/database"
	"my-gram/helpers"
	"my-gram/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	SocialMedia := models.SocialMedia{}
	userId := uint(userData["id"].(float64))
	var err error

	if contentType == appJSON {
		err = c.ShouldBindJSON(&SocialMedia)
	} else {
		err = c.ShouldBind(&SocialMedia)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	SocialMedia.UserId = userId

	err = db.Debug().Create(&SocialMedia).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               SocialMedia.Id,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"created_at":       SocialMedia.CreatedAt,
	})
}

func UpdateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	SocialMedia := models.SocialMedia{}
	var err error

	socialMediaId, err := strconv.Atoi(c.Param("socialMediaId"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Invalid parameter",
		})
		return
	}

	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		err = c.ShouldBindJSON(&SocialMedia)
	} else {
		err = c.ShouldBind(&SocialMedia)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	err = db.Select("user_id").First(&SocialMedia, socialMediaId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": "Data doesnt exist",
		})
		return
	}

	if SocialMedia.UserId != userId {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "You are not allowed to access this data",
		})
		return
	}

	SocialMedia.Id = uint(socialMediaId)

	err = db.Model(&SocialMedia).Where("id=?", socialMediaId).Updates(SocialMedia).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               SocialMedia.Id,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"updated_at":       SocialMedia.UpdatedAt,
	})
}

func GetSocialMedias(c *gin.Context) {
	db := database.GetDB()
	SocialMedias := []models.SocialMedia{}
	SocialMediasResponse := []models.SocialMediasResponse{}

	err := db.Preload("User").Find(&SocialMedias).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	for _, socialMedia := range SocialMedias {
		SocialMediasResponse = append(SocialMediasResponse, models.SocialMediasResponse{
			Id:             socialMedia.Id,
			Name:           socialMedia.Name,
			SocialMediaUrl: socialMedia.SocialMediaUrl,
			UserId:         socialMedia.UserId,
			CreatedAt:      socialMedia.CreatedAt,
			UpdatedAt:      socialMedia.UpdatedAt,
			User: models.UserSocialMedia{
				Id:       socialMedia.User.Id,
				Username: socialMedia.User.Username,
				Email:    socialMedia.User.Email,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"social_medias": SocialMediasResponse,
	})
}

func DeleteSocialMedia(c *gin.Context) {
	db := database.GetDB()
	socialMediaId, err := strconv.Atoi(c.Param("socialMediaId"))
	SocialMedia := models.SocialMedia{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := uint(userData["id"].(float64))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Invalid parameter",
		})
		return
	}

	err = db.Select("user_id").First(&SocialMedia, socialMediaId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": "Data doesnt exist",
		})
		return
	}

	if SocialMedia.UserId != userId {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "You are not allowed to access this data",
		})
		return
	}

	err = db.Delete(models.SocialMedia{}, "id", socialMediaId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Error deleting item",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfuly deleted",
	})
}

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

func CreateComment(c *gin.Context) {
	var (
		db          = database.GetDB()
		userData    = c.MustGet("userData").(jwt.MapClaims)
		userId      = uint(userData["id"].(float64))
		contentType = helpers.GetContentType(c)
		Comment     = models.Comment{}
		NewComment  = models.Comment{}
		err         error
	)

	if contentType == appJSON {
		err = c.ShouldBindJSON(&Comment)
	} else {
		err = c.ShouldBind(&Comment)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	if Comment.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your message is required",
		})
		return
	}

	if Comment.PhotoId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Photo id is required",
		})
		return
	}

	err = db.Model(models.Photo{}).Select("id").First(&NewComment, Comment.PhotoId).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": "Photo id doesnt exist",
		})
		return
	}

	Comment.UserId = userId

	err = db.Debug().Create(&Comment).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         Comment.Id,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoId,
		"user_id":    Comment.UserId,
		"created_at": Comment.CreatedAt,
	})
}

func GetComments(c *gin.Context) {
	var (
		db               = database.GetDB()
		Comment          = []models.Comment{}
		CommentsResponse = []models.CommentsResponse{}
		err              error
	)

	err = db.Preload("User").Preload("Photo").Find(&Comment).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	for _, comment := range Comment {
		CommentsResponse = append(CommentsResponse, models.CommentsResponse{
			Id:        comment.Id,
			Message:   comment.Message,
			PhotoId:   comment.PhotoId,
			UserId:    comment.UserId,
			UpdatedAt: comment.UpdatedAt,
			CreatedAt: comment.CreatedAt,
			User: models.UserComment{
				Id:       comment.User.Id,
				Email:    comment.User.Email,
				Username: comment.User.Username,
			},
			Photo: models.PhotoComment{
				Id:       comment.Photo.Id,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoUrl: comment.Photo.PhotoUrl,
				UserId:   comment.Photo.UserId,
			},
		})
	}

	c.JSON(http.StatusOK, CommentsResponse)
}

func UpdateComment(c *gin.Context) {
	var (
		db             = database.GetDB()
		userData       = c.MustGet("userData").(jwt.MapClaims)
		userId         = uint(userData["id"].(float64))
		commentId, err = strconv.Atoi(c.Param("commentId"))
		contentType    = helpers.GetContentType(c)
		Comment        = models.Comment{}
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Invalid parameter",
		})
		return
	}

	if contentType == appJSON {
		err = c.ShouldBindJSON(&Comment)
	} else {
		err = c.ShouldBind(&Comment)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	if Comment.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your message is required",
		})
		return
	}

	err = db.Select("user_id", "photo_id").First(&Comment, commentId).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": "Data doesnt exist",
		})
		return
	}

	if Comment.UserId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "You are not allowed to access this data",
		})
		return
	}

	Comment.Id = uint(commentId)

	err = db.Debug().Model(&Comment).Where("id=?", commentId).Updates(&Comment).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": "Data doesnt exist",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Comment.Id,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoId,
		"user_id":    Comment.UserId,
		"updated_at": Comment.UpdatedAt,
	})
}

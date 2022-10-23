package controllers

import (
	"my-gram/database"
	"my-gram/helpers"
	"my-gram/models"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

func RegisterUser(c *gin.Context) {
	var (
		db          = database.GetDB()
		contentType = helpers.GetContentType(c)
		User        = models.User{}
		NewUser     = models.User{}
		err         error
	)

	if contentType == appJSON {
		err = c.ShouldBindJSON(&User)
	} else {
		err = c.ShouldBind(&User)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	if User.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your email is required",
		})
		return
	}

	_, errEmail := mail.ParseAddress(User.Email)
	if errEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Invalid email format",
		})
		return
	}

	db.Where("email=?", User.Email).First(&NewUser)

	if NewUser.Email == User.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Email already used",
		})
		return
	}

	if User.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your username is required",
		})
		return
	}

	db.Where("username", User.Username).First(&NewUser)

	if NewUser.Username == User.Username {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Username already used",
		})
		return
	}

	if User.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your password is required",
		})
		return
	}

	if len(User.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Password has to have a minimum length of 6 characters",
		})
		return
	}

	if User.Age == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your age is required",
		})
		return
	}

	if User.Age <= 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Age must be more than 8 years",
		})
		return
	}

	err = db.Debug().Create(&User).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"age":      User.Age,
		"email":    User.Email,
		"id":       User.Id,
		"username": User.Username,
	})
}

func LoginUser(c *gin.Context) {
	var (
		db          = database.GetDB()
		contentType = helpers.GetContentType(c)
		User        = models.User{}
		password    = ""
		err         error
	)

	if contentType == appJSON {
		err = c.ShouldBindJSON(&User)
	} else {
		err = c.ShouldBind(&User)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	password = User.Password

	err = db.Debug().Where("email=?", User.Email).Take(&User).Error

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if err != nil || !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid password or passowrd",
		})
		return
	}

	token, err := helpers.GenerateToken(
		User.Id,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func UpdateUser(c *gin.Context) {
	var (
		db          = database.GetDB()
		userData    = c.MustGet("userData").(jwt.MapClaims)
		contentType = helpers.GetContentType(c)
		User        = models.User{}
		NewUser     = models.User{}
		userId      = userData["id"].(float64)
	)

	paramUserId, err := strconv.Atoi(c.Param("userId"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Invalid parameter",
		})
		return
	}

	if contentType == appJSON {
		err = c.ShouldBindJSON(&User)
	} else {
		err = c.ShouldBind(&User)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	if User.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your email is required",
		})
		return
	}

	_, errEmail := mail.ParseAddress(User.Email)
	if errEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Invalid email format",
		})
		return
	}

	db.Where("email=?", User.Email).First(&NewUser)

	if NewUser.Email == User.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Email already used",
		})
		return
	}

	if User.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Your username is required",
		})
		return
	}

	db.Where("username", User.Username).First(&NewUser)

	if NewUser.Username == User.Username {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Username already used",
		})
		return
	}

	err = db.Select("id", "age").First(&User, paramUserId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": "Data doesnt exist",
		})
		return
	}

	if User.Id != uint(userId) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "You are not allowed to access this data",
		})
		return
	}

	db.Where("email=?", User.Email).First(&NewUser)

	if User.Email == NewUser.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Email has been used",
		})
		return
	}

	db.Where("username=?", User.Username).First(&NewUser)

	if User.Username == NewUser.Username {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Username has been used",
		})
		return
	}

	err = db.Model(&User).Where("id=?", paramUserId).Updates(&User).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         User.Id,
		"email":      User.Email,
		"username":   User.Username,
		"age":        User.Age,
		"updated_at": User.UpdatedAt,
	})
}

func DeleteUser(c *gin.Context) {
	var (
		db               = database.GetDB()
		userData         = c.MustGet("userData").(jwt.MapClaims)
		userId           = uint(userData["id"].(float64))
		paramUserId, err = strconv.Atoi(c.Param("userId"))
		User             = models.User{}
	)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Invalid parameter",
		})
		return
	}

	err = db.Select("id").First(&User, paramUserId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": "Data doesnt exist",
		})
		return
	}

	if userId != uint(paramUserId) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "You are not allowed to access this data",
		})
		return
	}

	err = db.Delete(models.SocialMedia{}, "user_id", userId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Error deleting item",
			"message": err.Error(),
		})
		return
	}

	err = db.Delete(models.Comment{}, "user_id", userId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Error deleting item",
			"message": err.Error(),
		})
		return
	}

	err = db.Delete(models.Photo{}, "user_id", userId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Error deleting item",
			"message": err.Error(),
		})
		return
	}

	err = db.Delete(User, "id", paramUserId).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Error deleting item",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfuly deleted",
	})
}

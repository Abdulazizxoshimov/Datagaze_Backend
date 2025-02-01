package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Abdulazizxoshimov/Datagaze_Backend/entity"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Security  		BearerAuth
// @Summary   		Create User
// @Description 	Api for create a new user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			user body entity.UserCreateRequst true "Create User Model"
// @Success 		201 {object} entity.UserCreateResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/user [POST]
func (h *HandlerV1) CreateUser(c *gin.Context) {
	var user entity.UserCreateRequst

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: "Invalid request payload",
		})
		log.Println("Error parsing request body:", err.Error())
		return
	}
	id := uuid.New().String()

	password, err := validation.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Error hashing password",
		})
		log.Println("Error hashing password:", err.Error())
		return
	}
	user.Password = password
	_, err = h.Service.User().CreateUser(&entity.User{
		ID: id,
		Name: user.Name,
		SurName: user.Surname,
		UserName: user.Username,
		Email: user.Email,
		Password: password,
		Role: "user",

	})
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Error creating user",
		})
		log.Println("Error creating user:", err.Error())
		return
	}
	c.JSON(http.StatusAccepted, entity.UserCreateResponse{
		ID: id,
	})
}

// @Security  		BearerAuth
// @Summary   		Get User
// @Description 	Api for getting user by ID
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID"
// @Success 		200 {object} entity.UserResponse
// @Failure 		400 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/user/{id} [GET]
func (h *HandlerV1) GetUser(c *gin.Context) {
	id := c.Param("id")

	filter := map[string]interface{}{
		"id": id,
	}
	user, err := h.Service.User().GetUser(filter)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: "User not found",
		})
		log.Println("Error fetching user:", err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.UserResponse{
		ID: user.ID,
		Name: user.Name,
		SurName: user.SurName,
		UserName: user.UserName,
		Email: user.Email,
		Role: user.Role,
		RefreshT: "",
		AccessT: "",
	})
}

// @Security  		BearerAuth
// @Summary   		Get All Users with Paging
// @Description 	Api for getting all users with pagination
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			page query int false "Page number (default: 1)"
// @Param 			limit query int false "Number of users per page (default: 10)"
// @Success 		200 {object} entity.PaginatedResponse
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/users [GET]
func (h *HandlerV1) GetAllUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Invalid page parameter",
		})
		return
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Invalid limit parameter",
		})
		return
	}

	users,  err := h.Service.User().GetAllUsers(pageNum, limitNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Error fetching users",
		})
		log.Println("Error fetching users:", err.Error())
		return
	}

	response := entity.PaginatedResponse{
		Data:       *users,
		Page:       pageNum,
		Limit:      limitNum,
		TotalCount: users.UserCount,
		TotalPages: (users.UserCount + limitNum - 1) / limitNum,
	}

	c.JSON(http.StatusOK, response)
}

// @Security  		BearerAuth
// @Summary   		Update User
// @Description 	Api for updating user by ID
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID"
// @Param 			user body entity.UserCreateRequst true "Update User Model"
// @Success 		200 {object} string
// @Failure 		400 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/user/{id} [PUT]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user entity.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Invalid request payload",
		})
		log.Println("Error parsing request body:", err.Error())
		return
	}
	user.ID = id

	_, err := h.Service.User().UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: "User not found or update failed",
		})
		log.Println("Error updating user:", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// @Security  		BearerAuth
// @Summary   		Delete User
// @Description 	Api for deleting user by ID
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "User ID"
// @Success 		200 {object} string
// @Failure 		400 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/user/{id} [DELETE]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.Service.User().DeleteUser(&entity.DeleteRequest{
		ID: id,
		DeletedAt: time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: "User not found or deletion failed",
		})
		log.Println("Error deleting user:", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

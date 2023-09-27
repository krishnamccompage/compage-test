package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/krishnamccompage/compage-test/test-api/pkg/rest/server/daos/clients/nosqls"
	"github.com/krishnamccompage/compage-test/test-api/pkg/rest/server/models"
	"github.com/krishnamccompage/compage-test/test-api/pkg/rest/server/services"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() (*UserController, error) {
	userService, err := services.NewUserService()
	if err != nil {
		return nil, err
	}
	return &UserController{
		userService: userService,
	}, nil
}

func (userController *UserController) CreateUser(context *gin.Context) {
	// validate input
	var input models.User
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger user creation
	if _, err := userController.userService.CreateUser(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (userController *UserController) UpdateUser(context *gin.Context) {
	// validate input
	var input models.User
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// trigger user update
	if _, err := userController.userService.UpdateUser(context.Param("id"), &input); err != nil {
		log.Error(err)
		if errors.Is(err, nosqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, nosqls.ErrInvalidObjectID) {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (userController *UserController) FetchUser(context *gin.Context) {
	// trigger user fetching
	user, err := userController.userService.GetUser(context.Param("id"))
	if err != nil {
		log.Error(err)
		if errors.Is(err, nosqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, nosqls.ErrInvalidObjectID) {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(context.Request.Context())
		currentSpan.SetAttributes(attribute.String("user.id", user.ID))
	}

	context.JSON(http.StatusOK, user)
}

func (userController *UserController) DeleteUser(context *gin.Context) {
	// trigger user deletion
	if err := userController.userService.DeleteUser(context.Param("id")); err != nil {
		log.Error(err)
		if errors.Is(err, nosqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, nosqls.ErrInvalidObjectID) {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (userController *UserController) ListUsers(context *gin.Context) {
	// trigger all users fetching
	users, err := userController.userService.ListUsers()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, users)
}

func (*UserController) PatchUser(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "PATCH",
	})
}

func (*UserController) OptionsUser(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "OPTIONS",
	})
}

func (*UserController) HeadUser(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "HEAD",
	})
}

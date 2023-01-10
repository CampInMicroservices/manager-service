package api

import (
	"context"
	"log"
	"manager-service/db"
	"manager-service/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getUserListRequest struct {
	Offset int32 `form:"offset"`
	Limit  int32 `form:"limit" binding:"required,min=1,max=20"`
}

type createUserRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Activated *bool  `json:"activated" binding:"required"`
}

type LoginParam struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @BasePath /manager-service/v1

// User godoc
// @Summary Users by ID
// @Schemes
// @Description Returns user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} db.User
// @Router /v1/users/{id} [get]
func (server *Server) GetUserByID(ctx *gin.Context) {

	// Check if request has ID field in URI.
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	// Execute query.
	result, err := server.store.GetUserByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// @BasePath /manager-service/v1

// User godoc
// @Summary Users list
// @Schemes
// @Description Returns user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Success 200 {array} db.User
// @Router /v1/users [get]
func (server *Server) GetAllUsers(ctx *gin.Context) {

	// Check if request has parameters offset and limit for pagination.
	var req getUserListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	arg := db.ListUserParam{
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	// Execute query.
	result, err := server.store.GetAllUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// @BasePath /manager-service/v1

// User godoc
// @Summary Users create
// @Schemes
// @Description Creates a user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body db.User true "User"
// @Success 200 {array} db.User
// @Router /v1/users [post]
func (server *Server) CreateUser(ctx *gin.Context) {

	// Check if request has all required fields in json body.
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	arg := db.CreateUserParam{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		Activated: *req.Activated,
	}

	// Execute query.
	result, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// @BasePath /manager-service/v1

// Login godoc
// @Summary User Login
// @Schemes
// @Description Login for a user
// @Tags Login
// @Accept json
// @Produce json
// @Param request body api.LoginParam true "Login data"
// @Success 200 {array} proto.AuthResponse
// @Router /v1/login [post]
func (server *Server) LoginUser(ctx *gin.Context) {

	// Check if request has all required fields in json body.
	var req LoginParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	authRequest := &proto.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Create new payment via gRPC
	authResponse, err := server.grpcClient.Login(context.Background(), authRequest)

	if err != nil {

		st, _ := status.FromError(err)

		if st.Code() == codes.NotFound {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
			ctx.Abort()
			return
		}

		log.Println("Auth service not reachable.")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Auth service unavailable."})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, authResponse)
}

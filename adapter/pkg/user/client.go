package controller

import (
	"github.com/arvians-id/go-mongo/adapter/pb"
	"github.com/arvians-id/go-mongo/adapter/response"
	"github.com/arvians-id/go-mongo/post/cmd/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type UserController struct {
	UserService pb.UserServiceClient
}

func NewUserController(router *gin.Engine, configuration config.Config) *UserController {
	connection, err := grpc.Dial(configuration.Get("USER_SERVICE_URL"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	controller := &UserController{
		UserService: pb.NewUserServiceClient(connection),
	}

	routes := router.Group("/api")
	{
		routes.GET("/users", controller.FindAll)
		routes.GET("/users/:id", controller.FindByID)
		routes.POST("/users", controller.Create)
		routes.PATCH("/users/:id", controller.Update)
		routes.DELETE("/users/:id", controller.Delete)
	}

	return controller
}

func (controller *UserController) FindAll(ctx *gin.Context) {
	users, err := controller.UserService.FindAll(ctx.Request.Context(), new(emptypb.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", users)
}

func (controller *UserController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := controller.UserService.FindByID(ctx, &pb.GetUserByIDRequest{
		ID: id,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", user)
}

func (controller *UserController) Create(ctx *gin.Context) {
	var userRequest pb.CreateUserRequest
	err := ctx.ShouldBindJSON(&userRequest)
	if err != nil {
		response.ReturnErrorBadRequest(ctx, err, nil)
		return
	}

	user, err := controller.UserService.Create(ctx, &userRequest)
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", user)
}

func (controller *UserController) Update(ctx *gin.Context) {
	var userRequest pb.UpdateUserRequest
	err := ctx.ShouldBindJSON(&userRequest)
	if err != nil {
		response.ReturnErrorBadRequest(ctx, err, nil)
		return
	}

	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	userRequest.ID = id.Hex()
	user, err := controller.UserService.Update(ctx, &userRequest)
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", user)
}

func (controller *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := controller.UserService.Delete(ctx, &pb.GetUserByIDRequest{
		ID: id,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", nil)
}

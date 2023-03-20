package post

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

type Controller struct {
	PostService pb.PostServiceClient
}

func NewPostController(router *gin.Engine, configuration config.Config) *Controller {
	connection, err := grpc.Dial(configuration.Get("POST_SERVICE_URL"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	controller := &Controller{
		PostService: pb.NewPostServiceClient(connection),
	}

	routes := router.Group("/api")
	{
		routes.GET("/posts", controller.FindAll)
		routes.GET("/posts/:id", controller.FindByID)
		routes.POST("/posts", controller.Create)
		routes.PATCH("/posts/:id", controller.Update)
		routes.DELETE("/posts/:id", controller.Delete)
	}

	return controller
}

func (controller *Controller) FindAll(ctx *gin.Context) {
	posts, err := controller.PostService.FindAll(ctx.Request.Context(), new(emptypb.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", posts)
}

func (controller *Controller) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := controller.PostService.FindByID(ctx, &pb.GetPostByIDRequest{
		ID: id,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", user)
}

func (controller *Controller) Create(ctx *gin.Context) {
	var request pb.CreatePostRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.ReturnErrorBadRequest(ctx, err, nil)
		return
	}

	post, err := controller.PostService.Create(ctx, &request)
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", post)
}

func (controller *Controller) Update(ctx *gin.Context) {
	var request pb.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.ReturnErrorBadRequest(ctx, err, nil)
		return
	}

	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	request.ID = id.Hex()
	post, err := controller.PostService.Update(ctx, &request)
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", post)
}

func (controller *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := controller.PostService.Delete(ctx, &pb.GetPostByIDRequest{
		ID: id,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(ctx, err, nil)
		return
	}

	response.ReturnSuccessOK(ctx, "OK", nil)
}

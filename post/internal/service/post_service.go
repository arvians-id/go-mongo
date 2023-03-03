package service

import (
	"context"
	"github.com/arvians-id/go-mongo/post/internal/repository"
	"github.com/arvians-id/go-mongo/post/pb"
	util "github.com/arvians-id/go-mongo/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PostService struct {
	PostRepository repository.PostRepository
}

func NewPostService(PostRepository repository.PostRepository) pb.PostServiceServer {
	return &PostService{
		PostRepository: PostRepository,
	}
}

func (service *PostService) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListPostResponse, error) {
	posts, err := service.PostRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListPostResponse{
		Posts: posts,
	}, nil
}

func (service *PostService) FindByID(ctx context.Context, id *pb.GetPostByIDRequest) (*pb.GetPostResponse, error) {
	objectID, err := util.ConvertStringToHex(id.ID)
	if err != nil {
		return nil, err
	}

	post, err := service.PostRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return &pb.GetPostResponse{
		Post: post,
	}, nil
}

func (service *PostService) Create(ctx context.Context, request *pb.CreatePostRequest) (*pb.GetPostResponse, error) {
	var post pb.Post
	post.ID = util.GenerateID().Hex()
	post.Title = request.Title
	post.Content = request.Content
	post.CreatedAt = util.PrimitiveDateToTimestampPB()
	post.UpdatedAt = util.PrimitiveDateToTimestampPB()

	postCreated, err := service.PostRepository.Create(ctx, &post)
	if err != nil {
		return nil, err
	}

	return &pb.GetPostResponse{
		Post: postCreated,
	}, nil
}

func (service *PostService) Update(ctx context.Context, request *pb.UpdatePostRequest) (*pb.GetPostResponse, error) {
	id, err := util.ConvertStringToHex(request.ID)
	if err != nil {
		return nil, err
	}

	post, err := service.PostRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	post.ID = request.ID
	post.Title = request.Title
	post.Content = request.Content
	post.UpdatedAt = util.PrimitiveDateToTimestampPB()

	postUpdated, err := service.PostRepository.Update(ctx, post)
	if err != nil {
		return nil, err
	}

	return &pb.GetPostResponse{
		Post: postUpdated,
	}, nil
}

func (service *PostService) Delete(ctx context.Context, id *pb.GetPostByIDRequest) (*emptypb.Empty, error) {
	objectID, err := util.ConvertStringToHex(id.ID)
	if err != nil {
		return nil, err
	}

	_, err = service.PostRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	err = service.PostRepository.Delete(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return nil, err
}

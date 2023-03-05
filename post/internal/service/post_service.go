package service

import (
	"context"
	"github.com/arvians-id/go-mongo/post/internal/model"
	"github.com/arvians-id/go-mongo/post/internal/repository"
	"github.com/arvians-id/go-mongo/post/pb"
	"github.com/arvians-id/go-mongo/util"
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

	var postsPb []*pb.Post
	for _, post := range posts {
		postsPb = append(postsPb, post.ToPB())
	}

	return &pb.ListPostResponse{
		Posts: postsPb,
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
		Post: post.ToPB(),
	}, nil
}

func (service *PostService) Create(ctx context.Context, request *pb.CreatePostRequest) (*pb.GetPostResponse, error) {
	var post model.Post
	post.Title = request.Title
	post.Content = request.Content
	post.CreatedAt = util.PrimitiveDateTime()
	post.UpdatedAt = util.PrimitiveDateTime()

	postCreated, err := service.PostRepository.Create(ctx, &post)
	if err != nil {
		return nil, err
	}

	return &pb.GetPostResponse{
		Post: postCreated.ToPB(),
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

	post.Title = request.Title
	post.Content = request.Content
	post.UpdatedAt = util.PrimitiveDateTime()

	postUpdated, err := service.PostRepository.Update(ctx, post)
	if err != nil {
		return nil, err
	}

	return &pb.GetPostResponse{
		Post: postUpdated.ToPB(),
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

	return new(emptypb.Empty), nil
}

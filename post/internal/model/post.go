package model

import (
	"github.com/arvians-id/go-mongo/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	User      *User              `bson:"user" json:"user,omitempty"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (post *Post) ToPBWithUser() *pb.Post {
	return &pb.Post{
		ID:        post.ID.Hex(),
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID.Hex(),
		User:      post.User.ToPB(),
		CreatedAt: post.CreatedAt.Time().String(),
		UpdatedAt: post.UpdatedAt.Time().String(),
	}
}

func (post *Post) ToPB() *pb.Post {
	return &pb.Post{
		ID:        post.ID.Hex(),
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID.Hex(),
		CreatedAt: post.CreatedAt.Time().String(),
		UpdatedAt: post.UpdatedAt.Time().String(),
	}
}

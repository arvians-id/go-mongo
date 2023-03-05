package repository

import (
	"context"
	"github.com/arvians-id/go-mongo/post/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type PostRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.Post, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Post, error)
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) (*model.Post, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type PostRepository struct {
	DB *mongo.Database
}

func NewPostRepository(db *mongo.Database) PostRepository {
	return PostRepository{
		DB: db,
	}
}

func (repository *PostRepository) FindAll(ctx context.Context) ([]*model.Post, error) {
	rows, err := repository.DB.Collection("posts").Find(ctx, bson.M{})
	if err != nil {
		log.Println("[PostRepository][FindAll] problem querying to db, err: ", err.Error())
		return nil, err
	}
	defer func(rows *mongo.Cursor, ctx context.Context) {
		err := rows.Close(ctx)
		if err != nil {
			log.Println("[PostRepository][FindAll] problem closing db rows, err: ", err.Error())
			return
		}
	}(rows, ctx)

	var Posts []*model.Post
	for rows.Next(ctx) {
		var post model.Post
		err := rows.Decode(&post)
		if err != nil {
			log.Println("[PostRepository][FindAll] problem with scanning db row, err: ", err.Error())
			return nil, err
		}
		Posts = append(Posts, &post)
	}

	return Posts, nil
}

func (repository *PostRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	row := repository.DB.Collection("posts").FindOne(ctx, bson.M{"_id": id})

	var post model.Post
	err := row.Decode(&post)
	if err != nil {
		log.Println("[PostRepository][FindByID] problem with scanning db row, err: ", err.Error())
		return nil, err
	}

	return &post, nil
}

func (repository *PostRepository) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	row, err := repository.DB.Collection("posts").InsertOne(ctx, post)
	if err != nil {
		log.Println("[PostRepository][Create] problem querying to db, err: ", err.Error())
		return nil, err
	}

	post.ID = row.InsertedID.(primitive.ObjectID)

	return post, nil
}

func (repository *PostRepository) Update(ctx context.Context, post *model.Post) (*model.Post, error) {
	_, err := repository.DB.Collection("posts").UpdateOne(ctx, bson.M{
		"_id": post.ID,
	}, bson.M{"$set": post})
	if err != nil {
		log.Println("[PostRepository][Update] problem querying to db, err: ", err.Error())
		return nil, err
	}

	return post, nil
}

func (repository *PostRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := repository.DB.Collection("posts").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Println("[PostRepository][Delete] problem querying to db, err: ", err.Error())
		return err
	}

	return nil
}

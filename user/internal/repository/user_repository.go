package repository

import (
	"context"
	"github.com/arvians-id/go-mongo/user/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type UserRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type UserRepository struct {
	DB *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return UserRepository{
		DB: db,
	}
}

func (repository *UserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	rows, err := repository.DB.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		log.Println("[UserRepository][FindAll] problem querying to db, err: ", err.Error())
		return nil, err
	}
	defer func(rows *mongo.Cursor, ctx context.Context) {
		err := rows.Close(ctx)
		if err != nil {
			log.Println("[UserRepository][FindAll] problem closing db rows, err: ", err.Error())
			return
		}
	}(rows, ctx)

	var users []*model.User
	for rows.Next(ctx) {
		var user model.User
		err := rows.Decode(&user)
		if err != nil {
			log.Println("[UserRepository][FindAll] problem with scanning db row, err: ", err.Error())
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (repository *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	row := repository.DB.Collection("users").FindOne(ctx, bson.M{"_id": id})

	var user model.User
	err := row.Decode(&user)
	if err != nil {
		log.Println("[UserRepository][FindByID] problem with scanning db row, err: ", err.Error())
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	row, err := repository.DB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		log.Println("[UserRepository][Create] problem querying to db, err: ", err.Error())
		return nil, err
	}

	user.ID = row.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (repository *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := repository.DB.Collection("users").UpdateOne(ctx, bson.M{
		"_id": user.ID,
	}, bson.M{"$set": user})
	if err != nil {
		log.Println("[UserRepository][Update] problem querying to db, err: ", err.Error())
		return nil, err
	}

	return user, nil
}

func (repository *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := repository.DB.Collection("users").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Println("[UserRepository][Delete] problem querying to db, err: ", err.Error())
		return err
	}

	return nil
}

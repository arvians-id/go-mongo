package util

import (
	"github.com/arvians-id/go-mongo/post/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GenerateID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func ConvertStringToHex(id string) (primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return objectId, nil
}

func PrimitiveDateToTimestampPB() *timestamppb.Timestamp {
	dateTime := primitive.NewDateTimeFromTime(util.CurrentTime())
	result := timestamppb.New(dateTime.Time())

	return result
}

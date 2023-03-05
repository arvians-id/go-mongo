package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

func PrimitiveDateTime() primitive.DateTime {
	return primitive.NewDateTimeFromTime(CurrentTime())
}

func ConvertPrimitiveDateTimeToString(dateTime time.Time) string {
	dt := primitive.NewDateTimeFromTime(time.Now())
	tm := dt.Time()
	str := tm.Format("2006-01-02 15:04:05")
	return str
}

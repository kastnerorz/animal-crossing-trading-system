package testdata

import (
	"github.com/kastnerorz/animal-crossing-trading-system/backend/tools"
	"go.mongodb.org/mongo-driver/bson"
)

var ReviewerID = "5e82ffb52fc4557c9d343b46"
var ApplicantId = "5e82ffb52fc4557c9d343b47"

func TestUsers() []interface{} {
	return []interface{}{
		bson.M{
			"_id":              tools.ObjectID(ReviewerID),
			"username":         "zed",
			"password":         "b5609be1d47cca6abca4004253321f530fc4ff022e11683e9c36fc77769b24ae",
			"nickname":         "张豆",
			"switchFriendCode": "SW-1234-1234-1234",
			"jikeId":           "张豆",
		},
		bson.M{
			"_id":              tools.ObjectID(ApplicantId),
			"username":         "zed1",
			"password":         "b5609be1d47cca6abca4004253321f530fc4ff022e11683e9c36fc77769b24ae",
			"nickname":         "张豆",
			"switchFriendCode": "SW-1234-1234-1234",
			"jikeId":           "张豆",
		},
	}
}

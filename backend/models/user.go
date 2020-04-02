package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Username         string             `json:"username,omitempty" bson:"username"`
	Nickname         string             `json:"nickname" bson:"nickname"`
	Password         string             `json:"password,omitempty" bson:"password"`
	SwitchFriendCode string             `json:"switchFriendCode,omitempty" bson:"switchFriendCode"`
	JikeID           string             `json:"jikeId,omitempty" bson:"jikeId"`
}

type Credentials struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

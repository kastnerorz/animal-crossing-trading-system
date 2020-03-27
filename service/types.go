package main

import "time"

type User struct {
	ID               string `json:"id" bson:"_id"`
	Username         string `json:"username" bson:"username"`
	Nickname         string `json:"nickname" bson:"nickname"`
	Password         string `json:"password" bson:"password"`
	SwitchFriendCode string `json:"switchFriendCode" bson:"switchFriendCode"`
}

type Quotation struct {
	ID               string    `json:"id" bson:"_id"`
	Type             string    `json:"type" bson:"type"`
	Price            int       `json:"price" bson:"price"`
	Author           User      `json:"author" bson:"author"`
	ParticipantCount int       `json:"participantCount" bson:"participantCount"`
	Verified         bool      `json:"verified" bson:"verified"`
	LastModified     time.Time `json:"lastModified" bson:"lastModified"`
}

type QuotationParam struct {
	ID               string    `json:"id" bson:"_id"`
	Type             string    `json:"type" bson:"type"`
	Price            *int      `json:"price" bson:"price"`
	ParticipantCount *int      `json:"participantCount" bson:"participantCount"`
	Verified         *bool     `json:"verified" bson:"verified"`
	LastModified     time.Time `json:"lastModified" bson:"lastModified"`
}

var QuotationType = map[string]struct{}{
	"SELL": {},
	"BUY":  {},
}

type Credentials struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Application struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Applicant        User               `json:"applicant" bson:"applicant"`
	QuotationId      primitive.ObjectID `json:"quotationId" bson:"quotationId"`
	QuotationType    string             `json:"quotationType" bson:"quotationType"`
	Price            int                `json:"price" bson:"price"`
	ReviewerId       primitive.ObjectID `json:"reviewerId" bson:"reviewerId"`
	ReviewerNickname string             `json:"reviewerNickname" bson:"reviewerNickname"`
	PassCode         string             `json:"passCode" bson:"passCode"`
	SwitchFriendCode string             `json:"switchFriendCode" bson:"switchFriendCode"`
	Status           string             `json:"status" bson:"status"`
	LastModified     time.Time          `json:"lastModified" bson:"lastModified"`
}

type ApplicationParam struct {
	Type             string `json:"type" bson:"type"`
	QuotationId      string `json:"quotationId" bson:"quotationId"`
	Status           string `json:"status" bson:"status"`
	PassCode         string `json:"passCode" bson:"passCode"`
	SwitchFriendCode string `json:"switchFriendCode" bson:"switchFriendCode"`
}

var ApplicationType = map[string]struct{}{
	"APPLY":  {},
	"REVIEW": {},
}

var ApplicationStatus = map[string]struct{}{
	"PENDING": {},
	"ACCEPT":  {},
	"REJECT":  {},
}

package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID               string `json:"id" bson:"_id"`
	Username         string `json:"username,omitempty" bson:"username"`
	Nickname         string `json:"nickname" bson:"nickname"`
	Password         string `json:"password,omitempty" bson:"password"`
	SwitchFriendCode string `json:"switchFriendCode,omitempty" bson:"switchFriendCode"`
	JikeID           string `json:"jikeId,omitempty" bson:"jikeId"`
}

type Quotation struct {
	ID           string    `json:"id" bson:"_id"`
	Type         string    `json:"type" bson:"type"`
	Price        int       `json:"price" bson:"price"`
	Author       User      `json:"author" bson:"author"`
	ValidCount   int       `json:"validCount" bson:"validCount"`
	InvalidCount int       `json:"invalidCount" bson:"invalidCount"`
	OpenType     string    `json:"openType" bson:"openType"`
	PassCode     string    `json:"passCode,omitempty" bson:"passCode"`
	HandlingFee  int       `json:"handlingFee" bson:"handlingFee"`
	LastModified time.Time `json:"lastModified" bson:"lastModified"`
}

type QuotationParam struct {
	ID           string    `json:"id" bson:"_id"`
	Type         string    `json:"type" bson:"type"`
	Price        *int      `json:"price" bson:"price"`
	IsValid      *bool     `json:"isValid" bson:"isValid"`
	OpenType     string    `json:"openType" bson:"openType"`
	PassCode     string    `json:"passCode" bson:"passCode"`
	HandlingFee  *int      `json:"handlingFee" bson:"handlingFee"`
	LastModified time.Time `json:"lastModified" bson:"lastModified"`
}

var QuotationType = map[string]struct{}{
	"SELL": {},
	"BUY":  {},
}

var OpenType = map[string]struct{}{
	"PASS_CODE": {},
	"FRIENDS":   {},
}

type Application struct {
	Applicant        User               `json:"applicant" bson:"applicant"`
	Reviewer         primitive.ObjectID `json:"reviewer" bson:"reviewer"`
	PassCode         string             `json:"passCode" bson:"passCode"`
	SwitchFriendCode string             `json:"switchFriendCode" bson:"switchFriendCode"`
	Status           string             `json:"status" bson:"status"`
}

var ApplicationStatus = map[string]struct{}{
	"PENDING": {},
	"ACCEPT":  {},
	"REJECT":  {},
}

type Credentials struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

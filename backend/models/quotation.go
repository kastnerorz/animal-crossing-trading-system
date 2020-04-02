package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Quotation struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Type         string             `json:"type" bson:"type"`
	Price        int                `json:"price" bson:"price"`
	Author       User               `json:"author" bson:"author"`
	ValidCount   int                `json:"validCount" bson:"validCount"`
	InvalidCount int                `json:"invalidCount" bson:"invalidCount"`
	OpenType     string             `json:"openType" bson:"openType"`
	PassCode     string             `json:"passCode,omitempty" bson:"passCode"`
	HandlingFee  int                `json:"handlingFee" bson:"handlingFee"`
	LastModified time.Time          `json:"lastModified" bson:"lastModified"`
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

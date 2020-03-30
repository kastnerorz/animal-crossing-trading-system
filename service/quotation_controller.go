package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateQuotation(c *gin.Context) {
	o, _ := c.Get(IdentityKey)
	user := o.(User)

	var quotation Quotation
	err := c.BindJSON(&quotation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "内部错误！"})
		log.Println(err)
		return
	}

	if _, ok := QuotationType[quotation.Type]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "报价类型不正确！"})
		return
	}

	if quotation.Price == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": -3, "msg": "报价不合法！"})
		return
	}

	mongoCtx, collection := GetMongoContext("quotations")
	user.Password = ""
	_, err = collection.InsertOne(mongoCtx, bson.M{
		"type":             quotation.Type,
		"author":           user,
		"price":            quotation.Price,
		"participantCount": 0,
		"verified":         false,
		"lastModified":     time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -3, "msg": "Error while inserting into database."})
		log.Println(err)
		return
	}
}

func GetQuotations(c *gin.Context) {
	quotationType := c.Query("type")
	verified := c.Query("verified")
	available := c.Query("available")

	filter := bson.M{}

	if quotationType != "" {
		if _, ok := QuotationType[quotationType]; !ok {
			c.JSON(http.StatusOK, []struct{}{})
			return
		}

		filter["type"] = quotationType
	}

	if verified != "" {
		filter["verified"], _ = strconv.ParseBool(verified)
	}

	if available != "" {
		filter["available"], _ = strconv.ParseBool(available)
	}

	lowerBound, upperBound := GetValidDateLowerAndUpperBound()
	filter["lastModified"] = bson.M{
		"$gt":  lowerBound,
		"$lte": upperBound,
	}
	mongoCtx, collection := GetMongoContext("quotations")
	opts := options.Find()
	opts.SetSort(bson.D{{"price", -1}})
	opts.SetLimit(10)
	cursor, err := collection.Find(mongoCtx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}
	var res []Quotation
	if err = cursor.All(mongoCtx, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -2, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}
	cursor.Close(mongoCtx)
	if res == nil {
		res = []Quotation{}
	}
	c.JSON(http.StatusOK, res)
	return
}

func GetMyQuotation(c *gin.Context) {
	user, _ := c.Get(IdentityKey)
	username := user.(*User).Username

	filter := bson.M{
		"author.username": username,
	}

	lowerBound, upperBound := GetValidDateLowerAndUpperBound()
	filter["lastModified"] = bson.M{
		"$gt":  lowerBound,
		"$lte": upperBound,
	}

	mongoCtx, collection := GetMongoContext("quotations")
	opts := options.Find()
	opts.SetSort(bson.D{{"lastModified", -1}})
	opts.SetLimit(1)
	cursor, err := collection.Find(mongoCtx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}
	var res []Quotation
	if err = cursor.All(mongoCtx, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -2, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}
	cursor.Close(mongoCtx)
	if res == nil {
		res = []Quotation{}
	}
	c.JSON(http.StatusOK, res)
	return
}

func UpdateQuotation(c *gin.Context) {
	var param QuotationParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "内部错误！"})
		log.Println(err)
		return
	}
	participantCount := param.ParticipantCount
	verified := param.Verified

	set := bson.M{}

	if participantCount != nil {
		set["participantCount"] = participantCount
	}

	if verified != nil {
		set["verified"] = verified
	}

	var quotation Quotation
	mongoCtx, collection := GetMongoContext("quotations")
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	opt := options.FindOneAndUpdate()
	opt.SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(mongoCtx, bson.M{"_id": objectId}, bson.M{"$set": set}, opt).Decode(&quotation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "更新报价信息失败！"})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, quotation)
}

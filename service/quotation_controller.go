package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateQuotation(c *gin.Context) {
	user, _ := c.Get(IdentityKey)
	user = user.(*User)

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

		filter["quotationType"] = quotationType
	}

	if verified != "" {
		filter["verified"], _ = strconv.ParseBool(verified)
	}

	if available != "" {
		filter["available"], _ = strconv.ParseBool(available)
	}

	mongoCtx, collection := GetMongoContext("quotations")
	opts := options.Find()
	opts.SetSort(bson.D{{"price", -1}})
	sortCursor, err := collection.Find(mongoCtx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}
	var res []Quotation
	if err = sortCursor.All(mongoCtx, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -2, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}

}

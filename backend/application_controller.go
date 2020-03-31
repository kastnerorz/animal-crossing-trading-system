package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func CreateApplication(c *gin.Context) {
	user := GetUserFromContext(c)

	var applicationParam ApplicationParam
	err := c.BindJSON(&applicationParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误！"})
		log.Println(err)
		return
	}

	if applicationParam.QuotationId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "报价ID不能为空！"})
		log.Println(err)
		return
	}

	mongoCtx, collection := GetMongoContext("quotations")
	var quotation Quotation
	objectId, _ := primitive.ObjectIDFromHex(applicationParam.QuotationId)
	err = collection.FindOne(mongoCtx, bson.M{"_id": objectId}).Decode(&quotation)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误"})
		log.Println(err)
		return
	}

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "报价不存在！"})
		log.Println(err)
		return
	}

	mongoCtx, collection = GetMongoContext("applications")
	user.Password = ""
	user.Username = ""
	quotationId, _ := primitive.ObjectIDFromHex(quotation.ID)
	reviewerId, _ := primitive.ObjectIDFromHex(quotation.Author.ID)
	_, err = collection.InsertOne(mongoCtx, bson.M{
		"applicant":          user,
		"quotationId":        quotationId,
		"quotationType":      quotation.OpenType,
		"reviewerNickname":   quotation.Author.Nickname,
		"reviewerId":         reviewerId,
		"status":             "PENDING",
		"lastModified":       time.Now(),
		"passCode":           "",
		"switchFriendNumber": "",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -3, "msg": "（-3）内部错误"})
		log.Println(err)
		return
	}
	c.Status(http.StatusCreated)
}

func GetMyApplications(c *gin.Context) {
	user := GetUserFromContext(c)

	filter := bson.M{}
	applicationType := c.Query("type")
	if applicationType != "" {
		if _, ok := ApplicationType[applicationType]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "类型不合法"})
			return
		}
		if applicationType == "REVIEW" {
			filter["reviewerId"], _ = primitive.ObjectIDFromHex(user.ID)
		} else if applicationType == "APPLY" {
			filter["applicant"] = user.ID
		}
	}

	lowerBound, upperBound := GetValidDateLowerAndUpperBound()
	filter["lastModified"] = bson.M{
		"$gt":  lowerBound,
		"$lte": upperBound,
	}
	mongoCtx, collection := GetMongoContext("applications")
	opts := options.Find()
	opts.SetSort(bson.D{{"lastModified", -1}})
	cursor, err := collection.Find(mongoCtx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "申请查询失败！"})
		log.Println(err)
		return
	}
	var res []Application
	if err = cursor.All(mongoCtx, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -2, "msg": "申请查询失败！"})
		log.Println(err)
		return
	}
	cursor.Close(mongoCtx)
	if res == nil {
		res = []Application{}
	}

	c.JSON(http.StatusOK, res)
}

func UpdateApplication(c *gin.Context) {
	user := GetUserFromContext(c)

	var applicationParam ApplicationParam
	err := c.BindJSON(&applicationParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误！"})
		log.Println(err)
		return
	}

	mongoCtx, collection := GetMongoContext("applications")
	var application Application
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	err = collection.FindOne(mongoCtx, bson.M{"_id": objectId}).Decode(&application)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误"})
		log.Println(err)
		return
	}
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "申请不存在！"})
		log.Println(err)
		return
	}

	if user.ID != application.ReviewerId.Hex() {
		c.JSON(http.StatusForbidden, gin.H{"code": -2, "msg": "无权限！"})
		log.Println(err)
		return
	}

	if applicationParam.Status == "" {
		if _, ok := ApplicationStatus[applicationParam.Status]; !ok || applicationParam.Status == "PENDING" {
			c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "申请结果不合法！"})
			log.Println(err)
			return
		}
	}
	fmt.Println(applicationParam)
	set := bson.M{}
	if applicationParam.Status == "ACCEPT" {
		if application.QuotationType == "PASS_CODE" {
			if applicationParam.PassCode == "" {
				c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "密码不能为空！"})
				log.Println(err)
				return
			}
			set["passCode"] = applicationParam.PassCode
			mongoCtx, collection = GetMongoContext("quotations")
			_, err = collection.UpdateOne(mongoCtx, bson.M{"_id": application.QuotationId}, bson.M{"$set": bson.M{
				"passCode": applicationParam.PassCode,
			}})
			if err != nil && err != mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误"})
				log.Println(err)
				return
			}
		} else if application.QuotationType == "FRIENDS" {
			if applicationParam.SwitchFriendCode == "" {
				c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "好友编号不能为空！"})
				log.Println(err)
				return
			}
			set["switchFriendCode"] = applicationParam.SwitchFriendCode
		}
	}
	set["status"] = applicationParam.Status

	mongoCtx, collection = GetMongoContext("applications")
	_, err = collection.UpdateOne(mongoCtx, bson.M{"_id": objectId}, bson.M{"$set": set})
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误"})
		log.Println(err)
		return
	}

	c.Status(http.StatusOK)
}

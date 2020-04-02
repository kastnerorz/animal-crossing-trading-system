package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/models"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/pkg"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func CreateApplication(c *gin.Context) {
	user := tools.GetUserFromContext(c)

	var applicationParam models.ApplicationParam
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

	mongoCtx, collection := pkg.GetMongoContext("quotations")
	var quotation models.Quotation
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

	mongoCtx, collection = pkg.GetMongoContext("applications")
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
	user := tools.GetUserFromContext(c)

	filter := bson.M{}
	applicationType := c.Query("type")
	if applicationType != "" {
		if _, ok := models.ApplicationType[applicationType]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "类型不合法"})
			return
		}
		if applicationType == "REVIEW" {
			filter["reviewerId"], _ = primitive.ObjectIDFromHex(user.ID)
		} else if applicationType == "APPLY" {
			filter["applicant._id"] = user.ID
		}
	}

	lowerBound, upperBound := tools.GetValidDateLowerAndUpperBound()
	filter["lastModified"] = bson.M{
		"$gt":  lowerBound,
		"$lte": upperBound,
	}
	mongoCtx, collection := pkg.GetMongoContext("applications")
	opts := options.Find()
	opts.SetSort(bson.D{{"lastModified", -1}})
	cursor, err := collection.Find(mongoCtx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "申请查询失败！"})
		log.Println(err)
		return
	}
	var res []models.Application
	if err = cursor.All(mongoCtx, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -2, "msg": "申请查询失败！"})
		log.Println(err)
		return
	}
	cursor.Close(mongoCtx)
	if res == nil {
		res = []models.Application{}
	}

	c.JSON(http.StatusOK, res)
}

func UpdateApplication(c *gin.Context) {
	user := tools.GetUserFromContext(c)

	var applicationParam models.ApplicationParam
	err := c.BindJSON(&applicationParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误！"})
		log.Println(err)
		return
	}

	mongoCtx, collection := pkg.GetMongoContext("applications")
	var application models.Application
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
		if _, ok := models.ApplicationStatus[applicationParam.Status]; !ok || applicationParam.Status == "PENDING" {
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
			mongoCtx, collection = pkg.GetMongoContext("quotations")
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

	mongoCtx, collection = pkg.GetMongoContext("applications")
	_, err = collection.UpdateOne(mongoCtx, bson.M{"_id": objectId}, bson.M{"$set": set})
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误"})
		log.Println(err)
		return
	}

	c.Status(http.StatusOK)
}

func DeleteApplication(c *gin.Context) {
	user := tools.GetUserFromContext(c)

	mongoCtx, collection := pkg.GetMongoContext("applications")
	_, err := collection.DeleteOne(mongoCtx, bson.M{"_id": tools.ObjectID(c.Param("id")), "applicant._id": user.ID})
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误"})
		log.Println(err)
		return
	}
	if err == mongo.ErrNoDocuments {
		c.Status(http.StatusNotFound)
		return
	}
	c.Status(http.StatusOK)
}

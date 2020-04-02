package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/middlewares"
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

func CreateQuotation(c *gin.Context) {
	o, _ := c.Get(middlewares.IdentityKey)
	user := o.(*models.User)

	var quotation models.Quotation
	err := c.BindJSON(&quotation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误！"})
		log.Println(err)
		return
	}

	if _, ok := models.QuotationType[quotation.Type]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "报价类型不正确！"})
		return
	}

	if _, ok := models.OpenType[quotation.OpenType]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "岛屿开放类型不正确！"})
		return
	}

	if quotation.Price == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": -3, "msg": "报价不合法！"})
		return
	}

	mongoCtx, collection := pkg.GetMongoContext("quotations")
	user.Password = ""
	user.SwitchFriendCode = ""
	_, err = collection.InsertOne(mongoCtx, bson.M{
		"type":         quotation.Type,
		"author":       user,
		"price":        quotation.Price,
		"validCount":   0,
		"invalidCount": 0,
		"openType":     quotation.OpenType,
		"passCode":     quotation.PassCode,
		"handlingFee":  quotation.HandlingFee,
		"lastModified": time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -3, "msg": "Error while inserting into database."})
		log.Println(err)
		return
	}
	c.Status(http.StatusCreated)
}

func GetQuotations(c *gin.Context) {
	quotationType := c.Query("type")
	openType := c.Query("openType")
	//isValid := c.Query("isValid")

	filter := bson.M{}

	if quotationType != "" {
		if _, ok := models.QuotationType[quotationType]; !ok {
			c.JSON(http.StatusOK, []struct{}{})
			return
		}

		filter["type"] = quotationType
	}

	if openType != "" {
		if _, ok := models.OpenType[openType]; !ok {
			c.JSON(http.StatusOK, []struct{}{})
			return
		}
		filter["openType"] = openType
	}

	//if isValid != "" {
	//	v, _ := strconv.ParseBool(isValid);
	//	if v {
	//		filter["$where"] = "this.validCount > this.inValidCount"
	//	} else {
	//		filter["$where"] = "this.validCount < this.inValidCount"
	//	}
	//}

	lowerBound, upperBound := tools.GetValidDateLowerAndUpperBound()
	filter["lastModified"] = bson.M{
		"$gt":  lowerBound,
		"$lte": upperBound,
	}
	mongoCtx, collection := pkg.GetMongoContext("quotations")
	opts := options.Find()
	opts.SetSort(bson.D{{"price", -1}})
	opts.SetLimit(10)
	opts.SetProjection(bson.D{
		{"passCode", 0},
	})
	cursor, err := collection.Find(mongoCtx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}
	var res []models.Quotation
	if err = cursor.All(mongoCtx, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -2, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}
	cursor.Close(mongoCtx)
	if res == nil {
		res = []models.Quotation{}
	}
	c.JSON(http.StatusOK, res)
	return
}

func GetQuotation(c *gin.Context) {
	mongoCtx, collection := pkg.GetMongoContext("quotations")
	var res models.Quotation
	err := collection.FindOne(mongoCtx, bson.M{"_id": tools.ObjectID(c.Param("id"))}).Decode(&res)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误"})
		log.Println(err)
		return
	}
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, struct{}{})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
func GetMyQuotation(c *gin.Context) {
	o, _ := c.Get(middlewares.IdentityKey)
	userId := o.(*models.User).ID
	filter := bson.M{
		"author._id": userId,
	}

	lowerBound, upperBound := tools.GetValidDateLowerAndUpperBound()
	filter["lastModified"] = bson.M{
		"$gt":  lowerBound,
		"$lte": upperBound,
	}

	mongoCtx, collection := pkg.GetMongoContext("quotations")
	opts := options.Find()
	opts.SetSort(bson.D{{"lastModified", -1}})
	opts.SetLimit(1)
	cursor, err := collection.Find(mongoCtx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}
	var res []models.Quotation
	if err = cursor.All(mongoCtx, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -2, "msg": "报价查询失败！"})
		log.Println(err)
		return
	}
	cursor.Close(mongoCtx)
	if res == nil {
		res = []models.Quotation{}
	}
	c.JSON(http.StatusOK, res)
	return
}

func UpdateQuotation(c *gin.Context) {
	o, _ := c.Get(middlewares.IdentityKey)
	user := o.(*models.User)

	var param models.QuotationParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "内部错误！"})
		log.Println(err)
		return
	}
	price := param.Price
	openType := param.OpenType
	passCode := param.PassCode
	handlingFee := param.HandlingFee

	set := bson.M{}

	if price != nil {
		set["price"] = price
	}

	if openType != "" {
		if _, ok := models.OpenType[openType]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "岛屿开放类型不正确！"})
			return
		}
		set["openType"] = openType
	}

	if passCode != "" {
		set["passCode"] = passCode
	}

	if handlingFee != nil {
		set["handlingFee"] = handlingFee
	}

	var quotation models.Quotation
	mongoCtx, collection := pkg.GetMongoContext("quotations")
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	opt := options.FindOneAndUpdate()
	opt.SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(mongoCtx, bson.M{"_id": objectId, "author._id": user.ID}, bson.M{"$set": set}, opt).Decode(&quotation)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "更新报价信息失败！"})
		log.Println(err)
		return
	}

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": "没有这个报价信息或无权限更改！"})
		return
	}
	c.JSON(http.StatusOK, quotation)
}

func DeleteQuotation(c *gin.Context) {
	user := tools.GetUserFromContext(c)

	mongoCtx, collection := pkg.GetMongoContext("quotations")
	_, err := collection.DeleteOne(mongoCtx, bson.M{"_id": tools.ObjectID(c.Param("id")), "author._id": user.ID})
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
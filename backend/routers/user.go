package routers

import (
	"crypto/sha256"
	"encoding/hex"
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
)

func CreateUser(c *gin.Context) {

	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "（-1）内部错误！"})
		log.Println(err)
		return
	}
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "用户名不能为空！"})
		log.Println(err)
		return
	}

	mongoCtx, collection := pkg.GetMongoContext("users")
	var res models.User
	err = collection.FindOne(mongoCtx, bson.M{"username": user.Username}).Decode(&res)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -3, "msg": "（-3）内部错误！"})
		log.Println(err)
		return
	}
	if res.Username == user.Username {
		c.JSON(http.StatusBadRequest, gin.H{"code": -4, "msg": "用户名已被使用！"})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -5, "msg": "密码不能为空！"})
		log.Println(err)
		return
	}

	if user.Nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -6, "msg": "昵称不能为空！"})
		log.Println(err)
		return
	}

	if user.SwitchFriendCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -7, "msg": "好友编号不能为空！"})
		log.Println(err)
		return
	}

	if user.SwitchNickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -8, "msg": "NS昵称不能为空！"})
		log.Println(err)
		return
	}

	h := sha256.New()
	h.Write([]byte(user.Password))
	mongoCtx, collection = pkg.GetMongoContext("users")
	insertResult, err := collection.InsertOne(mongoCtx, bson.M{
		"username":         user.Username,
		"password":         hex.EncodeToString(h.Sum(nil)),
		"nickname":         user.Nickname,
		"switchFriendCode": user.SwitchFriendCode,
		"switchNickname":   user.SwitchNickname,
		"jikeId":           user.JikeID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -8, "msg": "（-8）内部错误！"})
		log.Println(err)
		return
	}
	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	c.JSON(http.StatusCreated, gin.H{"username": user.Username, "id": id})
}

func GetUser(c *gin.Context) {
	mongoCtx, collection := pkg.GetMongoContext("users")
	var res models.User
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	opt := options.FindOne()
	opt.SetProjection(bson.D{
		{"_id", 1},
		{"username", 1},
		{"nickname", 1},
		{"jikeId", 1},
	})
	err := collection.FindOne(mongoCtx, bson.M{"_id": objectId}, opt).Decode(&res)
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

func UpdateUser(c *gin.Context) {
	user := tools.GetUserFromContext(c)

	var param models.User
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误！"})
		log.Println(err)
		return
	}

	set := bson.M{}
	applicationSet := bson.M{}
	quotationSet := bson.M{}
	if param.Nickname != "" {
		set["nickname"] = param.Nickname
		applicationSet["applicant.nickname"] = param.Nickname
		quotationSet["author.nickname"] = param.Nickname
	}

	if param.SwitchFriendCode != "" {
		set["switchFriendCode"] = param.SwitchFriendCode
		applicationSet["applicant.switchFriendCode"] = param.SwitchFriendCode
		quotationSet["author.switchFriendCode"] = param.SwitchFriendCode
	}

	if param.JikeID != "" {
		set["jikeId"] = param.JikeID
		applicationSet["applicant.jikeId"] = param.JikeID
		quotationSet["author.jikeId"] = param.JikeID
	}

	mongoCtx, collection := pkg.GetMongoContext("users")
	_, err = collection.UpdateOne(mongoCtx, bson.M{"_id": user.ID}, bson.M{"$set": set})
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-1）内部错误"})
		log.Println(err)
		return
	}

	mongoCtx, collection = pkg.GetMongoContext("applications")
	_, err = collection.UpdateMany(mongoCtx, bson.M{"applicant._id": user.ID}, bson.M{"$set": applicationSet})
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-2）内部错误"})
		log.Println(err)
		return
	}

	mongoCtx, collection = pkg.GetMongoContext("quotations")
	_, err = collection.UpdateMany(mongoCtx, bson.M{"author._id": user.ID}, bson.M{"$set": quotationSet})
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "（-3）内部错误"})
		log.Println(err)
		return
	}

	c.Status(http.StatusOK)
}
func GetMyInfo(c *gin.Context) {
	user := tools.GetUserFromContext(c)
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

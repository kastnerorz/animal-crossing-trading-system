package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func CreateUser(c *gin.Context) {

	var user User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "内部错误！"})
		log.Println(err)
		return
	}
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "用户名不能为空！"})
		log.Println(err)
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -3, "msg": "密码不能为空！"})
		log.Println(err)
		return
	}

	if user.Nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -4, "msg": "昵称不能为空！"})
		log.Println(err)
		return
	}

	if user.SwitchFriendCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -5, "msg": "好友编号不能为空！"})
		log.Println(err)
		return
	}

	h := sha256.New()
	h.Write([]byte(user.Password))
	mongoCtx, collection := GetMongoContext("users")
	_, err = collection.InsertOne(mongoCtx, bson.M{
		"username":         user.Username,
		"password":         hex.EncodeToString(h.Sum(nil)),
		"nickname":         user.Nickname,
		"switchFriendCode": user.SwitchFriendCode,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -3, "msg": "Error while inserting into database."})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"username": user.Username})
}

func GetUser(c *gin.Context) {
	mongoCtx, collection := GetMongoContext("users")
	var res User
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	err := collection.FindOne(mongoCtx, bson.M{"_id": objectId}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "用户获取失败."})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":               res.ID,
		"username":         res.Username,
		"nickname":         res.Nickname,
		"switchFriendCode": res.SwitchFriendCode,
	})
}

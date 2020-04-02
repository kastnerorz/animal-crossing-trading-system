package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/models"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var IdentityKey = "username"

func AuthMiddleware() *jwt.GinJWTMiddleware {

	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					IdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			mongoCtx, collection := pkg.GetMongoContext("users")
			var user models.User

			err := collection.FindOne(mongoCtx, bson.M{"username": claims[IdentityKey].(string)}).Decode(&user)
			if err != nil && err != mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "用户获取失败."})
				log.Println(err)
				return nil
			}
			user.Password = ""
			return &user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var credentials models.Credentials
			if err := c.BindJSON(&credentials); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			h := sha256.New()
			h.Write([]byte(credentials.Password))
			passwordHash := hex.EncodeToString(h.Sum(nil))
			mongoCtx, collection := pkg.GetMongoContext("users")
			var res models.User

			err := collection.FindOne(mongoCtx, bson.M{"username": credentials.Username}).Decode(&res)
			if err != nil && err != mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "msg": "用户获取失败."})
				log.Println(err)
				return nil, jwt.ErrFailedTokenCreation
			}
			if res.Password == passwordHash {
				return &res, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code": -1,
				"msg":  message,
			})
		},
		HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
			switch e {
			case jwt.ErrExpiredToken:
				return "Token已过期！请重新登录"
			case jwt.ErrEmptyAuthHeader:
				return "无效的Token"
			case jwt.ErrMissingLoginValues:
				return "用户名或密码不能为空！"
			case jwt.ErrFailedAuthentication:
				return "用户名或密码错误！"
			}
			return e.Error()
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return middleware
}

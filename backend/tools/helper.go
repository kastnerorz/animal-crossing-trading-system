package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/middlewares"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func GetValidDateLowerAndUpperBound() (time.Time, time.Time) {
	now := time.Now()
	var lowerBound time.Time
	var upperBound time.Time
	if now.Hour() <= 12 {
		lowerBound = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		upperBound = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local)
	} else {
		lowerBound = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local)
		upperBound = time.Date(now.Year(), now.Month(), now.Day(), 24, 0, 0, 0, time.Local)
	}
	return lowerBound, upperBound
}

func GetUserFromContext(c *gin.Context) *models.User {
	o, _ := c.Get(middlewares.IdentityKey)
	user := o.(*models.User)
	return user
}

func ObjectID(id string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(id)
	return oid
}

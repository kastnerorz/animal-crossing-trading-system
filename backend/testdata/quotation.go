package testdata

import (
	"github.com/kastnerorz/animal-crossing-trading-system/backend/tools"
	"go.mongodb.org/mongo-driver/bson"
)

var QuotationId = "5e8359a0f2f2d8c751c92b39"
var QuotationIdPassCode = "5e8300cd6dab4e49d227ca6b"

func TestQuotations() []interface{} {
	return []interface{}{
		bson.M{
			"_id":   tools.ObjectID(QuotationIdPassCode),
			"type":  "SELL",
			"price": 40,
			"author": bson.M{
				"_id":      tools.ObjectID(ReviewerID),
				"nickname": "张豆",
				"jikeId":   "张豆",
			},
			"validCount":   0,
			"invalidCount": 0,
			"openType":     "PASS_CODE",
			"handlingFee":  "100000",
			"lastModified": "2020-03-31T08:35:25.963Z",
		},
		bson.M{
			"_id":   tools.ObjectID(QuotationId),
			"type":  "SELL",
			"price": 40,
			"author": bson.M{
				"_id":      tools.ObjectID(ReviewerID),
				"nickname": "张豆",
				"jikeId":   "张豆",
			},
			"validCount":   0,
			"invalidCount": 0,
			"openType":     "FRIENDS",
			"handlingFee":  "100000",
			"lastModified": "2020-03-31T14:54:24.143Z",
		},
	}

}

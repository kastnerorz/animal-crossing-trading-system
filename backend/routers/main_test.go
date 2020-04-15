package routers

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/configs"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/pkg"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/testdata"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.StringVar(&configs.MongoURI, "mongo-url", "mongodb://localhost:27017", "")
	flag.StringVar(&configs.MongoCollection, "mongo-collection", "acts-test", "")
	flag.StringVar(&configs.JWTToken, "token", "secret key", "")
	flag.StringVar(&configs.Port, "port", "8080", "")
	flag.Parse()

	//inputOptions := mongoimport.InputOptions{
	//    File: "./testdata/users.json",
	//}
	//opt := mongoimport.Options{
	//    ToolOptions:  nil,
	//    InputOptions: &inputOptions,
	//    ParsedArgs:   nil,
	//}
	//mi, _ := mongoimport.New(opt)
	//mi.SessionProvider, _ = db.NewSessionProvider(options.ToolOptions{
	//    URI: &options.URI{
	//        ConnectionString: "mongodb://localhost:27017/acts-test",
	//        ConnString:       connstring.ConnString{},
	//    },
	//})
	//success, failure, err := mi.ImportDocuments()
	//fmt.Println(success, failure, err)

	//load data
	mongoCtx, collection := pkg.GetMongoContext("users")
	_, err := collection.InsertMany(mongoCtx, testdata.TestUsers())
	if err != nil {
		fmt.Println(err)
	}

	mongoCtx, collection = pkg.GetMongoContext("quotations")
	_, err = collection.InsertMany(mongoCtx, testdata.TestQuotations())
	if err != nil {
		fmt.Println(err)
	}

	// login
	body := []byte(`{"username":"zed","password":"01db71ab8048f74a4b92c26ba77285ade0687ac192758e8185ad52701f649ef2"}`)
	r := PerformRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	var res map[string]string
	json.Unmarshal([]byte(r.Body.String()), &res)
	ReviewerToken, _ = res["token"]

	body = []byte(`{"username":"zed1","password":"01db71ab8048f74a4b92c26ba77285ade0687ac192758e8185ad52701f649ef2"}`)
	r = PerformRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	json.Unmarshal([]byte(r.Body.String()), &res)
	ApplicantToken, _ = res["token"]

	testResult := m.Run()

	//clean database
	mongoCtx, collection = pkg.GetMongoContext("users")
	collection.DeleteMany(mongoCtx, bson.M{})

	mongoCtx, collection = pkg.GetMongoContext("quotations")
	collection.DeleteMany(mongoCtx, bson.M{})

	mongoCtx, collection = pkg.GetMongoContext("applications")
	collection.DeleteMany(mongoCtx, bson.M{})

	os.Exit(testResult)
}

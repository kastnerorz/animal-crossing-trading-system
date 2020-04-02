package routers

import (
	"bytes"
	"encoding/json"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var userId string

func TestCreateUser(t *testing.T) {
	// normal
	body := []byte(`{"username":"zed2","nickname":"\u5f20\u8c46","password":"01db71ab8048f74a4b92c26ba77285ade0687ac192758e8185ad52701f649ef2","switchFriendCode":"SW-1234-1234-1234","jikeId":"\u5f20\u8c46"}`)
	r := PerformRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	assert.Equal(t, http.StatusCreated, r.Code)

	var res map[string]string
	err := json.Unmarshal([]byte(r.Body.String()), &res)
	userId, _ = res["id"]
	t.Log("UserId: " + userId)
	assert.Nil(t, err)

	// 400 username is already exist
	body = []byte(`{"username":"zed","nickname":"\u5f20\u8c46","password":"01db71ab8048f74a4b92c26ba77285ade0687ac192758e8185ad52701f649ef2","switchFriendCode":"SW-1234-1234-1234","jikeId":"\u5f20\u8c46"}`)
	r = PerformRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, r.Code)

}

func TestGetUser(t *testing.T) {

	t.Log("UserId: " + userId)
	// normal
	r := PerformRequest("GET", "/api/v1/users/"+userId, nil)
	assert.Equal(t, http.StatusOK, r.Code)

	var user models.User
	err := json.Unmarshal([]byte(r.Body.String()), &user)
	assert.Nil(t, err)
	assert.Equal(t, userId, user.ID.Hex())
	assert.Equal(t, "zed2", user.Username)
	assert.Equal(t, "张豆", user.Nickname)

	// 404
	r = PerformRequest("GET", "/api/v1/users/r"+userId, nil)
	assert.Equal(t, http.StatusNotFound, r.Code)
}

func TestGetMyInfo(t *testing.T) {
	// normal
	r := PerformRequestWithAuth("GET", "/api/v1/me", nil, ReviewerToken)
	assert.Equal(t, http.StatusOK, r.Code)

	var user models.User
	err := json.Unmarshal([]byte(r.Body.String()), &user)
	assert.Nil(t, err)
	assert.Equal(t, "zed", user.Username)
	assert.Equal(t, "张豆", user.Nickname)
	assert.Equal(t, "SW-1234-1234-1234", user.SwitchFriendCode)
	assert.Equal(t, "张豆", user.JikeID)
}

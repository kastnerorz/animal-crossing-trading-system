package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/models"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/testdata"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateApplication(t *testing.T) {
	t.Log("token", ApplicantToken)
	t.Log("id", testdata.QuotationId)
	body := []byte(`{"QuotationId":"` + testdata.QuotationId + `"}`)
	r := PerformRequestWithAuth("POST", "/api/v1/applications", bytes.NewBuffer(body), ApplicantToken)
	fmt.Println(r.Body)
	assert.Equal(t, http.StatusCreated, r.Code)
}

func TestGetMyApplicationsApplicant(t *testing.T) {
	r := PerformRequestWithAuth("GET", "/api/v1/applications?type=APPLY", nil, ApplicantToken)
	assert.Equal(t, http.StatusOK, r.Code)

	var applications []models.Application
	err := json.Unmarshal([]byte(r.Body.String()), &applications)
	application := applications[0]
	ApplicationId = application.ID
	assert.Nil(t, err)
	assert.Equal(t, testdata.QuotationId, application.QuotationId.Hex())
	assert.Equal(t, "FRIENDS", application.QuotationType)
	assert.Equal(t, "", application.PassCode)
	assert.Equal(t, "", application.SwitchFriendCode)
	assert.Equal(t, "PENDING", application.Status)
}

func TestUpdateApplication(t *testing.T) {
	body := []byte(`{"status":"ACCEPT","switchFriendCode":"SW-1234-1234-1234"}`)
	r := PerformRequestWithAuth("PUT", "/api/v1/applications/"+ApplicationId, bytes.NewBuffer(body), ReviewerToken)
	assert.Equal(t, http.StatusOK, r.Code)
}

func TestGetMyApplicationsReviewer(t *testing.T) {
	r := PerformRequestWithAuth("GET", "/api/v1/applications?type=REVIEW", nil, ReviewerToken)
	assert.Equal(t, http.StatusOK, r.Code)

	var applications []models.Application
	err := json.Unmarshal([]byte(r.Body.String()), &applications)
	application := applications[0]
	ApplicationId = application.ID
	assert.Nil(t, err)
	assert.Equal(t, testdata.QuotationId, application.QuotationId.Hex())
	assert.Equal(t, "FRIENDS", application.QuotationType)
	assert.Equal(t, "", application.PassCode)
	assert.Equal(t, "SW-1234-1234-1234", application.SwitchFriendCode)
	assert.Equal(t, "ACCEPT", application.Status)
}

// type PASS_CODE

func TestCreateApplicationPassCode(t *testing.T) {
	body := []byte(`{"QuotationId":"` + testdata.QuotationIdPassCode + `"}`)
	r := PerformRequestWithAuth("POST", "/api/v1/applications", bytes.NewBuffer(body), ApplicantToken)
	assert.Equal(t, http.StatusCreated, r.Code)
}

func TestGetMyApplicationsApplicantPassCode(t *testing.T) {
	r := PerformRequestWithAuth("GET", "/api/v1/applications?type=APPLY", nil, ApplicantToken)
	assert.Equal(t, http.StatusOK, r.Code)

	var applications []models.Application
	err := json.Unmarshal([]byte(r.Body.String()), &applications)
	application := applications[0]
	ApplicationId = application.ID
	assert.Nil(t, err)
	assert.Equal(t, testdata.QuotationIdPassCode, application.QuotationId.Hex())
	assert.Equal(t, "PASS_CODE", application.QuotationType)
	assert.Equal(t, "", application.PassCode)
	assert.Equal(t, "", application.SwitchFriendCode)
	assert.Equal(t, "PENDING", application.Status)
}

func TestUpdateApplicationPassCode(t *testing.T) {
	body := []byte(`{"status":"ACCEPT","passCode":"56HMS"}`)
	r := PerformRequestWithAuth("PUT", "/api/v1/applications/"+ApplicationId, bytes.NewBuffer(body), ReviewerToken)
	assert.Equal(t, http.StatusOK, r.Code)
}

func TestGetMyApplicationsReviewerPassCode(t *testing.T) {
	r := PerformRequestWithAuth("GET", "/api/v1/applications?type=REVIEW", nil, ReviewerToken)
	assert.Equal(t, http.StatusOK, r.Code)

	var applications []models.Application
	err := json.Unmarshal([]byte(r.Body.String()), &applications)
	application := applications[0]
	ApplicationId = application.ID
	assert.Nil(t, err)
	assert.Equal(t, testdata.QuotationIdPassCode, application.QuotationId.Hex())
	assert.Equal(t, "PASS_CODE", application.QuotationType)
	assert.Equal(t, "56HMS", application.PassCode)
	assert.Equal(t, "", application.SwitchFriendCode)
	assert.Equal(t, "ACCEPT", application.Status)
}

func TestDeleteApplication(t *testing.T) {
	r := PerformRequestWithAuth("DELETE", "/api/v1/applications/"+ApplicationId, nil, ApplicantToken)
	assert.Equal(t, http.StatusOK, r.Code)
}

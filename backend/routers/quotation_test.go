package routers

import (
	"bytes"
	"encoding/json"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/models"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/testdata"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var quotationId string

func TestCreateQuotation(t *testing.T) {
	body := []byte(`{"type":"SELL","price":40,"openType":"PASS_CODE","passCode":"56HMS","handlingFee":100000}`)
	r := PerformRequestWithAuth("POST", "/api/v1/quotations", bytes.NewBuffer(body), ReviewerToken)
	assert.Equal(t, http.StatusCreated, r.Code)
}

func TestGetQuotations(t *testing.T) {
	r := PerformRequest("GET", "/api/v1/quotations?type=SELL", nil)
	assert.Equal(t, http.StatusOK, r.Code)

	var quotations []models.Quotation
	err := json.Unmarshal([]byte(r.Body.String()), &quotations)
	quotation := quotations[0]
	quotationId = quotation.ID.Hex()
	assert.Nil(t, err)
	assert.Equal(t, 40, quotation.Price)
	assert.Equal(t, "PASS_CODE", quotation.OpenType)
	assert.Equal(t, "", quotation.PassCode)
	assert.Equal(t, 100000, quotation.HandlingFee)
}

func TestGetQuotation(t *testing.T) {
	r := PerformRequest("GET", "/api/v1/quotations/"+testdata.QuotationId, nil)
	assert.Equal(t, http.StatusOK, r.Code)

	var quotation models.Quotation
	err := json.Unmarshal([]byte(r.Body.String()), &quotation)
	assert.Nil(t, err)
	assert.Equal(t, 40, quotation.Price)
	assert.Equal(t, "FRIENDS", quotation.OpenType)
	assert.Equal(t, "", quotation.PassCode)
	assert.Equal(t, 100000, quotation.HandlingFee)
}

func TestUpdateQuotation(t *testing.T) {
	body := []byte(`{"price":90,"handlingFee":0,"openType":"FRIENDS"}`)
	r := PerformRequestWithAuth("PUT", "/api/v1/quotations/"+quotationId, bytes.NewBuffer(body), ReviewerToken)
	assert.Equal(t, http.StatusOK, r.Code)
}

func TestGetMyQuotation(t *testing.T) {
	r := PerformRequestWithAuth("GET", "/api/v1/my-quotations?type=SELL", nil, ReviewerToken)
	assert.Equal(t, http.StatusOK, r.Code)

	var quotations []models.Quotation
	err := json.Unmarshal([]byte(r.Body.String()), &quotations)
	assert.Nil(t, err)
	quotation := quotations[0]
	assert.Equal(t, 90, quotation.Price)
	assert.Equal(t, "FRIENDS", quotation.OpenType)
	assert.Equal(t, 0, quotation.HandlingFee)
}

func TestCreateQuotationPassCode(t *testing.T) {
	body := []byte(`{"type":"SELL","price":40,"openType":"PASS_CODE","passCode":"56HMS","handlingFee":100000}`)
	r := PerformRequestWithAuth("POST", "/api/v1/quotations", bytes.NewBuffer(body), ReviewerToken)
	assert.Equal(t, http.StatusCreated, r.Code)
}

func TestGetMyQuotationPassCode(t *testing.T) {
	r := PerformRequestWithAuth("GET", "/api/v1/my-quotations", nil, ReviewerToken)
	assert.Equal(t, http.StatusOK, r.Code)

	var quotations []models.Quotation
	err := json.Unmarshal([]byte(r.Body.String()), &quotations)
	quotation := quotations[0]
	assert.Nil(t, err)
	assert.Equal(t, 40, quotation.Price)
	assert.Equal(t, "PASS_CODE", quotation.OpenType)
	assert.Equal(t, 100000, quotation.HandlingFee)
}

func TestDeleteQuotation(t *testing.T) {
	r := PerformRequestWithAuth("DELETE", "/api/v1/quotations/"+quotationId, nil, ReviewerToken)
	assert.Equal(t, http.StatusOK, r.Code)
}

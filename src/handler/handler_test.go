package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/clear-street/backend-screening-parthingle/src/db"
	"github.com/clear-street/backend-screening-parthingle/src/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	for k := range db.AllTrades {
		delete(db.AllTrades, k)
	}
}
func TestTradesHandlerFuncHappyPath(t *testing.T) {
	defer cleanup()

	// GET /v1/trades
	reqGET, err := http.NewRequest("GET", "/v1/trades", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TradesHandlerFunc)

	handler.ServeHTTP(rr, reqGET)
	status := rr.Code
	assert.Equal(t, status, http.StatusOK, "Status code for GET /v1/trades should always be 200")

	// POST /v1/trades GoodPosts()
	reqPOST, err := http.NewRequest("POST", "/v1/trades", strings.NewReader(string(GoodPosts())))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, reqPOST)
	status = rr.Code
	assert.Equal(t, status, http.StatusOK, "Status code for POST GoodPosts() should be 200 on the first attempt")

}

func TestTradesHandlerFuncBadRequestPostsWrongFormat(t *testing.T) {
	defer cleanup()
	// POST /v1/trades GoodPosts()

	handler := http.HandlerFunc(TradesHandlerFunc)
	rr := httptest.NewRecorder()
	reqPOST, err := http.NewRequest("POST", "/v1/trades", strings.NewReader(string(BadRequestBadPriceTypePosts())))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, reqPOST)
	status := rr.Code
	assert.Equal(t, status, http.StatusBadRequest, "Status code for POST BadRequestBadPriceTypePosts should be 400")

}

func TestTradesHandlerFuncPostsMissingRequiredField(t *testing.T) {
	defer cleanup()
	// POST /v1/trades MissingRequiredJSONParseErrorPosts()

	handler := http.HandlerFunc(TradesHandlerFunc)
	rr := httptest.NewRecorder()
	reqPOST, err := http.NewRequest("POST", "/v1/trades", strings.NewReader(string(MissingRequiredJSONParseErrorPosts())))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, reqPOST)
	status := rr.Code
	// Bad or missing ticker type
	assert.Equal(t, status, http.StatusUnprocessableEntity, "Status code for POST MissingRequiredJSONParseErrorPosts should be 422")
}

func getParsedTradeObjects(b []byte) ([]model.TradeSubmitted, error) {
	t := []model.TradeSubmitted{}
	err := json.Unmarshal(b, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}
func TestTradeHandlerFuncLookupByIDAfterPostSuccess(t *testing.T) {
	defer cleanup()
	handler := http.HandlerFunc(TradesHandlerFunc)
	rr := httptest.NewRecorder()
	reqPOST, err := http.NewRequest("POST", "/v1/trades", strings.NewReader(string(GoodPosts())))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, reqPOST)
	status := rr.Code
	assert.Equal(t, status, http.StatusOK, "Should be a successful POST")

	trades, err := getParsedTradeObjects(rr.Body.Bytes())
	reqGET, err := http.NewRequest("GET", "/v1/trades/"+trades[0].TradeID, nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, reqGET)
	status = rr.Code
	assert.Equal(t, status, http.StatusOK)

}

func TestTradeHandlerFuncDeleteByIDAfterPostSuccess(t *testing.T) {
	defer cleanup()
	handler := http.HandlerFunc(TradesHandlerFunc)
	rr := httptest.NewRecorder()
	reqPOST, err := http.NewRequest("POST", "/v1/trades", strings.NewReader(string(GoodPosts())))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, reqPOST)
	status := rr.Code
	assert.Equal(t, status, http.StatusOK, "Should be a successful POST")

	trades, err := getParsedTradeObjects(rr.Body.Bytes())
	reqGET, err := http.NewRequest("DELETE", "/v1/trades/"+trades[0].TradeID, nil)
	assert.Nil(t, err)
	handler.ServeHTTP(rr, reqGET)
	status = rr.Code
	assert.Equal(t, status, http.StatusOK)

}

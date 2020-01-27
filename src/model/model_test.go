package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTradeToJSONValid(t *testing.T) {
	trade := Trade{ClientTradeID: "12345", Date: 20010101, Quantity: "10", Price: "5.67", Ticker: "PRTH"}
	json, err := trade.ToJSON()

	assert.Equal(t, `{"client_trade_id":"12345","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"}`,
		string(json), "Trade marshalling gone wrong")
	assert.Nil(t, err)

}

func TestTradeToJSONInvalid(t *testing.T) {
	trade := Trade{ClientTradeID: "12345", Date: 120010101, Quantity: "10", Price: "5.67", Ticker: "PRTH"}
	json, err := trade.ToJSON()

	assert.NotEqual(t, `{"client_trade_Wrong","JSONdate":"OUTPUT_120010101","quantity":"10","price":"5.67","ticker":"PRTH"}`,
		string(json), "Trade marshalling gone wrong")
	assert.Nil(t, err)
}

func TestFromJSONToTradeValid(t *testing.T) {
	json := []byte(`[{"client_trade_id":"12345","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"}]`)
	trade, err := FromJSON(json)

	assert.Equal(t, trade[0], Trade{ClientTradeID: "12345", Date: 20010101, Quantity: "10", Price: "5.67", Ticker: "PRTH"})
	assert.Nil(t, err)
}

func TestFromJSONToTradeInvalidBadQuantity(t *testing.T) {
	json := []byte(`[{"client_trade_id":"12345","date":21000101,"quantity":"1q0","price":"5.67","ticker":"PRTH"}]`)
	_, err := FromJSON(json)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "bad or missing quantity format")
}

func TestFromJSONToTradeInvalidBadDate(t *testing.T) {
	json := []byte(`[{"client_trade_id":"12345","date":2,"quantity":"12","price":"5.67","ticker":"PRTH"}]`)
	_, err := FromJSON(json)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "bad or missing date")
}

func TestFromJSONWithBadTypesFail(t *testing.T) {
	json1 := []byte(`[{"client_trade_id":"12345","date":20010101,"quantity":"10","price":5.67,"ticker":"PRTH"}]`)
	_, err := FromJSON(json1)
	assert.NotNil(t, err, "We should get a bad Price type error here")
	assert.Equal(t, err.Error(), "bad price type")

	json2 := []byte(`[{"client_trade_id":"12345","date":"20010101","quantity":"10","price":"5.67","ticker":"PRTH"}]`)
	_, err = FromJSON(json2)
	assert.NotNil(t, err, "We should get a bad Date type error here")
	assert.Equal(t, err.Error(), "bad date type")

	json3 := []byte(`[{"date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"}]`)
	_, err = FromJSON(json3)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "bad or missing client_trade_id", "json3 has missing client_trade_id")
}

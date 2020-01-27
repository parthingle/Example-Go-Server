package db

import (
	"encoding/json"
	"testing"

	"github.com/clear-street/backend-screening-parthingle/src/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	for k := range AllTrades {
		delete(AllTrades, k)
	}
}

func TestGetAfterInsertSuccess(t *testing.T) {
	defer cleanup()
	JSON := []byte(`[{"client_trade_id":"123456","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"}]`)

	key, err := AtomicInsertTradesFromJSONArray(JSON)
	assert.Nil(t, err, "We shouldn't be getting an error here")
	assert.Equal(t, key[0].ClientTradeID, "123456")

	v := AllTrades[key[0].TradeID]

	assert.Equal(t, key[0].TradeID, GenKey(v))

	expectedInternalTrade, _ := GetTradeByID(GenKey(v))
	expectedTradeJSON, _ := expectedInternalTrade.Trade.ToJSON()
	assert.Equal(t, JSON[1:len(JSON)-1], expectedTradeJSON, "The Trade JSONs should match")

	insertedTrade := AllTrades[GenKey(v)]
	assert.Equal(t, insertedTrade, expectedInternalTrade.Trade, "The Trade objects should match")
	assert.True(t, len(GenKey(v)) < 256, "Length of the key should be < 256")
}

func TestDeleteTradeByIDSuccess(t *testing.T) {
	defer cleanup()

	JSON := []byte(`[{"client_trade_id":"12345","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"}]`)

	key, err := AtomicInsertTradesFromJSONArray(JSON)
	assert.Nil(t, err, "We shouldn't be getting an error here")
	assert.Equal(t, key[0].ClientTradeID, "12345")
	v := AllTrades[key[0].TradeID]
	assert.Equal(t, key[0].TradeID, GenKey(v))

	_, exists := AllTrades[GenKey(v)]
	assert.True(t, exists, "This trade should exist in DB")

	err = DeleteTradeByID(GenKey(v))
	assert.Nil(t, err, "We shouldn't be getting an error here")

	_, exists = AllTrades[GenKey(v)]
	assert.False(t, exists, "This trade shouldn't exist in DB")
}

func TestDeleteTradeNonExistentIDFail(t *testing.T) {
	defer cleanup()

	JSON := []byte(`[{"client_trade_id":"12345","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"}]`)

	key, err := AtomicInsertTradesFromJSONArray(JSON)
	assert.Nil(t, err, "We shouldn't be getting an error here")
	assert.Equal(t, key[0].ClientTradeID, "12345")
	v := AllTrades[key[0].TradeID]
	assert.Equal(t, key[0].TradeID, GenKey(v))

	_, exists := AllTrades[GenKey(v)]
	assert.True(t, exists, "This trade should exist in DB")

	err = DeleteTradeByID("non-existent-key")
	assert.NotNil(t, err, "We should get a trade not found error here")
	assert.Equal(t, err.Error(), "trade not found")

}

func TestGetTradeNonExistentIDFail(t *testing.T) {
	defer cleanup()
	JSON := []byte(`[{"client_trade_id":"12345","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"}]`)

	key, err := AtomicInsertTradesFromJSONArray(JSON)
	assert.Nil(t, err, "We shouldn't be getting an error here")
	assert.Equal(t, key[0].ClientTradeID, "12345")
	v := AllTrades[key[0].TradeID]
	assert.Equal(t, key[0].TradeID, GenKey(v))

	_, exists := AllTrades[GenKey(v)]
	assert.True(t, exists)

	_, err = GetTradeByID("non-existent-key")
	assert.NotNil(t, err, "we should get an error here about a non existent trade")
	assert.Equal(t, err.Error(), "trade not found")
}

func TestInsertTradeWithDuplicateTickerFail(t *testing.T) {
	defer cleanup()

	JSON1 := []byte(`[{"client_trade_id":"12345","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"}]`)

	key, err := AtomicInsertTradesFromJSONArray(JSON1)
	assert.Nil(t, err, "We shouldn't be getting an error here")
	assert.Equal(t, key[0].ClientTradeID, "12345")
	v := AllTrades[key[0].TradeID]
	assert.Equal(t, key[0].TradeID, GenKey(v))

	JSON2 := []byte(`[{"client_trade_id":"23456","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"}]`)

	key, err = AtomicInsertTradesFromJSONArray(JSON2)
	assert.NotNil(t, err, "We should get an error here about an existing trade with this ticker")
	assert.Equal(t, err.Error(), "contains trade with already existing ticker or client ID")

	JSON3 := []byte(`[{"client_trade_id":"12345","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"},{"client_trade_id":"23456","date":20010102,"quantity":"20","price":"87.1","ticker":"PRTH"},{"client_trade_id":"34567","date":20010102,"quantity":"7","price":"5.4","ticker":"AMZN"}]`)

	key, err = AtomicInsertTradesFromJSONArray(JSON3)
	assert.NotNil(t, err, "We should get an error here about an existing trade with this ticker")
	assert.Equal(t, err.Error(), "contains trade with already existing ticker or client ID")

}

func TestUpdateExistingTradesFailThenSuccess(t *testing.T) {
	defer cleanup()
	JSONs := []byte(`[{"client_trade_id":"12345","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"},{"client_trade_id":"23456","date":20010102,"quantity":"20","price":"87.1","ticker":"AAPL"},{"client_trade_id":"34567","date":20010102,"quantity":"7","price":"5.4","ticker":"AMZN"}]`)

	_, err := AtomicInsertTradesFromJSONArray(JSONs)
	assert.Nil(t, err, "There should be no errors here")

	newTrade := []byte(`{"client_trade_id":"01234","date":20010102,"quantity":"10","price":"5.67","ticker":"QWER"}`)
	newTradeWithExistingTicker := []byte(`{"client_trade_id":"23456","date":20010102,"quantity":"10","price":"5.67","ticker":"AAPL"}`)
	_, err = UpdateExistingTrade(newTrade, "non-existent-trade-id")

	assert.NotNil(t, err, "We should get a trade not found error here")
	assert.Equal(t, err.Error(), "trade not found")

	key := GenKey(model.Trade{ClientTradeID: "12345", Date: 20010101, Quantity: "10", Price: "5.67", Ticker: "PRTH"})

	_, err = UpdateExistingTrade(newTradeWithExistingTicker, key)
	assert.NotNil(t, err, "We should get a existing ticker error here")
	assert.Equal(t, err.Error(), "contains trade with already existing ticker or client ID")

	ret, err := UpdateExistingTrade(newTrade, key)
	assert.Nil(t, err, "There should be no error here")
	trade := model.Trade{}
	err = json.Unmarshal(newTrade, &trade)

	assert.Equal(t, ret.Trade, trade)
	assert.NotEqual(t, key, ret.ID)
}

func TestGetAllTradesHappyPathSuccess(t *testing.T) {
	defer cleanup()
	JSONs := []byte(`[{"client_trade_id":"12345","date":20010101,"quantity":"10","price":"5.67","ticker":"PRTH"},{"client_trade_id":"23456","date":20010102,"quantity":"20","price":"87.1","ticker":"AAPL"},{"client_trade_id":"34567","date":20010102,"quantity":"7","price":"5.4","ticker":"AMZN"}]`)

	_, err := AtomicInsertTradesFromJSONArray(JSONs)
	assert.Nil(t, err, "There should be no errors here")

	allDbItems, _ := GetAllTrades()
	JSONsFromDB, err := json.Marshal(allDbItems)
	assert.Nil(t, err, "There should be no errors here")

	tradesFromJSONsFromDB, err := model.FromJSON(JSONs)
	assert.Nil(t, err, "There should be no errors here")

	internalTradesFromTest := []model.InternalTrade{}
	err = json.Unmarshal(JSONsFromDB, &internalTradesFromTest)
	assert.Nil(t, err, "There should be no errors here")

	tradesFromTest := []model.Trade{}

	for _, it := range internalTradesFromTest {
		tradesFromTest = append(tradesFromTest, it.Trade)
	}
	assert.Nil(t, err, "There should be no errors here")

	assert.Equal(t, tradesFromJSONsFromDB, tradesFromTest, "These values should be equal")
}

package db

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"

	"github.com/clear-street/backend-screening-parthingle/src/model"
)

// AllTrades is a mock DB as a key-value store
var AllTrades = map[string]model.Trade{}

// GetAllTrades ... used by HandleFunc GET /v1/trades
func GetAllTrades() ([]model.InternalTrade, error) {
	trades := []model.InternalTrade{}

	for k, v := range AllTrades {
		t := model.InternalTrade{ID: k, Trade: v}
		trades = append(trades, t)
	}

	return trades, nil
}

// GetTradeByID ...used by HandleFunc GET /v1/trades/{trade_id}
func GetTradeByID(id string) (model.InternalTrade, error) {
	ret := model.InternalTrade{}
	if val, ok := AllTrades[id]; ok {
		ret.Trade = val
		ret.ID = id
		return ret, nil
	}

	return ret, errors.New("trade not found")
}

// DeleteTradeByID ...used by HandleFunc DELETE /v1/trades/{trade_id}
func DeleteTradeByID(id string) error {
	if _, ok := AllTrades[id]; ok {
		delete(AllTrades, id)
		return nil
	}

	return errors.New("trade not found")
}

// GenKey is a deterministic way to generate internal db ID for lookups
func GenKey(t model.Trade) string {
	j, _ := json.Marshal(t)
	hash := md5.Sum(j)
	s := hex.EncodeToString(hash[:])

	return s[:int(math.Min(float64(len(s)), 255.9))]
}

func lookupByTickerAndClientID(t model.Trade) bool {
	for _, v := range AllTrades {
		if t.Ticker == v.Ticker {
			return true
		}
		if t.ClientTradeID == v.ClientTradeID {
			return true
		}
	}

	return false
}

// GetTradesFromJSONArraySafe ..used by AtomicInsertTradesFromJSONArray
func GetTradesFromJSONArraySafe(js []byte) ([]model.Trade, error) {
	trades, err := model.FromJSON(js)
	if err != nil {
		return trades, err
	}
	for _, t := range trades {
		if lookupByTickerAndClientID(t) {
			return trades, errors.New("contains trade with already existing ticker or client ID")
		}
	}
	// TODO: Use a set() here maybe
	for i, t := range trades {
		for j, q := range trades {
			if i != j && (t.ClientTradeID == q.ClientTradeID || t.Ticker == q.Ticker) {
				return trades, errors.New("contains trade with already existing ticker or client ID")
			}
		}
	}

	return trades, nil
}

// AtomicInsertTradesFromJSONArray ...used by HandleFunc POST /v1/trades
func AtomicInsertTradesFromJSONArray(ts []byte) ([]model.TradeSubmitted, error) {
	res := []model.TradeSubmitted{}
	trades, err := GetTradesFromJSONArraySafe(ts)
	if err != nil {
		return res, err
	}
	for _, t := range trades {
		tradeID := GenKey(t)
		AllTrades[tradeID] = t
		res = append(res, model.TradeSubmitted{ClientTradeID: t.ClientTradeID, TradeID: tradeID})
	}

	return res, nil
}

// UpdateExistingTrade ...used by HandleFunc PUT /v1/trades/
func UpdateExistingTrade(t []byte, tradeID string) (model.InternalTrade, error) {
	ret := model.InternalTrade{}
	if _, ok := AllTrades[tradeID]; !ok {
		return ret, errors.New("trade not found")

	}
	trade, err := model.FromJSON(t)
	if err != nil {
		return ret, err
	}
	temp := AllTrades[tradeID]
	delete(AllTrades, tradeID)
	if lookupByTickerAndClientID(trade[0]) {
		AllTrades[tradeID] = temp
		return ret, errors.New("contains trade with already existing ticker or client ID")
	}
	newTradeID := GenKey(trade[0])
	AllTrades[newTradeID] = trade[0]
	ret.Trade = trade[0]
	ret.ID = newTradeID

	return ret, nil
}

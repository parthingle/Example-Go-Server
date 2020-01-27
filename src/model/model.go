package model

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
)

// Trade ...Base trade details; common amongst all trade types
type Trade struct {
	ClientTradeID string `json:"client_trade_id"`
	Date          int32  `json:"date"`
	Quantity      string `json:"quantity"`
	Price         string `json:"price"`
	Ticker        string `json:"ticker"`
}

// InternalTrade ...Internal representation of trade including id
type InternalTrade struct {
	ID    string `json:"id"`
	Trade Trade  `json:"trade"`
}

// TradeSubmitted ...Submitted trade details
type TradeSubmitted struct {
	ClientTradeID string `json:"client_trade_id"`
	TradeID       string `json:"trade_id"`
}

// Error ...human readable description of error
type Error struct {
	Message string `json:"message"`
}

// ToJSON to be used for marshalling of Trade type
func (t Trade) ToJSON() ([]byte, error) {
	ToJSON, err := json.Marshal(t)
	if err != nil {
		return nil, errors.New("error parsing Trade to JSON")
	}
	return ToJSON, nil
}

func validTrade(trade Trade) (bool, error) {

	var validQuantity = regexp.MustCompile(`^[-]?[0-9]*\.?[0-9]+$`)
	var validPrice = regexp.MustCompile(`^[-]?[0-9]*\.?[0-9]+$`)

	if len(trade.ClientTradeID) < 1 || len(trade.ClientTradeID) > 256 {
		return false, errors.New("bad or missing client_trade_id")
	}

	if trade.Date < 20010101 || trade.Date > 21000101 {
		return false, errors.New("bad or missing date")
	}

	if !validQuantity.MatchString(trade.Quantity) {
		return false, errors.New("bad or missing quantity format")
	}

	if !validPrice.MatchString(trade.Price) {
		return false, errors.New("bad or missing price format")
	}

	if len(trade.Ticker) < 1 {
		return false, errors.New("bad or missing ticker format")
	}
	return true, nil
}

func parseBadMapType(m map[string]interface{}) (bool, string) {
	if v, ok := m["client_trade_id"]; ok {
		if reflect.TypeOf(v).String() != "string" {
			return true, "bad client_trade_id type"
		}
	}
	if v, ok := m["date"]; ok {
		if reflect.TypeOf(v).String() != "float64" {
			return true, "bad date type"
		}
	}
	if v, ok := m["quantity"]; ok {
		if reflect.TypeOf(v).String() != "string" {
			return true, "bad quantity type"
		}
	}
	if v, ok := m["price"]; ok {
		if reflect.TypeOf(v).String() != "string" {
			return true, "bad price type"
		}
	}
	if v, ok := m["ticker"]; ok {
		if reflect.TypeOf(v).String() != "string" {
			return true, "bad ticker type"
		}
	}
	return false, ""
}

// FromJSON to be used for unmarshalling of Trade type
func FromJSON(data []byte) ([]Trade, error) {

	trades := []Trade{}
	err := json.Unmarshal(data, &trades)

	if err != nil {
		t := Trade{}
		err = json.Unmarshal(data, &t)
		if err != nil { // Exclusive error territory. Now need to find out what type of error
			var emap []map[string]interface{}
			err = json.Unmarshal(data, &emap)
			if err != nil {
				return trades, errors.New("bad JSON format")
			}
			for _, m := range emap {
				typeerror, msg := parseBadMapType(m)
				if typeerror {
					return trades, errors.New(msg)
				}
			}
		} else {
			trades = append(trades, t)
		}
	}

	for _, trade := range trades {
		if valid, err := validTrade(trade); !valid {
			return trades, err
		}
	}
	return trades, nil
}

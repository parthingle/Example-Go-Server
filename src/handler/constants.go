package handler

// GoodPosts returns properly formatted json bytearray for 200 response
func GoodPosts() []byte {
	// 200
	return []byte(`[
		{
		  "client_trade_id": "T-50264430-bc41",
		  "date": 20200101,
		  "quantity": "100",
		  "price": "10.00",
		  "ticker": "AAPL"
		},
		{
			"client_trade_id": "Q-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "AMZN"
		},
		{
			"client_trade_id": "P-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "PRTH"
		}
	  ]`)
}

// BadRequestBadPriceTypePosts returns json with incorrect Price type
func BadRequestBadPriceTypePosts() []byte {
	// 400
	return []byte(`[
		{
			"client_trade_id": "T-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "AAPL"
		},
		{
			"client_trade_id": "Q-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "AMZN"
		},
		{
			"client_trade_id": "P-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": 1234,
			"ticker": "PRTH"
		}
	  ]`)
}

// BadRequestBadTickerTypePosts returns posts with Ticker as a float instead of string
func BadRequestBadTickerTypePosts() []byte {
	// 400
	return []byte(`[
		{
			"client_trade_id": "T-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "AAPL"
		},
		{
			"client_trade_id": "Q-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "AMZN"
		},
		{
			"client_trade_id": "P-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "1234",
			"ticker": 1.0
		}
	  ]`)
}

// MissingRequiredJSONParseErrorPosts returns posts where one of them has a missing Ticker field
func MissingRequiredJSONParseErrorPosts() []byte {
	// 422
	return []byte(`[
		{
			"client_trade_id": "T-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "AAPL"
		},
		{
			"client_trade_id": "Q-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "AMZN"
		},
		{
			"client_trade_id": "P-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00"
		}
	  ]`)

}

// DuplicateTickerFieldPosts returns posts with two of them having the same Ticker
func DuplicateTickerFieldPosts() []byte {
	// 422
	return []byte(`[
		{
		  "client_trade_id": "T-50264430-bc41",
		  "date": 20200101,
		  "quantity": "100",
		  "price": "10.00",
		  "ticker": "AAPL"
		},
		{
			"client_trade_id": "Q-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "PRTH"
		},
		{
			"client_trade_id": "P-50264430-bc41",
			"date": 20200101,
			"quantity": "100",
			"price": "10.00",
			"ticker": "PRTH"
		}
	  ]`)

}

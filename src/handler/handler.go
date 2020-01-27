package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/clear-street/backend-screening-parthingle/src/db"
	"github.com/clear-street/backend-screening-parthingle/src/model"
)

func writeJSON(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// TradesHandlerFunc ...handles GET and POST /v1/trades endpoint
func TradesHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		trades, err := db.GetAllTrades()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			e, _ := json.Marshal(model.Error{Message: err.Error()})
			fmt.Fprint(w, e)
			break
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		writeJSON(w, trades)
		break
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		submissions, err := db.AtomicInsertTradesFromJSONArray(body)
		if err != nil {
			errString := err.Error()
			if strings.Contains(errString, "bad JSON format") {
				w.WriteHeader(http.StatusBadRequest)

			} else if strings.Contains(errString, "bad or missing") {
				w.WriteHeader(http.StatusUnprocessableEntity)

			} else if strings.Contains(errString, "bad") {
				w.WriteHeader(http.StatusBadRequest)

			} else {
				w.WriteHeader(http.StatusInternalServerError)

			}
			writeJSON(w, model.Error{Message: errString})
			break
		}
		writeJSON(w, submissions)
		break
	}
}

// TradeHandlerFunc ...handles GET, DELETE, and PUT /v1/trades/ endpoint
func TradeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/trades/"):]
	switch method := r.Method; method {
	case http.MethodGet:
		trade, err := db.GetTradeByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, model.Error{Message: err.Error()})
			break
		}
		writeJSON(w, trade)
		break

	case http.MethodDelete:
		err := db.DeleteTradeByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, model.Error{Message: err.Error()})
		}
		break

	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		ret, err := db.UpdateExistingTrade(body, id)
		if err != nil {
			switch err.Error() {
			case "trade not found":
				w.WriteHeader(http.StatusNotFound)
				break
			default:
				w.WriteHeader(http.StatusInternalServerError)
				break
			}
			writeJSON(w, model.Error{Message: err.Error()})
			return
		}
		writeJSON(w, ret)
		break

	}

}

package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Instrument struct {
	// ExchangeToken   string `json:"exchange_token"`
	// Name            string `json:"name"`
	// LastPrice       string `json:"last_price"`
	// TickSize        string `json:"tick_size"`
	InstrumentToken int     `json:"instrument_token"`
	Tradingsymbol   string  `json:"tradingsymbol"`
	Expiry          string  `json:"expiry"`
	Strike          float64 `json:"strike"`
	LotSize         int     `json:"lot_size"`
	InstrumentType  string  `json:"instrument_type"`
	Segment         string  `json:"segment"`
	Exchange        string  `json:"exchange"`
}

func InstrumentsDetailsHandler(c echo.Context) error {

	// get query params
	var query struct {
		Instruments []string
		Tokens      []string
	}

	// creates query params binder that stops binding at first error
	err := echo.QueryParamsBinder(c).
		Strings("i", &query.Instruments).
		Strings("t", &query.Tokens).
		BindError() // returns first binding error

	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// process the query params
	queryInstruments := query.Instruments
	queryTokens := query.Tokens

	// initialize the database connection
	db := ConnectToDB()

	// tx is a database transaction
	var tx *gorm.DB

	// Query instruments by instrument names
	if len(queryInstruments) > 0 {
		instrumentsFound := []InstrumentsModel{}
		tx = db.Raw("SELECT * FROM api_instruments WHERE instrument IN (?)", queryInstruments).Scan(&instrumentsFound)
		if tx.Error != nil {
			return SendError(c, http.StatusBadRequest, tx.Error.Error())
		}
		CloseDB(db)

		responseData := map[string]InstrumentsModel{}
		for _, record := range instrumentsFound {
			responseData[record.Instrument] = record
		}
		return SendResponse(c, http.StatusOK, responseData)

		// Query instruments by instrument tokens
	} else if len(queryTokens) > 0 {
		tokensFound := []InstrumentsModel{}
		tx = db.Raw("SELECT * FROM api_instruments WHERE instrument_token IN (?)", queryTokens).Scan(&tokensFound)
		if tx.Error != nil {
			return SendError(c, http.StatusBadRequest, tx.Error.Error())
		}
		CloseDB(db)

		responseData := map[int]InstrumentsModel{}
		for _, record := range tokensFound {
			responseData[record.InstrumentToken] = record
		}
		return SendResponse(c, http.StatusOK, responseData)

		// send error if no query params are provided
	} else {
		return SendError(c, http.StatusBadRequest, "instruments or tokens query param is required")
	}
}

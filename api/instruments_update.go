package api

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// InstrumentsUpdateHandler updates the instruments
func InstrumentsUpdateHandler(c echo.Context) error {

	// Download csv file from the internet and save to instruments table
	// Download the CSV file
	instrumentsUrl := "https://api.kite.trade/instruments"
	resp, err := http.Get(instrumentsUrl)
	if err != nil {
		return err
	}

	// Close the response body on function exit
	defer resp.Body.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return SendError(c, http.StatusInternalServerError, err.Error())
	}

	// exchanges to update
	exchangesToUpdate := []string{"NSE", "NFO", "MCX"}

	// create a slice of InstrumentsModel
	instruments := []InstrumentsModel{}

	// iterate over the records and populate the struct
	for _, record := range records {

		exchange := record[11]
		if exchangesInSlice(exchange, exchangesToUpdate) {
			// exchange_token := record[1]
			// name := record[3]
			// last_price := record[4]
			// tick_size := record[7]
			instrument_token := record[0]
			tradingsymbol := record[2]
			expiry := record[5]
			strike := record[6]
			lot_size := record[8]
			instrument_type := record[9]
			segment := record[10]
			// exchange := record[11]

			// convert strings to data types
			instrumentToken, err := strconv.Atoi(instrument_token)
			if err != nil {
				return SendError(c, http.StatusInternalServerError, err.Error())
			}

			strikeValue, err := strconv.ParseFloat(strike, 64)
			if err != nil {
				return SendError(c, http.StatusInternalServerError, err.Error())
			}

			lotSize, err := strconv.Atoi(lot_size)
			if err != nil {
				return SendError(c, http.StatusInternalServerError, err.Error())
			}

			exTs := exchange + ":" + tradingsymbol

			// create an instance of InstrumentsModel
			instrument := InstrumentsModel{
				// ExchangeToken:   exchangeToken,
				// Name:            name,
				// LastPrice: last_price,
				// TickSize:        tick_size,
				Instrument:      exTs,
				InstrumentToken: instrumentToken,
				Exchange:        exchange,
				Tradingsymbol:   tradingsymbol,
				Expiry:          expiry,
				Strike:          strikeValue,
				LotSize:         lotSize,
				InstrumentType:  instrument_type,
				Segment:         segment,
			}

			// append the instrument to the slice
			instruments = append(instruments, instrument)
		}
	}

	// fmt.Println("Total instruments to update: ", len(instruments))

	// initialize the database connection
	db := ConnectToDB()

	// truncate table before inserting new records
	db.Exec("TRUNCATE TABLE api_instruments RESTART IDENTITY CASCADE;")
	db.Exec("VACUUM;")

	// insert records in batches
	tx := db.CreateInBatches(&instruments, 1000)
	rowsInserted := tx.RowsAffected

	// send response
	responseData := map[string]interface{}{
		"exchanges":   exchangesToUpdate,
		"instruments": rowsInserted,
		"updated_at":  time.Now().Format("2006-01-02 15:04:05"),
	}

	// close the database connection
	// CloseDB(db)

	return SendResponse(c, http.StatusOK, responseData)
}

func exchangesInSlice(exchange string, exchangesToUpdate []string) bool {
	for _, e := range exchangesToUpdate {
		if e == exchange {
			return true
		}
	}
	return false
}

package api

import (
	"fmt"

	"gorm.io/gorm"
)

// InstrumentInfo represent information about an instrument
type InstrumentInfo struct {
	Exchange        string `json:"exchange"`
	Tradingsymbol   string `json:"tradingsymbol"`
	Instrument      string `json:"instrument"`
	InstrumentToken uint32 `json:"instrument_token"`
}

// InstrumentsMap is a map of instrument to InstrumentInfo
type InstrumentsMap map[string]InstrumentInfo

// InstrumentTokensMap is a map of instrument_token to InstrumentInfo
type InstrumentTokensMap map[uint32]InstrumentInfo

func GetInstrumentsMap(instruments []string) *InstrumentsMap {

	// initialize the database connection
	db := ConnectToDB()

	// tx is a database transaction
	var tx *gorm.DB

	instrumentsFound := []InstrumentsModel{}
	queryInstruments := instruments
	tx = db.Raw("SELECT exchange, tradingsymbol, instrument_token, instrument FROM api_instruments WHERE instrument IN (?)", queryInstruments).Scan(&instrumentsFound)
	if tx.Error != nil {
		fmt.Println(tx.Error.Error())
	}

	defer CloseDB(db)

	// Get the instrument id map
	var instrumentsMap = make(InstrumentsMap)
	for _, record := range instrumentsFound {
		instrumentsMap[record.Instrument] = InstrumentInfo{
			Exchange:        record.Exchange,
			Tradingsymbol:   record.Tradingsymbol,
			Instrument:      record.Instrument,
			InstrumentToken: uint32(record.InstrumentToken),
		}
	}

	return &instrumentsMap

}

func GetInstrumentTokensMap(instrumentsMap *InstrumentsMap) *InstrumentTokensMap {
	var instrumentTokensMap = make(InstrumentTokensMap)
	for _, instrumentInfo := range *instrumentsMap {
		instrumentTokensMap[instrumentInfo.InstrumentToken] = instrumentInfo
	}
	return &instrumentTokensMap
}

// GetInstrumentTokens returns the slice of instrument tokens in the instruments map
func GetInstrumentTokens(instrumentsMap *InstrumentsMap) []uint32 {
	instrumentTokens := []uint32{}
	for _, instrumentInfo := range *instrumentsMap {
		instrumentTokens = append(instrumentTokens, instrumentInfo.InstrumentToken)
	}
	return instrumentTokens
}

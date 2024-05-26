package api

import "time"

// Users struct represents the user model
// fields instrument_token,exchange_token,tradingsymbol,name,last_price,expiry,strike,tick_size,lot_size,instrument_type,segment,exchange

type InstrumentsModel struct {
	// Name            string    `gorm:"type:varchar(100)"`
	// LastPrice       float64   `gorm:"type:real"`
	// TickSize        float64   `gorm:"type:real"`
	// ExchangeToken   int       `gorm:"type:integer"`
	ID              uint      `gorm:"primarykey" json:"-"`
	Instrument      string    `gorm:"type:varchar(120);index;unique" json:"instrument"`
	InstrumentToken int       `gorm:"type:integer;index;unique;not null" json:"instrument_token"`
	Exchange        string    `gorm:"type:varchar(20)" json:"exchange"`
	Tradingsymbol   string    `gorm:"type:varchar(100)" json:"tradingsymbol"`
	Expiry          string    `gorm:"type:varchar(10)" json:"expiry"`
	Strike          float64   `gorm:"type:real" json:"strike"`
	LotSize         int       `gorm:"type:integer" json:"lot_size"`
	InstrumentType  string    `gorm:"type:varchar(20)" json:"instrument_type"`
	Segment         string    `gorm:"type:varchar(20)" json:"segment"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"-"`
}

// TableName returns the table name of the instruments model
func (i *InstrumentsModel) TableName() string {
	return "api_instruments"
}

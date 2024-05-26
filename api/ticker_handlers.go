package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	kiteconnect "github.com/nsvirk/gokiteticker/kiteconnect"
	kitemodels "github.com/nsvirk/gokiteticker/models"
	kiteticker "github.com/nsvirk/gokiteticker/ticker"
)

var (
	ticker                    *kiteticker.Ticker
	tickerStatus              string
	STATUS_STARTED            string = "started"
	STATUS_STOPPED            string = "stopped"
	tickerUserid              string
	tickerInstruments         []string
	tickerInstrumentsMap      *InstrumentsMap
	tickerInstrumentTokensMap *InstrumentTokensMap
	tickerTokens              []uint32
	tickerPublishedTicks      int
	tickerChannel             string
	// tickerTokens []uint32 = append([]uint32{}, 256265, 264969, 5633, 779521, 408065, 738561, 895745)
)

type TickerStartInput struct {
	Instruments []string `json:"instruments" form:"instruments" query:"instruments"`
}

// TickerStopHandler stops the ticker
func TickerStopHandler(c echo.Context) error {

	// Check if ticker is not started
	if tickerStatus != STATUS_STARTED {
		return SendError(c, http.StatusBadRequest, "Ticker is not started")
	}

	// --------------------------------------------------------------
	// Stop the ticker
	// --------------------------------------------------------------
	ticker.Unsubscribe(tickerTokens)
	// --------------------------------------------------------------

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("Channel: ", tickerChannel)
	fmt.Println("Unsubscribing from", len(tickerTokens), "instrument tokens")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("")

	// Set the ticker status
	tickerStatus = STATUS_STOPPED

	// Send response
	respDate := map[string]interface{}{
		"user_id":         tickerUserid,
		"published_ticks": tickerPublishedTicks,
		"stopped_at":      time.Now().Format("2006-01-02 15:04:05"),
	}
	return SendResponse(c, http.StatusOK, respDate)
}

// TickerStartHandler starts the ticker
func TickerStartHandler(c echo.Context) error {

	// Check if ticker is already started
	if tickerStatus == STATUS_STARTED {
		return SendError(c, http.StatusBadRequest, "Ticker is already started")
	}

	// Get form inputs
	requestJson := new(TickerStartInput)
	err := c.Bind(requestJson)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err.Error())
	}

	// Get the instruments from the request
	tickerInstruments = requestJson.Instruments
	if tickerInstruments == nil {
		return SendError(c, http.StatusBadRequest, "`instruments` are required")
	}
	// Get the user_id and enctoken from the context
	userId := c.Get("user_id").(string)
	enctoken := c.Get("enctoken").(string)
	tickerUserid = userId // for global access

	// Get ticker instruments map, instrument tokens map and tokens
	tickerInstrumentsMap = GetInstrumentsMap(tickerInstruments)
	tickerInstrumentTokensMap = GetInstrumentTokensMap(tickerInstrumentsMap)
	tickerTokens = GetInstrumentTokens(tickerInstrumentsMap)
	tickerChannel = fmt.Sprintf("CH:API:TICKS:%s", tickerUserid)

	// --------------------------------------------------------------
	// Start the ticker
	// --------------------------------------------------------------
	// Create new Kite ticker instance
	ticker = kiteticker.New(userId, enctoken)

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)

	// Start the tick server
	go ticker.Serve()
	// --------------------------------------------------------------

	// Set the ticker status
	tickerStatus = STATUS_STARTED

	// Send response
	respDate := map[string]interface{}{
		"user_id":         userId,
		"instruments_ct":  len(tickerInstruments),
		"instruments_map": tickerInstrumentsMap,
		"started_at":      time.Now().Format("2006-01-02 15:04:05"),
	}

	return SendResponse(c, http.StatusOK, respDate)
}

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	// fmt.Println("Connected")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("Channel: ", tickerChannel)
	fmt.Println("Subscribing to", len(tickerTokens), "instrument tokens")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("")
	err := ticker.Subscribe(tickerTokens)
	if err != nil {
		fmt.Println("Connect: ", err)
	}
	// Set subscription mode for given list of tokens
	// Default mode is Quote
	err = ticker.SetMode(kiteticker.ModeFull, tickerTokens)
	// err = ticker.SetMode(kiteticker.ModeLTP, tickerTokens)
	if err != nil {
		fmt.Println("Connect: ", err)
	}
}

// Tick represents a single packet in the market feed.
type CustomTick struct {
	Tick          kitemodels.Tick
	Instrument    string
	Exchange      string
	Tradingsymbol string
}

// Triggered when tick is recevived
func onTick(tick kitemodels.Tick) {
	// Increment the tickerPublishedTicks
	tickerPublishedTicks = tickerPublishedTicks + 1

	// Convert to CustomTick
	instrumentInfo := (*tickerInstrumentTokensMap)[tick.InstrumentToken]
	instrument := instrumentInfo.Instrument
	exchange := instrumentInfo.Exchange
	tradingsymbol := instrumentInfo.Tradingsymbol
	// CustomTick
	customTick := CustomTick{
		Tick:          tick,
		Instrument:    instrument,
		Exchange:      exchange,
		Tradingsymbol: tradingsymbol,
	}

	// fmt.Println("--------------------------------------------------------------")
	// tickJsonIndented, err := json.MarshalIndent(customTick, "", " \t")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(tickJsonIndented))
	// fmt.Println("--------------------------------------------------------------")

	// Convert to JSON
	tickJson, err := json.Marshal(customTick)
	if err != nil {
		fmt.Println(err)
	}

	// publish the tick
	go PublishTicks([]byte(tickJson))

}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	fmt.Printf("Order: %s", order.OrderID)
}

// PublishTicks publishes the ticks to the clients
func PublishTicks(payload []byte) {
	// Publish ticks to the clients
	rdb := ConnectToRedis()
	ctx := context.Background()
	err := rdb.Publish(ctx, tickerChannel, payload).Err()
	if err != nil {
		fmt.Println("Error publising ticks: ", err.Error())
	}
}

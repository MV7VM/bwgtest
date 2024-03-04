package models

type DB struct {
	Ticker string
	Prices []byte
}

type Db struct {
	Ticker string   `db:"ticker"`
	Prices []Prices `db:"prices"`
}

type Prices struct {
	Prices []byte `json:"Prices"`
	Date   string `json:"Date"`
}

type NewTicker struct {
	Ticker string
}

type TicketDifference struct {
	Ticker     string
	Price      float32
	Difference float32
}

type TickerInfo struct {
	Ticker   string
	DateFrom string
	DateTo   string
}

type TickerInfoResponse struct {
	Info []TickerResponse
}

type TickerResponse struct {
	Ticker string `json:"symbol"`
	Price  string `json:"price"`
}

package common

import (
	"bwg/internal/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

func Trim(tickers []string) string {
	separator := string('"') + string(',') + string('"')
	result := strings.Join(tickers, separator)
	slog.Info("TRIM" + result)
	return "[" + string('"') + result + string('"') + "]"
}

func UpdateJson(oldPrice []models.Prices, update []byte, date string) []byte {
	//slog.Info("Update Json oldPrice", oldPrice)
	var newPrice models.Db
	newPrice.Prices = append(newPrice.Prices, oldPrice...)
	newPrice.Prices = append(newPrice.Prices, models.Prices{Prices: string(update), Date: date})
	//fmt.Println(newPrice.Prices)
	//fmt.Println(newPrice.Prices)
	//fmt.Println(newPrice.Prices)
	//price = append(price, byte(10))
	//price = append(price, update...)
	//rawDataOut, err := json.MarshalIndent(&settings, "", "  ")
	rawDataOut, err := json.Marshal(&newPrice.Prices)
	slog.Info("JSON out of UpdateJson", rawDataOut)
	if err != nil {
		slog.Error("Error in UpdateJson cann't marshal:", err)
	}
	return rawDataOut
}

func PriceDifference(ticker, dateFrom, dateTo string, prices []models.Prices) models.TicketDifference {
	dFrom, err := time.Parse("2006-01-02 15:04:05", dateFrom)
	if err != nil {
		slog.Error("Error in time from:", err)
	}
	dTo, err := time.Parse("2006-01-02 15:04:05", dateTo)
	if err != nil {
		slog.Error("Error in time from:", err)
	}
	var priceIn, priceOut string
	for _, v := range prices {
		dIn, err := time.Parse("2006-01-02 15:04:05", v.Date)
		if err != nil {
			slog.Error("Error in time from:", err)
		}
		fmt.Println(dIn.Sub(dFrom))
		if dIn.Sub(dFrom) <= 0 {
			priceIn = v.Prices
		} else if dTo.Sub(dIn) <= 0 {
			priceOut = v.Prices
			break
		}
	}
	fmt.Println(priceIn, priceOut)
	pOut, _ := strconv.ParseFloat(priceOut, 64)
	pIn, _ := strconv.ParseFloat(priceIn, 64)
	slog.Info("return Getter", models.TicketDifference{Ticker: "", Price: float32(pOut), Difference: float32((pIn - pOut) / pIn)})
	return models.TicketDifference{Ticker: ticker, Price: float32(pOut), Difference: float32((pOut - pIn) / pIn)}
}

//func main() {
//	tickers := []string{"USDTBTC", "BTNUSDT"}
//	fmt.Println(Trim(tickers))
//}

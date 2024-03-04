package common

import (
	"bwg/internal/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

func Trim(tickers []string) string {
	separator := string('"') + string(',') + string('"')
	result := strings.Join(tickers, separator)
	slog.Info("TRIM" + result)
	return "[" + string('"') + result + string('"') + "]"
}

func UpdateJson(oldPrice models.DB, update []byte, date string) []byte {
	slog.Info("Update Json oldPrice", oldPrice)
	var newPrice models.Db
	newPrice.Prices = append(newPrice.Prices, models.Prices{Prices: oldPrice.Prices, Date: date})
	newPrice.Prices = append(newPrice.Prices, models.Prices{Prices: update, Date: date})
	fmt.Println(newPrice.Prices)
	fmt.Println(newPrice.Prices)
	fmt.Println(newPrice.Prices)
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

//func main() {
//	tickers := []string{"USDTBTC", "BTNUSDT"}
//	fmt.Println(Trim(tickers))
//}

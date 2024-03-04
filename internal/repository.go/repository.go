package repository_go

import (
	"bwg/config"
	"bwg/database"
	"bwg/internal/common"
	"bwg/internal/models"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io/ioutil"
	"log/slog"
	"net/http"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(cfg config.Config) *Repository {
	pool, err := database.ConnectDB(cfg)
	if err != nil {
		slog.Error("error in create pool connects: ", err)
		return nil
	}
	return &Repository{pool: pool}
}

func (r *Repository) NewTicker(ticker string) error {
	_, err := r.pool.Exec(context.Background(), "INSERT INTO tickers VALUES ($1)", ticker)
	if err != nil {
		slog.Error("Fail to exec new ticker: ", ticker, " error: ", err)
		return err
	}
	return nil
}

func (r *Repository) GetPriceDifference(info models.TickerInfo) (models.TicketDifference, error) {
	var TicketDifference models.TicketDifference
	rows, _ := r.pool.Query(context.Background(), "SELECT * FROM tickers WHERE ticker=($1);", info.Ticker)
	slog.Info("GetPriceDifference Row:", rows)
	tickerInfo, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.DB])
	if err != nil {
		slog.Error("fail to CollectOneRow tickerInfo")
		return TicketDifference, err
	}
	var prices []models.Prices
	err = json.Unmarshal(tickerInfo.Prices, &prices)
	slog.Info("GetPriceDifference prices:", prices)
	return common.PriceDifference(info.Ticker, info.DateFrom, info.DateTo, prices), nil
}

func (r *Repository) CurrentPrice() error {
	rows, err := r.pool.Query(context.Background(), "SELECT ticker FROM tickers")
	if err != nil {
		return err
	}
	defer rows.Close()
	var tickers []string
	var newPrices models.TickerInfoResponse
	for rows.Next() {
		var ticker []byte
		err := rows.Scan(&ticker)
		slog.Info("Row:", ticker)
		if err != nil {
			slog.Error("Error in Scan tickers:", err)
		}
		tickers = append(tickers, string(ticker))
		slog.Info("CuttentPrice all tickers", tickers)
	}
	newPrices, err = r.GetPrice(tickers)
	err = database.UpdateDb(newPrices, r.pool, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		slog.Error("Error in UpdateDb: ", err)
		return err
	}
	if err = rows.Err(); err != nil {
		slog.Error("Something goes wrong in CurrentPrice rows: ", err)
		return err
	}
	return nil
}

func (r *Repository) GetPrice(tickers []string) (models.TickerInfoResponse, error) {
	//https://www.binance.com/api/v3/ticker/price?symbol=BNBBTC
	//https://www.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22,%22BNBUSDT%22%5D
	//resp, err := http.Get("https://www.binance.com/api/v3/ticker/price?symbols=%5B%22BTCUSDT%22,%22BNBUSDT%22%5D")
	resp, err := http.Get("https://www.binance.com/api/v3/ticker/price?symbols=" + common.Trim(tickers))
	if err != nil {
		slog.Error("error in Get: ", err)
		return models.TickerInfoResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	var m []models.TickerResponse
	err = json.Unmarshal(body, &m)
	if err != nil {
		slog.Error("Fail to unmarshal:", err)
	}
	slog.Info("GetPrice unmarshal", string(body))
	slog.Info("GetPrice unmarshal->model", models.TickerInfoResponse{Info: m})
	return models.TickerInfoResponse{Info: m}, nil
}

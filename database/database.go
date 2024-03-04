package database

import (
	"bwg/config"
	"bwg/internal/common"
	"bwg/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

func ConnectDB(cfg config.Config) (pool *pgxpool.Pool, err error) {
	//var cfg config.Config
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	DbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	pool, err = pgxpool.New(context.Background(), DbUrl)
	if err != nil {
		log.Error("Unable to connect to database: %v\n", err)
		//os.Exit(1)
		return nil, err
	}
	//defer conn.Close(context.Background())
	return pool, nil
}

func UpdateDb(newPrices models.TickerInfoResponse, pool *pgxpool.Pool, date string) error {
	//slog.Info("UpdateDb newPrices:", newPrices.Info)
	//var users []*User
	//pgxscan.Select(ctx, db, &users, `SELECT id, name, email, age FROM users`)
	for i, v := range newPrices.Info {
		slog.Info("UpdeteDb i,v:", i, v.Ticker)
		var tickerInfo models.DB
		var err error
		rows, _ := pool.Query(context.Background(), "SELECT * FROM tickers WHERE ticker=($1);", v.Ticker)
		slog.Info("UpdeteDb Row:", rows)
		tickerInfo, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[models.DB])
		if err != nil {
			slog.Error("fail to CollectOneRow tickerInfo")
			return err
		}
		slog.Info("UpdateDb CollectOneRow to tickerInfo", tickerInfo)
		fmt.Println(tickerInfo.Prices, string(tickerInfo.Prices))
		//var OldPrice []models.Prices

		var prices []models.Prices
		err = json.Unmarshal(tickerInfo.Prices, &prices)
		if err != nil {
			slog.Error("fail to inmarshal:", err)
		}
		//err = json.Unmarshal(tickerInfo.Prices, &OldPrice)
		//if err != nil {
		//	slog.Error("fail to unmarshal tickerInfo to oldPrice", err)
		//}
		slog.Info("UpdateDb Marshal to OldPrice", prices)
		//tickerInfo.Prices = OldPrice
		//err := row.Scan(&tickerInfo)
		//slog.Info("UpdeteDb Row->tickerInfo:", tickerInfo)

		NewPricesJson := common.UpdateJson(prices, []byte(v.Price), date)
		slog.Info("tickerInfo to commit:", tickerInfo.Ticker, NewPricesJson)
		_, err = pool.Exec(context.Background(), "UPDATE tickers SET prices=($1) WHERE ticker=($2);", NewPricesJson, v.Ticker)
		if err != nil {
			slog.Error("fail to update tickerInfo in DB ticker="+v.Ticker, err)
			return err
		}
	}
	return nil
}

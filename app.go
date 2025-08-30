package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// FetchPairs fetches the trading pairs from Binance.
func (a *App) FetchPairs(marketType string) ([]string, error) {
	client := binance.NewClient("", "")
	var symbols []string

	if marketType == "spot" {
		res, err := client.NewExchangeInfoService().Do(context.Background())
		if err != nil {
			return nil, err
		}
		for _, s := range res.Symbols {
			symbols = append(symbols, s.Symbol)
		}
	} else {
		fClient := binance.NewFuturesClient("", "")
		res, err := fClient.NewExchangeInfoService().Do(context.Background())
		if err != nil {
			return nil, err
		}
		for _, s := range res.Symbols {
			symbols = append(symbols, s.Symbol)
		}
	}
	return symbols, nil
}

// DownloadOHLC downloads OHLC data for a given symbol, timeframe and market type.
func (a *App) DownloadOHLC(symbol, timeframe, marketType, startTimeStr, endTimeStr, outputFormat string) (string, error) {
	var _ futures.Kline // Prevent "imported and not used" error

	var startTime, endTime, initialStartTime int64
	var err error

	if startTimeStr != "" {
		t, err := time.Parse("2006-01-02", startTimeStr)
		if err == nil {
			startTime = t.UnixMilli()
			initialStartTime = startTime
		}
	}
	if endTimeStr != "" {
		t, err := time.Parse("2006-01-02", endTimeStr)
		if err == nil {
			endTime = t.UnixMilli()
		}
	} else {
		// If no end time, use now
		endTime = time.Now().UnixMilli()
	}

	totalTimeRange := float64(endTime - initialStartTime)

	var allKlines []*binance.Kline

	for {
		var klines []*binance.Kline
		if marketType == "spot" {
			client := binance.NewClient("", "")
			service := client.NewKlinesService().Symbol(symbol).Interval(timeframe).Limit(1000)
			if startTime != 0 {
				service.StartTime(startTime)
			}
			if endTime != 0 {
				service.EndTime(endTime)
			}
			klines, err = service.Do(a.ctx)
		} else {
			fClient := binance.NewFuturesClient("", "")
			service := fClient.NewKlinesService().Symbol(symbol).Interval(timeframe).Limit(1500)
			if startTime != 0 {
				service.StartTime(startTime)
			}
			if endTime != 0 {
				service.EndTime(endTime)
			}
			futuresKlines, err := service.Do(a.ctx)
			if err == nil {
				for _, k := range futuresKlines {
					klines = append(klines, &binance.Kline{
						OpenTime:                 k.OpenTime,
						Open:                     k.Open,
						High:                     k.High,
						Low:                      k.Low,
						Close:                    k.Close,
						Volume:                   k.Volume,
						CloseTime:                k.CloseTime,
						QuoteAssetVolume:         k.QuoteAssetVolume,
						TradeNum:                 k.TradeNum,
						TakerBuyBaseAssetVolume:  k.TakerBuyBaseAssetVolume,
						TakerBuyQuoteAssetVolume: k.TakerBuyQuoteAssetVolume,
					})
				}
			}
		}

		if err != nil {
			return "", err
		}

		if len(klines) == 0 {
			break
		}

		allKlines = append(allKlines, klines...)
		startTime = klines[len(klines)-1].CloseTime + 1

		if endTime != 0 && startTime > endTime {
			break
		}

		// Calculate and emit progress
		if totalTimeRange > 0 {
			progress := (float64(startTime-initialStartTime) / totalTimeRange) * 100
			if progress > 100 {
				progress = 100
			}
			runtime.EventsEmit(a.ctx, "download_progress", int(progress))
		}

		// Add a small delay to avoid hitting API rate limits
		time.Sleep(200 * time.Millisecond)
	}

	fileName := fmt.Sprintf("%s-%s-%s.%s", symbol, marketType, timeframe, outputFormat)
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if outputFormat == "csv" {
		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		header := []string{"Date", "Open", "High", "Low", "Close", "Volume"}
		if err := writer.Write(header); err != nil {
			return "", err
		}

		// Write data
		for _, kline := range allKlines {
			timestamp, _ := strconv.ParseInt(fmt.Sprintf("%d", kline.OpenTime), 10, 64)
			t := time.Unix(timestamp/1000, 0)
			record := []string{
				t.Format("2006-01-02 15:04:05"),
				kline.Open,
				kline.High,
				kline.Low,
				kline.Close,
				kline.Volume,
			}
			if err := writer.Write(record); err != nil {
				return "", err
			}
		}
	} else {
		encoder := json.NewEncoder(file)
		err = encoder.Encode(allKlines)
		if err != nil {
			return "", err
		}
	}

	runtime.EventsEmit(a.ctx, "download_progress", 100)
	return fmt.Sprintf("Saved %d records to %s", len(allKlines), fileName), nil
}

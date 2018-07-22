package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"strconv"
	"strings"

	"github.com/go-resty/resty"
	"github.com/ppincak/rysen/api"
	"github.com/ppincak/rysen/binance/model"
	"github.com/ppincak/rysen/core"
	"github.com/ppincak/rysen/monitor"
	log "github.com/sirupsen/logrus"
)

var _ monitor.Reporter = (*Client)(nil)

type Client struct {
	url     string
	secret  *api.Secret
	metrics *core.ApiMetrics
}

func NewClient(url string, secret *api.Secret) *Client {
	return &Client{
		url:     url,
		secret:  secret,
		metrics: core.NewApiMetrics(),
	}
}

func (client *Client) Statistics() []*monitor.Statistic {
	return []*monitor.Statistic{
		client.metrics.ToStatistic("binanceClientCalls"),
	}
}

func (client *Client) assembleUrl(url string) string {
	var urlBuilder strings.Builder
	urlBuilder.WriteString(client.url)
	urlBuilder.WriteString(url)
	return urlBuilder.String()
}

func (client *Client) addApiKey(request *resty.Request) *resty.Request {
	if client.secret == nil {
		panic("Api secrets are mssing is missing")
	}
	return request.SetHeader("X-MBX-APIKEY", client.secret.ApiKey)
}

func (client *Client) signQuery(request *resty.Request, value []byte) *resty.Request {
	if client.secret == nil {
		panic("Api secrets are mssing is missing")
	}
	mac := hmac.New(sha256.New, []byte(client.secret.SecretKey))
	mac.Write(value)
	hash := mac.Sum(nil)

	return request.SetQueryParam("signature", string(hash))
}

func (client *Client) baseGetCall(url string, queryParams map[string]string) (api.ApiResponse, error) {
	resp, err := resty.R().
		SetQueryParams(queryParams).
		Get(client.assembleUrl(url))

	apiResp, err := client.handleResponse(resp, err)
	defer client.handleMetrics(err)
	return apiResp, err
}

func (client *Client) baseQueryCall(url string, symbol string, limit uint32) (api.ApiResponse, error) {
	queryParams := make(map[string]string)
	if symbol != "" {
		queryParams["symbol"] = symbol
	}
	if limit != 0 {
		queryParams["limit"] = strconv.FormatUint(uint64(limit), 10)
	}

	return client.baseGetCall(url, queryParams)
}

func (client *Client) handleResponse(resp *resty.Response, err error) (api.ApiResponse, error) {
	if err != nil {
		return nil, err
	}

	switch statusCode := resp.StatusCode(); {
	case statusCode == 429:
		log.Error("Request limit reached")
		fallthrough
	case statusCode > 201:
		log.Error("Request failed with statusCode %d", resp.StatusCode())

		return nil, api.NewError("Request failed with status code %d", resp.StatusCode())
	}

	m, err := api.Unmarshall(resp.Body())
	return m, err
}

func (client *Client) handleMetrics(err error) {
	if err != nil {
		client.metrics.Inc(core.SuccessfullCalls)
	} else {
		client.metrics.Inc(core.FailedCalls)
	}
}

func (client *Client) ExchangeInfo() (*model.ExchangeInfo, error) {
	resp, err := resty.R().Get(client.assembleUrl(Endpoints.ExchangeInfo))
	if err != nil {
		return nil, err
	}

	var exchangeInfo model.ExchangeInfo
	err = api.UnmarshallAs(resp.Body(), &exchangeInfo)
	return &exchangeInfo, err
}

func (client *Client) OrderBook(symbol string, limit uint32) (api.ApiResponse, error) {
	return client.baseQueryCall(Endpoints.OrderBook, symbol, limit)
}

func (client *Client) OrderBookTicker(symbol string) (api.ApiResponse, error) {
	return client.baseQueryCall(Endpoints.OrderBookTicker, symbol, 0)
}

func (client *Client) AggregateTrades(
	symbol string,
	limit uint32,
	fromId uint64,
	startTime uint64,
	endTime uint64) (api.ApiResponse, error) {

	return client.baseGetCall(Endpoints.AggregateTrades, map[string]string{
		"symbol":    symbol,
		"limit":     strconv.FormatUint(uint64(limit), 10),
		"fromId":    strconv.FormatUint(uint64(fromId), 10),
		"startTime": strconv.FormatUint(startTime, 10),
		"endTime":   strconv.FormatUint(endTime, 10),
	})
}

func (client *Client) HistoricalTrades(
	symbol string,
	limit uint32,
	fromId uint64) (api.ApiResponse, error) {

	return client.baseGetCall(Endpoints.AggregateTrades, map[string]string{
		"symbol": symbol,
		"limit":  strconv.FormatUint(uint64(limit), 10),
		"fromId": strconv.FormatUint(uint64(fromId), 10),
	})
}

func (client *Client) Trades(symbol string, limit uint32) (api.ApiResponse, error) {
	return client.baseQueryCall(Endpoints.Trades, symbol, limit)
}

func (client *Client) Candlesticks(symbol string,
	limit uint32,
	interval string,
	startTime uint64,
	endTime uint64) (api.ApiResponse, error) {

	return client.baseGetCall(Endpoints.Candlesticks, map[string]string{
		"symbol":    symbol,
		"limit":     strconv.FormatUint(uint64(limit), 10),
		"interval":  interval,
		"startTime": strconv.FormatUint(startTime, 10),
		"endTime":   strconv.FormatUint(endTime, 10),
	})
}

func (client *Client) Ticker24h(symbol string) (api.ApiResponse, error) {
	return client.baseGetCall(Endpoints.Ticker24, map[string]string{
		"symbol": symbol,
	})
}

func (client *Client) TickerPrice(symbol string) (api.ApiResponse, error) {
	return client.baseGetCall(Endpoints.TickerPrice, map[string]string{
		"symbol": symbol,
	})
}

// TODO implement
func (client *Client) NewOrder(symbol string) (api.ApiResponse, error) {
	return client.baseGetCall(Endpoints.TickerPrice, map[string]string{
		"symbol": symbol,
	})
}

package main

import (
	"encoding/json"

	"github.com/diantanjung/order-book-matching/orderbook"
	"github.com/labstack/echo/v4"
)

func main()  {
	e := echo.New()
	ex := NewExchange()

	e.POST("/order", ex.handlePlaceOrder)

	e.Start(":3000")
}

type OrderType string

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder OrderType = "LIMIT"
)

type Market string

const (
	MarketETH Market ="ETH"
)

func NewExchange() *Exchange {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderBook()

	return &Exchange{
		orderbooks: orderbooks,
	}
}

type Exchange struct {
	orderbooks map[Market]*orderbook.Orderbook
}

type PlaceOrderRequest struct {
	Type OrderType // limit or market
	Bid bool
	Size float64
	Price float64
	Market Market
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)
	ob := ex.orderbooks[market]
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size)
	ob.PlaceLimitOrder(placeOrderData.Price, order)

	return c.JSON(200, map[string]any{"msg": "order placed"})
}


package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"home-broker/assets"
	"home-broker/assetwallets"
	"home-broker/core"
	"home-broker/money"
	"home-broker/users"
	"home-broker/wallets"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"
)

// OrderUseCases represents the order use cases.
type OrderUseCases struct {
	db            OrderDBInterface
	walletUC      wallets.WalletUseCases
	assetWalletUC assetwallets.AssetWalletUseCases
	orderBookHost string
}

// NewOrderUseCases returns a new OrderUseCases.
func NewOrderUseCases(db OrderDBInterface, walletUC wallets.WalletUseCases, assetWalletUC assetwallets.AssetWalletUseCases, orderBookHost string) OrderUseCases {
	return OrderUseCases{db: db, walletUC: walletUC, assetWalletUC: assetWalletUC, orderBookHost: orderBookHost}
}

// ExchangeOrderResponse represents a send order response of a exchange.
type exchangeOrderResponse struct {
	id        ExternalOrderID
	timestamp time.Time
	status    OrderStatus
}

// GetOrder returns an order by ID.
func (uc OrderUseCases) GetOrder(orderID OrderID) (*Order, error) {
	if orderID <= 0 {
		return nil, core.NewErrValidation("Invalid order ID.")
	}
	entity, err := uc.db.GetByID(orderID)
	return entity, err
}

// BuyOrder adds a buying order.
func (uc OrderUseCases) BuyOrder(userID users.UserID, assetID assets.AssetID, price money.Money, amount assets.AssetUnit) (*Order, error) {
	if userID <= 0 {
		return nil, core.NewErrValidation("Invalid user ID.")
	}
	if assetID == "" {
		return nil, core.NewErrValidation("Invalid asset ID.")
	}
	if price <= 0 {
		return nil, core.NewErrValidation("Invalid price.")
	}
	if amount <= 0 {
		return nil, core.NewErrValidation("Invalid amount.")
	}

	wallet, _, _, err := uc.walletUC.GetWallet(userID)
	if err != nil {
		return nil, err
	}
	if wallet.Balance < money.Money(int64(price)*int64(amount)) {
		return nil, core.NewErrValidation("No funds.")
	}

	entity := NewBuyOrder(assetID, amount, price)
	entity.UserID = userID
	entity.Status = OrderStatusPending

	newEntity, err := uc.db.Insert(entity)
	if err != nil {
		switch err {
		case assets.ErrAssetDoesNotExist:
			return nil, core.NewErrValidation("Asset does not exist.")
		case users.ErrUserDoesNotExist:
			return nil, core.NewErrValidation("User does not exist.")
		default:
			return nil, err
		}
	}

	// TODO: We can do this using a message broker to be asynchronous.
	// - Send to a topic/queue and leaves the order with "pending" status.
	// - Retuns this method with the order as "pending".
	// The Order entity already has a "status" for this:
	// - "pending" for sent or being send to the exchange
	// - "accepted" for orders accepted by the exchange
	// - "denied" for order denied by the exchange
	// For now we are just mocking the exchange response, so all requests are "accepted".
	response := uc.sendOrderToExchange(newEntity) // fake call

	// After that we update the order with the ID generate by the exchange (external ID).
	err = uc.db.UpdateExternalResponse(newEntity.ID, response.id, response.timestamp, response.status)
	if err != nil {
		return newEntity, err
	}

	// Refresh the order because it has a new state.
	newEntity, err = uc.db.GetByID(newEntity.ID)

	return newEntity, err
}

// SellOrder adds a selling order.
func (uc OrderUseCases) SellOrder(userID users.UserID, assetID assets.AssetID, price money.Money, amount assets.AssetUnit) (*Order, error) {
	if userID <= 0 {
		return nil, core.NewErrValidation("Invalid user ID.")
	}
	if assetID == "" {
		return nil, core.NewErrValidation("Invalid asset ID.")
	}
	if price <= 0 {
		return nil, core.NewErrValidation("Invalid price.")
	}
	if amount <= 0 {
		return nil, core.NewErrValidation("Invalid amount.")
	}

	assetWallet, _, _, err := uc.assetWalletUC.GetAssetWallet(userID, assetID)
	if err != nil {
		return nil, err
	}
	if assetWallet.Balance < amount {
		return nil, core.NewErrValidation("No assets.")
	}

	entity := NewSellOrder(assetID, amount, price)
	entity.UserID = userID
	entity.Status = OrderStatusPending

	newEntity, err := uc.db.Insert(entity)
	if err != nil {
		switch err {
		case assets.ErrAssetDoesNotExist:
			return nil, core.NewErrValidation("Asset does not exist.")
		case users.ErrUserDoesNotExist:
			return nil, core.NewErrValidation("User does not exist.")
		default:
			return nil, err
		}
	}

	// TODO: We can do this using a message broker to be asynchronous.
	// - Send to a topic/queue and leaves the order with "pending" status.
	// - Retuns this method with the order as "pending".
	// The Order entity already has a "status" for this:
	// - "pending" for sent or being send to the exchange
	// - "accepted" for orders accepted by the exchange
	// - "denied" for order denied by the exchange
	// For now we are just mocking the exchange response, so all requests are "accepted".
	response := uc.sendOrderToExchange(newEntity) // fake call

	// After that we update the order with the ID generate by the exchange (external ID).
	err = uc.db.UpdateExternalResponse(newEntity.ID, response.id, response.timestamp, response.status)
	if err != nil {
		return newEntity, err
	}

	// Refresh the order because it has a new state.
	newEntity, err = uc.db.GetByID(newEntity.ID)

	return newEntity, err
}

// CancelOrder returns an order by ID.
func (uc OrderUseCases) CancelOrder(orderID OrderID) (*Order, error) {
	if orderID <= 0 {
		return nil, core.NewErrValidation("Invalid order ID.")
	}
	entity, err := uc.db.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}
	// TODO: This not the best approach.
	// We should have a status group for "orders accepted" and "orders in canceling" processes.
	// A order should be able to go back to the previous status if the cancel fail for some reasons.
	// We lost the previous status if we ovewrite the value with "canceling".
	if entity.Status == OrderStatusAccepted {
		response := uc.cancelOrderOnExchange(entity.ID)
		// After that we update the order status.
		err = uc.db.UpdateStatus(entity.ID, response.status)
		if err != nil {
			return entity, err
		}

		// Refresh the order because it has a new state.
		entity, err = uc.db.GetByID(entity.ID)
	}
	return entity, err
}

// ProcessExternalUpdate process an order update received from an exchange service.
// These updates change the state of an Order Book.
func (uc OrderUseCases) ProcessExternalUpdate(externalUp ExternalUpdate) error {
	err := uc.updateOrderBook(externalUp)
	if err != nil {
		log.Printf("error to update order book: %v\n", err)
	}
	if externalUp.Action == ExternalUpdateActionTraded {
		err = uc.processExternalUpdateTraded(externalUp)
	}
	return err
}

func (uc OrderUseCases) updateOrderBook(externalUp ExternalUpdate) error {
	url := fmt.Sprintf("%s/api/v1/orderbooks/%s/webhook", uc.orderBookHost, externalUp.AssetID)
	log.Printf("sending to order book on %s...\n", url)
	json, err := json.Marshal(externalUp)
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(json))
	if err != nil {

		return err
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	log.Printf("response: %s\n", bodyString)
	return nil
}

func (uc OrderUseCases) processExternalUpdateAdded(externalUp ExternalUpdate) error {
	return nil
}

func (uc OrderUseCases) processExternalUpdateDeleted(externalUp ExternalUpdate) error {
	return nil
}

func (uc OrderUseCases) processExternalUpdateTraded(externalUp ExternalUpdate) error {
	entity, err := uc.db.GetByExternalIDAssetID(externalUp.ID, externalUp.AssetID)
	if err != nil {
		return err
	}
	if entity == nil {
		// The user and/or asset is not from this home broker.
		return nil
	}

	// TODO: This is not the best way.
	// We should put everything inside a transaction.
	switch entity.Type {
	case OrderTypeBuy:
		// On buy, we remove money and add assets.
		valueMoney := money.Money(math.Abs(float64(int64(externalUp.Amount) * int64(externalUp.Price))))
		valueMoney *= -1
		_, err = uc.walletUC.IncBalanceByUserID(entity.UserID, valueMoney)
		if err != nil {
			return err
		}
		amountAssets := assets.AssetUnit(math.Abs(float64(externalUp.Amount)))
		_, err = uc.assetWalletUC.IncBalanceByUserIDAssetID(entity.UserID, entity.AssetID, amountAssets)
		if err != nil {
			return err
		}
	case OrderTypeSell:
		// On sell, we add money and remove assets.
		valueMoney := money.Money(math.Abs(float64(int64(entity.Amount) * int64(entity.Price))))
		_, err = uc.walletUC.IncBalanceByUserID(entity.UserID, valueMoney)
		if err != nil {
			return err
		}
		amountAssets := assets.AssetUnit(math.Abs(float64(entity.Amount)))
		amountAssets *= -1
		_, err = uc.assetWalletUC.IncBalanceByUserIDAssetID(entity.UserID, entity.AssetID, amountAssets)
		if err != nil {
			return err
		}
	}
	return nil
}

// sendOrderToExchange sends an order to the exchange.
//   For now we are just mocking the exchange response, so all requests are "accepted".
func (uc OrderUseCases) sendOrderToExchange(order *Order) exchangeOrderResponse {
	/*
		This part depends of the exchange API (ex B3/Nasdaq).
		- Connect to the exchange.
		- Send the order.
		- Receive the generated exchange order ID (external ID).

		We can use the ExchangeID inside the asset (order.asset.ExchangeID) to choose the correct exchange.
	*/

	// We are just mocking the response here.
	return exchangeOrderResponse{
		id:        ExternalOrderID(fmt.Sprintf("EX-%v", order.ID)), // fake value
		timestamp: time.Now(),                                      // fake value
		status:    OrderStatusAccepted,
	}
}

// cancelOrderOnExchange cancel an order on the exchange.
//   For now we are just mocking the exchange response, so all requests are "accepted".
func (uc OrderUseCases) cancelOrderOnExchange(orderID OrderID) exchangeOrderResponse {
	/*
		This part depends of the exchange API (ex B3/Nasdaq).
		- Connect to the exchange.
		- Cancel the order.
		- Receive a accepted/denied status.
	*/

	// We are just mocking the response here.
	return exchangeOrderResponse{
		id:        ExternalOrderID(fmt.Sprintf("EX-%v", orderID)), // fake value
		timestamp: time.Now(),                                     // fake value
		status:    OrderStatusCanceled,
	}
}

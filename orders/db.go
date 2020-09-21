package orders

import "time"

// OrderDBInterface is an interface that handles database commands for Order entity.
type OrderDBInterface interface {

	// GetByID must return an order by ID.
	// If the record does not exist a nil entity will be returned.
	GetByID(id OrderID) (*Order, error)

	// Insert must insert a new order.
	// A nil entity will be returned if an error occurs.
	// The following errors can happen: ErrUserDoesNotExist, ErrAssetDoesNotExist.
	Insert(entity Order) (*Order, error)

	// UpdateExternalResponse updates an order base on a exchange response.
	UpdateExternalResponse(orderID OrderID, externalID ExternalOrderID, externalTimestamp time.Time, status OrderStatus) error

	// UpdateStatus updates an order status.
	UpdateStatus(orderID OrderID, status OrderStatus) error
}

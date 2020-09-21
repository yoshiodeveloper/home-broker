package assets

// AssetDBInterface is an interface that handles database commands for Asset entity.
type AssetDBInterface interface {
	// GetByID must return an asset by ID.
	// A nil entity will be returned if it does not exist.
	GetByID(id AssetID) (*Asset, error)

	// Insert must insert a new asset.
	// A nil entity will be returned if an error occurs.
	// The following errors can happen: ErrUserDoesNotExist.
	Insert(entity Asset) (*Asset, error)
}

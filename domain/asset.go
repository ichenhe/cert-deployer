package domain

import (
	"errors"
	"fmt"
)

var _ Asseter = &Asset{}

type AssetType string

const (
	// Universal Type -------------------------

	TypeCdn AssetType = "cdn"
)

func (at AssetType) String() string {
	return string(at)
}

func AssetTypeFromString(t string) (AssetType, error) {
	switch t {
	case TypeCdn.String():
		return TypeCdn, nil
	default:
		return "", errors.New("invalid asset type")
	}
}

// Asset represents an item of a product or a service in your account, such as
// a host, a bucket in object storage or a CDN domain.
type Asset struct {
	Name      string
	Id        string
	Type      AssetType // One of constants 'domain.Type*', e.g. TypeCdn.
	Provider  string    // e.g. tencent.Provider
	Available bool      // This asset is ready to set up.
}

func (a *Asset) String() string {
	return fmt.Sprintf("{Asset provider=%s, type=%s, name=%s, id=%s, available=%v}",
		a.Provider, a.Type, a.Name, a.Id, a.Available)
}

func (a *Asset) GetBaseInfo() *Asset {
	return a
}

type Asseter interface {
	GetBaseInfo() *Asset
}

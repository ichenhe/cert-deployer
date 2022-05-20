package search

import (
	"github.com/ichenhe/cert-deployer/asset"
	"github.com/ichenhe/cert-deployer/config"
	"github.com/ichenhe/cert-deployer/utils"
	"go.uber.org/zap"
)

// UnionSearch is a special AssetSearcher. It will not be registered to the list. Instead, it is
// responsible for instantiating registered searchers based on configuration and call each of them.
//
// The main program usually wants to use this structure to perform the operation.
type UnionSearch struct {
	logger         *zap.SugaredLogger
	assetSearchers map[string][]AssetSearcher
}

func NewUnionSearch(logger *zap.SugaredLogger, providersConfig []config.CloudProvider) *UnionSearch {
	uniSearch := &UnionSearch{
		logger:         logger,
		assetSearchers: make(map[string][]AssetSearcher),
	}

	if providersConfig != nil && len(providersConfig) > 0 {
		for i, conf := range providersConfig {
			if constructor, ok := assetSearcherConstructors[conf.Provider]; ok {
				options := map[string]interface{}{
					"secretId":  conf.SecretId,
					"secretKey": conf.SecretKey,
				}
				if searcher, err := constructor(options); err != nil {
					logger.Errorf("failed to create asset searcher for provider '%s': %v, index=%d",
						conf.Provider, err, i)
				} else {
					uniSearch.addSearcher(conf.Provider, searcher)
					logger.Debugf("add asset searcher $%d for provider '%s'", i, conf.Provider)
				}
			} else {
				logger.Errorf("can not find asset searcher constructor for provider '%s', index=%d",
					conf.Provider, i)
			}
		}
	}

	return uniSearch
}

// List calls each of registered searcher and return all the results.
//
// Note: Returning error does not mean that other return values are invalid.
// It may just be some searcher execution failures.
func (u *UnionSearch) List(assetType string) ([]asset.Asseter, *utils.ErrorCollection) {
	r := make([]asset.Asseter, 0, 64)
	errs := make([]error, 0)
	for _, v := range u.assetSearchers {
		for _, searcher := range v {
			if l, e := searcher.List(assetType); e != nil {
				errs = append(errs, e)
			} else {
				r = append(r, l...)
			}
		}
	}
	return r, utils.NewErrorCollection(errs)
}

// ListApplicable calls each of registered searcher and return all the results.
//
// Note: Returning error does not mean that other return values are invalid.
// It may just be some searcher execution failures.
func (u *UnionSearch) ListApplicable(assetType string, cert []byte) ([]asset.Asseter,
	*utils.ErrorCollection) {
	r := make([]asset.Asseter, 0, 64)
	errs := make([]error, 0)
	for _, v := range u.assetSearchers {
		for _, searcher := range v {
			if l, e := searcher.ListApplicable(assetType, cert); e != nil {
				errs = append(errs, e)
			} else {
				r = append(r, l...)
			}
		}
	}
	return r, utils.NewErrorCollection(errs)
}

func (u *UnionSearch) addSearcher(provider string, searcher AssetSearcher) {
	if searchers, ok := u.assetSearchers[provider]; ok {
		u.assetSearchers[provider] = append(searchers, searcher)
	} else {
		u.assetSearchers[provider] = []AssetSearcher{searcher}
	}
}

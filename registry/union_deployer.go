package registry

import (
	"github.com/ichenhe/cert-deployer/domain"
	"go.uber.org/zap"
)

// UnionDeployer will not be registered to the list. Instead, it is
// responsible for instantiating registered deployers based on configuration and call each of them.
//
// The main program usually wants to use this structure to perform the operation.
type UnionDeployer struct {
	logger    *zap.SugaredLogger
	deployers map[string][]domain.Deployer
}

func NewUnionDeployer(logger *zap.SugaredLogger, providersConfig map[string]domain.CloudProvider) *UnionDeployer {
	uniDeployer := &UnionDeployer{
		logger:    logger,
		deployers: make(map[string][]domain.Deployer),
	}

	if providersConfig != nil && len(providersConfig) > 0 {
		for name, conf := range providersConfig {
			if constructor, ok := assetDeployerConstructors[conf.Provider]; ok {
				options := map[string]interface{}{
					"secretId":  conf.SecretId,
					"secretKey": conf.SecretKey,
					"logger":    logger,
				}
				if searcher, err := constructor(options); err != nil {
					logger.Errorf("failed to create asset deployer for provider '%s': %v",
						name, err)
				} else {
					uniDeployer.addDeployer(conf.Provider, searcher)
					logger.Debugf("add asset deployer for provider '%s'", conf.Provider)
				}
			} else {
				logger.Errorf("can not find asset deployer constructor for provider '%s'",
					name)
			}
		}
	}

	return uniDeployer
}

// ListAssets calls each of registered deployer and return all the results.
//
// Note: Returning error does not mean that other return values are invalid.
// It may just be some searcher execution failures.
func (u *UnionDeployer) ListAssets(assetType domain.AssetType) ([]domain.Asseter, *domain.ErrorCollection) {
	r := make([]domain.Asseter, 0, 64)
	errs := make([]error, 0)
	for _, v := range u.deployers {
		for _, searcher := range v {
			if l, e := searcher.ListAssets(assetType); e != nil {
				errs = append(errs, e)
			} else {
				r = append(r, l...)
			}
		}
	}
	return r, domain.NewErrorCollection(errs)
}

// ListApplicableAssets calls each of registered deployer and return all the results.
//
// Note: Returning error does not mean that other return values are invalid.
// It may just be some searcher execution failures.
func (u *UnionDeployer) ListApplicableAssets(assetType domain.AssetType, cert []byte) ([]domain.Asseter,
	*domain.ErrorCollection) {
	r := make([]domain.Asseter, 0, 64)
	errs := make([]error, 0)
	for _, v := range u.deployers {
		for _, searcher := range v {
			if l, e := searcher.ListApplicableAssets(assetType, cert); e != nil {
				errs = append(errs, e)
			} else {
				r = append(r, l...)
			}
		}
	}
	return r, domain.NewErrorCollection(errs)
}

// Deploy calls each of registered deployer. All errors will be printed to logger with ERROR level.
func (u *UnionDeployer) Deploy(assets []domain.Asseter, cert []byte,
	key []byte) (deployedAsseters []domain.Asseter, hasError bool) {
	for _, assetItem := range assets {
		for _, deployer := range u.deployers[assetItem.GetBaseInfo().Provider] {
			deployed, deployErrs := deployer.Deploy([]domain.Asseter{assetItem}, cert, key)
			if deployed != nil {
				deployedAsseters = append(deployedAsseters, deployed...)
			}
			if deployErrs != nil && len(deployErrs) > 0 {
				hasError = true
				for _, e := range deployErrs {
					u.logger.Error(e.Error())
				}
			}
		}
	}
	return
}

func (u *UnionDeployer) addDeployer(provider string, deployer domain.Deployer) {
	if searchers, ok := u.deployers[provider]; ok {
		u.deployers[provider] = append(searchers, deployer)
	} else {
		u.deployers[provider] = []domain.Deployer{deployer}
	}
}

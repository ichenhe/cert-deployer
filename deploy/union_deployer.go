package deploy

import (
	"github.com/ichenhe/cert-deployer/asset"
	"github.com/ichenhe/cert-deployer/config"
	"go.uber.org/zap"
)

// UnionDeployer will not be registered to the list. Instead, it is
// responsible for instantiating registered deployers based on configuration and call each of them.
//
// The main program usually wants to use this structure to perform the operation.
type UnionDeployer struct {
	logger    *zap.SugaredLogger
	deployers map[string][]Deployer
}

func NewUnionDeployer(logger *zap.SugaredLogger, providersConfig []config.CloudProvider) *UnionDeployer {
	uniDeployer := &UnionDeployer{
		logger:    logger,
		deployers: make(map[string][]Deployer),
	}

	if providersConfig != nil && len(providersConfig) > 0 {
		for i, conf := range providersConfig {
			if constructor, ok := assetDeployerConstructors[conf.Provider]; ok {
				options := map[string]interface{}{
					"secretId":  conf.SecretId,
					"secretKey": conf.SecretKey,
					"logger":    logger,
				}
				if searcher, err := constructor(options); err != nil {
					logger.Errorf("failed to create asset deployer for provider '%s': %v, index=%d",
						conf.Provider, err, i)
				} else {
					uniDeployer.addDeployer(conf.Provider, searcher)
					logger.Debugf("add asset deployer $%d for provider '%s'", i, conf.Provider)
				}
			} else {
				logger.Errorf("can not find asset deployer constructor for provider '%s', index=%d",
					conf.Provider, i)
			}
		}
	}

	return uniDeployer
}

// Deploy calls each of registered deployer. All errors will be printed to logger with ERROR level.
func (u *UnionDeployer) Deploy(assets []asset.Asseter, cert []byte,
	key []byte) (deployedAsseters []asset.Asseter, hasError bool) {
	for _, assetItem := range assets {
		for _, deployer := range u.deployers[assetItem.GetBaseInfo().Provider] {
			deployed, deployErrs := deployer.Deploy([]asset.Asseter{assetItem}, cert, key)
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

func (u *UnionDeployer) addDeployer(provider string, deployer Deployer) {
	if searchers, ok := u.deployers[provider]; ok {
		u.deployers[provider] = append(searchers, deployer)
	} else {
		u.deployers[provider] = []Deployer{deployer}
	}
}

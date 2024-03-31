package domain

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"reflect"
)

type DeployResultCallbackFunc = func(asset Asseter, err error)

type DeployPreExecuteCallbackFunc = func(asset Asseter)

type DeployCallback struct {
	ResultCallback     DeployResultCallbackFunc
	PreExecuteCallback DeployPreExecuteCallbackFunc
}

type Deployer interface {
	// IsAssetTypeSupported checks whether the given asset type is supported by the deployer.
	//
	// This function must always return the consistent result for the same input.
	IsAssetTypeSupported(assetType string) bool

	// ListAssets fetches all assets that match the given type.
	// The assetType should be one of constants 'asset.Type*', e.g. asset.TypeCdn.
	ListAssets(ctx context.Context, assetType string) ([]Asseter, error)

	// ListApplicableAssets fetch all assets that match the given type and cert.
	// The assetType should be one of constants 'asset.Type*', e.g. asset.TypeCdn.
	ListApplicableAssets(ctx context.Context, assetType string, cert []byte) ([]Asseter, error)

	// Deploy the given pem cert to the all assets.
	//
	// The result of each deployment is returned by the callback which can be nil if not interested.
	//
	// The return value itself indicates the general error rather than deployment result.
	Deploy(ctx context.Context, assets []Asseter, cert []byte, key []byte, callback *DeployCallback) error
}

type DeployerFactory interface {
	// NewDeployer creates a deployer corresponding to the given cloudProvider.
	NewDeployer(logger *zap.SugaredLogger, cloudProvider CloudProvider) (Deployer, error)
}

// DeployerHelper provides some useful functions for deployer implementation.
type DeployerHelper struct {
}

// OnDeployResult calls the callback function if it is not nil.
func (d *DeployerHelper) OnDeployResult(cb *DeployCallback, asset Asseter, err error) {
	if cb != nil && cb.ResultCallback != nil {
		cb.ResultCallback(asset, err)
	}
}

// OnPreDeploy calls the callback function if it is not nil.
func (d *DeployerHelper) OnPreDeploy(cb *DeployCallback, asset Asseter) {
	if cb != nil && cb.PreExecuteCallback != nil {
		cb.PreExecuteCallback(asset)
	}
}

type Options map[string]interface{}

const OptionsKeyLogger = "cert-deployer.logger"

// MustReadString reads a string value from options. See MustReadOption.
func (o Options) MustReadString(key string) string {
	return MustReadOption[string](o, key)
}

// MustReadLogger gets the logger from options. See MustReadOption.
func (o Options) MustReadLogger() *zap.SugaredLogger {
	return MustReadOption[*zap.SugaredLogger](o, OptionsKeyLogger)
}

// MustReadOption reads a value from options. Panic with InvalidOptionError if key does not exist or
// the type of value is not T.
func MustReadOption[T any | *zap.SugaredLogger](opt Options, key string) T {
	if v, ok := opt[key]; !ok {
		panic(NewInvalidOptionError(fmt.Sprintf("option '%s' does not exist", key)))
	} else if s, ok := v.(T); !ok {
		var tmp T
		panic(NewInvalidOptionError(fmt.Sprintf("type of option '%s' should be %T, actual is %s", key,
			reflect.TypeOf(tmp), reflect.TypeOf(v).String())))
	} else {
		return s
	}
}

var _ error = &InvalidOptionError{}

type InvalidOptionError struct {
	Message string
}

func NewInvalidOptionError(message string) *InvalidOptionError {
	return &InvalidOptionError{Message: message}
}

func (e *InvalidOptionError) Error() string {
	return e.Message
}

// RecoverFromInvalidOptionError catches the InvalidOptionError and handles it to the given handler.
//
// Usage:
//
//	defer RecoverFromInvalidOptionError(func(e *domain.InvalidOptionError) {
//		// deal with the error...
//	})
func RecoverFromInvalidOptionError(handler func(err *InvalidOptionError)) {
	if v := recover(); v != nil {
		if e, ok := v.(*InvalidOptionError); ok {
			handler(e)
		} else {
			panic(v)
		}
	}
}

// DeploymentExecutor is responsible for executing deployment defined in the profile.
// Typically, implements may contain a field to store domain.CloudProvider s.
type DeploymentExecutor interface {
	// ExecuteDeployment executes a deployment.
	//
	// If a deployer is created, it is considered a successful execution, even if no assets were
	// deployed successfully. Because one deployment may contain many assets, it's confused to say
	// whether it is success.
	ExecuteDeployment(ctx context.Context, deployment Deployment) error
}

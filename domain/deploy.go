package domain

import (
	"fmt"
	"go.uber.org/zap"
	"reflect"
)

type Deployer interface {
	// ListAssets fetches all assets that match the given type.
	// The assetType should be one of constants 'asset.Type*', e.g. asset.TypeCdn.
	ListAssets(assetType string) ([]Asseter, error)

	// ListApplicableAssets fetch all assets that match the given type and cert.
	// The assetType should be one of constants 'asset.Type*', e.g. asset.TypeCdn.
	ListApplicableAssets(assetType string, cert []byte) ([]Asseter, error)

	// Deploy the given pem cert to the all assets.
	//
	// Returns assets that were successfully deployed and errors. Please note that there is no
	// guarantee that len(deployedAsseters)+len(deployErrs)=len(assets), because some minor
	// problems do not count as errors, such as provider mismatch.
	Deploy(assets []Asseter, cert []byte, key []byte) (deployedAssets []Asseter,
		deployErrs []*DeployError)
}

type DeployerFactory interface {
	// NewDeployer creates a deployer corresponding to the given cloudProvider.
	NewDeployer(logger *zap.SugaredLogger, cloudProvider CloudProvider) (Deployer, error)
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

var _ error = &DeployError{}

type DeployError struct {
	Asset Asseter
	Err   error
}

func (d *DeployError) Error() string {
	return fmt.Sprintf("failed to deploy %v: %v", d.Asset, d.Err)
}

func NewDeployError(asset Asseter, err error) *DeployError {
	return &DeployError{
		Asset: asset,
		Err:   err,
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
	ExecuteDeployment(deployment Deployment) error
}

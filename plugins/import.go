// Package plugins is used to import all plug-ins centrally.
package plugins

import (
	_ "github.com/ichenhe/cert-deployer/plugins/alibaba"
	_ "github.com/ichenhe/cert-deployer/plugins/aws"
	_ "github.com/ichenhe/cert-deployer/plugins/tencent"
)

package config

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_mustUnmarshalSpecificTrigger(t *testing.T) {
	type args struct {
		triggerDef  map[string]any
		triggerName string
	}
	tests := []struct {
		name  string
		args  args
		panic bool
	}{
		{name: "ok", panic: false, args: args{
			triggerDef: map[string]any{
				"type":         "file_monitoring",
				"deployments":  []string{"a", "b", "c"},
				"options.file": "/path/to/file",
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := koanf.New(".")
			for path, v := range tt.args.triggerDef {
				err := k.Set(fmt.Sprintf("triggers.%s.%s", tt.args.triggerName, path), v)
				require.Nil(t, err, "failed to load preset trigger defs")
			}

			trigger := domain.FileMonitoringTriggerDef{}
			if tt.panic {
				assert.Panics(t, func() {
					mustUnmarshalSpecificTrigger(k, tt.args.triggerName, &trigger)
				})
				return
			}

			// no panic
			assert.NotPanics(t, func() {
				mustUnmarshalSpecificTrigger(k, tt.args.triggerName, &trigger)
			})
			assert.Equal(t, tt.args.triggerName, trigger.GetName(), "not expected trigger's name")
			assert.Equal(t, tt.args.triggerDef["type"], trigger.Type, "not expected trigger's type")
			assert.Equal(t, tt.args.triggerDef["deployments"], trigger.GetDeploymentIds(), "not expected deployments")
		})
	}
}

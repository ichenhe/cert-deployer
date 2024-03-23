package config

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/knadh/koanf/v2"
	"reflect"
)

// parseTriggers parses raw trigger configuration and load into AppConfig.
// If an error occurred, AppConfig won't change.
//
// The default value of triggers' options are set in this function.
func parseTriggers(k *koanf.Koanf, config *domain.AppConfig) (err error) {
	//triggerDefs := make(map[string]domain.TriggerDef)
	//if err = k.Unmarshal("triggers", &triggerDefs); err != nil {
	//	return fmt.Errorf("failed to unmarshal triggers: %w", err)
	//}
	triggerNames := k.MapKeys("triggers")
	parsedTriggers := make(map[string]domain.TriggerDefiner, len(triggerNames))

	defer func() {
		// catch any panics, typically raised by mustUnmarshalSpecificTrigger().
		if v := recover(); v != nil {
			if e, ok := v.(error); ok {
				err = fmt.Errorf("failed to parse trigger: %w", e)
			} else {
				err = fmt.Errorf("failed to parse trigger: %v", e)
			}
		}
	}()

	for _, name := range triggerNames {
		triggerType := k.String("triggers." + name + ".type")
		var t domain.TriggerDefiner
		switch triggerType {
		case "file_monitoring":
			t = mustUnmarshalSpecificTrigger(k, name, &domain.FileMonitoringTriggerDef{
				Options: domain.FileMonitoringTriggerOptions{
					Event: "content_change",
					Delay: 1000,
				},
			})
		default:
			return fmt.Errorf("invalid triggers[%s]'s type: %s", name, triggerType)
		}
		parsedTriggers[name] = t
	}

	config.Triggers = parsedTriggers
	return nil
}

// mustUnmarshalSpecificTrigger parses the Koanf to a specific type of trigger struct which is
// provided by the dst. The dst must be a pointer to the corresponding struct. Panic if any error.
func mustUnmarshalSpecificTrigger(k *koanf.Koanf, triggerName string, dst domain.TriggerDefiner) domain.TriggerDefiner {
	v := reflect.ValueOf(dst)
	ele := v.Elem()
	if v.Type().Kind() != reflect.Pointer || ele.Kind() != reflect.Struct {
		panic("dst must be an ptr to struct")
	}

	if err := k.Unmarshal("triggers."+triggerName, v.Interface()); err != nil {
		panic(err)
	}
	// set trigger's name
	nameField := ele.FieldByName("Name")
	nameField.SetString(triggerName)
	return dst
}

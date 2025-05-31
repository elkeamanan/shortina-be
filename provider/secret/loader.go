package secret

import (
	"context"
	"fmt"
	"reflect"
)

func LoadValueFromSecret(ctx context.Context, manager SecretsManager, cfg interface{}) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be a pointer to struct")
	}

	return populateSecrets(ctx, manager, v.Elem())
}

func populateSecrets(ctx context.Context, manager SecretsManager, v reflect.Value) error {
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		val := v.Field(i)

		if !val.CanSet() {
			continue
		}

		if val.Kind() == reflect.Struct {
			if err := populateSecrets(ctx, manager, val); err != nil {
				return err
			}
			continue
		}

		if secretName := field.Tag.Get("secret"); secretName != "" {
			secretVal, err := manager.GetSecret(ctx, secretName)
			if err != nil {
				return fmt.Errorf("failed to get secret for field %s: %w", field.Name, err)
			}

			if val.Kind() == reflect.String {
				val.SetString(secretVal)
			} else {
				return fmt.Errorf("field %s must be a string to hold secret", field.Name)
			}
		}
	}

	return nil
}

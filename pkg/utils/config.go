package utils

import (
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Loader[T any] interface {
	Load() (T, error)
}

type ConfigLoader[T any] struct {
	viper      *viper.Viper
	hasSecrets bool
}

const SECRETS_PATH = "/run/secrets"

func NewConfigLoader[T any](configPath string, configType string, hasSecrets bool) (*ConfigLoader[T], error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType(configType)

	loader := &ConfigLoader[T]{}
	loader.viper = v
	loader.hasSecrets = hasSecrets

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	content := os.ExpandEnv(string(file))
	if err := v.ReadConfig(strings.NewReader(content)); err != nil {
		return nil, err
	}

	return loader, nil
}

func (l *ConfigLoader[T]) Load() (T, error) {
	var config T

	if err := l.viper.Unmarshal(&config); err != nil {
		return config, err
	}
	if l.hasSecrets {
		if err := l.ReadSecrets(&config); err != nil {
			return config, err
		}
	}

	return config, nil
}

func (l *ConfigLoader[T]) ReadSecrets(config any) error {
	val := reflect.ValueOf(config)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if !field.CanSet() {
			continue
		}

		if field.Kind() == reflect.Struct {
			if err := l.ReadSecrets(field.Addr().Interface()); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.String {
			fieldValue := field.String()

			if strings.Contains(fieldValue, SECRETS_PATH) {
				fileContent, err := os.ReadFile(fieldValue)
				if err != nil {
					return err
				}
				field.SetString(strings.TrimSpace(string(fileContent)))
			}
		}
	}

	return nil
}

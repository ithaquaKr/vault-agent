package main

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configVar [TODO:description]
func configVar(cmd *cobra.Command, key string, defaultValue interface{}, description string, vp *viper.Viper) error {
	flags := cmd.PersistentFlags()
	flagSet := reflect.ValueOf(defaultValue)
	flagType := flagSet.Type()

	switch flagType.Kind() {
	case reflect.Bool:
		flags.Bool(key, defaultValue.(bool), description)
	case reflect.Int:
		flags.Int(key, defaultValue.(int), description)
	case reflect.String:
		flags.String(key, defaultValue.(string), description)
	case reflect.Slice:
		if flagType.Elem().Kind() == reflect.String {
			flags.StringSlice(key, defaultValue.([]string), description)
		} else {
			return fmt.Errorf("unsupported slice type for key: %s", key)
		}
	// TODO: Support another type if needed.
	default:
		return fmt.Errorf("unsupported flag type for key: %s", key)
	}

	err := vp.BindPFlag(key, flags.Lookup(key))
	if err != nil {
		return fmt.Errorf("failed to bind flag '%s': %w", key, err)
	}

	return nil
}

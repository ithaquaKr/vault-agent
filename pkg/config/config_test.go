package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestBindFlag(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue interface{}
		description  string
		expectedErr  bool
	}{
		{
			name:         "bool flag",
			key:          "test-bool",
			defaultValue: true,
			description:  "A test boolean flag",
			expectedErr:  false,
		},
		{
			name:         "int flag",
			key:          "test-int",
			defaultValue: 123,
			description:  "A test integer flag",
			expectedErr:  false,
		},
		{
			name:         "string flag",
			key:          "test-string",
			defaultValue: "hello",
			description:  "A test string flag",
			expectedErr:  false,
		},
		{
			name:         "string slice flag",
			key:          "test-string-slice",
			defaultValue: []string{"a", "b"},
			description:  "A test string slice flag",
			expectedErr:  false,
		},
		{
			name:         "unsupported type",
			key:          "test-unsupported",
			defaultValue: 1.23,
			description:  "An unsupported type flag",
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "test"}
			vp := viper.New()

			err := BindFlag(cmd, tt.key, tt.defaultValue, tt.description, vp)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify that the flag was added
				flag := cmd.PersistentFlags().Lookup(tt.key)
				assert.NotNil(t, flag)
				assert.Equal(t, tt.description, flag.Usage)

				// Verify that viper can get the value
				assert.Equal(t, tt.defaultValue, vp.Get(tt.key))
			}
		})
	}
}


package gateway

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config ...
type Config struct {
	Name     string          `yaml:"name"`
	Location string          `yaml:"location,omitempty"`
	Devices  []*DeviceConfig `yaml:"devices"`
}

// DeviceConfig ...
type DeviceConfig struct {
	Kind       string `yaml:"kind"`
	DeviceID   string `yaml:"deviceID"`
	PrivateKey string `yaml:"privateKey"`
	Algorithm  string `yaml:"algorithm"`
}

// unmarshalConfig unmarshals and validates the given payload.
func unmarshalConfig(in []byte, out *Config) error {
	err := yaml.Unmarshal(in, out)
	if err != nil {
		return err
	}
	if out.Name == "" {
		return errors.New("missing value for field 'name'")
	}
	if len(out.Devices) == 0 {
		return errors.New("missing value for field 'devices'")
	}
	for _, dev := range out.Devices {
		if dev.DeviceID == "" {
			return errors.New("missing value for field 'deviceID'")
		}
		if !(dev.Kind == "opcua" || dev.Kind == "modbus") {
			return errors.New("value for field 'kind' must be one of 'opcua', 'modbus'")
		}
		//TODO: finish rest of validation
	}
	//TODO: validate there are no duplicate deviceIDs

	return nil
}

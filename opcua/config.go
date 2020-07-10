package opcua

import (
	"errors"
	"time"

	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Name             string        `yaml:"name"`
	Location         string        `yaml:"location,omitempty"`
	EndpointURL      string        `yaml:"endpointURL"`
	SecurityPolicy   string        `yaml:"securityPolicy,omitempty"`
	Username         string        `yaml:"username,omitempty"`
	Password         string        `yaml:"password,omitempty"`
	PublishInterval  float64       `yaml:"publishInterval"`
	SamplingInterval float64       `yaml:"samplingInterval"`
	QueueSize        uint32        `yaml:"queueSize"`
	PublishTTL       time.Duration `yaml:"publishTTL"`
	Metrics          []*Metric     `yaml:"metrics"`
}

// Metric ...
type Metric struct {
	MetricID         string  `yaml:"metricID"`
	NodeID           string  `yaml:"nodeID"`
	AttributeID      uint32  `yaml:"attributeID,omitempty"`
	IndexRange       string  `yaml:"indexRange,omitempty"`
	SamplingInterval float64 `yaml:"samplingInterval,omitempty"`
	QueueSize        uint32  `yaml:"queueSize,omitempty"`
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
	if out.EndpointURL == "" {
		return errors.New("missing value for field 'endpointURL'")
	}
	if out.PublishInterval == 0 {
		return errors.New("missing value for field 'publishInterval'")
	}
	if out.SamplingInterval == 0 {
		return errors.New("missing value for field 'samplingInterval'")
	}
	if out.QueueSize == 0 {
		return errors.New("missing value for field 'queueSize'")
	}
	if out.PublishTTL == 0 {
		out.PublishTTL = 24 * time.Hour
	}
	if len(out.Metrics) == 0 {
		return errors.New("missing value for field 'metrics'")
	}

	for _, metric := range out.Metrics {
		if metric.MetricID == "" {
			return errors.New("missing value for field 'metricID'")
		}
		if metric.NodeID == "" {
			return errors.New("missing value for field 'nodeID'")
		}
		if metric.AttributeID == 0 {
			metric.AttributeID = 13
		}
		if metric.SamplingInterval == 0 {
			metric.SamplingInterval = out.SamplingInterval
		}
		if metric.QueueSize == 0 {
			metric.QueueSize = out.QueueSize
		}
	}
	//TODO: validate there are no duplicate metricIDs

	return nil
}

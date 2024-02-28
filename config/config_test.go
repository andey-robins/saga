package config

import "testing"

func TestParseExampleConfig(t *testing.T) {
	config := ParseConfig("../input/config/test.json")
	if config == nil {
		t.Error("expected config to be non-nil")
		return
	}

	if len(config.Populations) != 2 {
		t.Error("expected 2 populations, got", len(config.Populations))
	}

	if len(config.Jobs) != 2 {
		t.Error("expected 2 job, got", len(config.Jobs))
	}
}

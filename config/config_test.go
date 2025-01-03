package config_test

import (
	"maker-checker/config"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	want := "golang-api-template-test"

	viper.SetDefault("APP_NAME", want)

	cfg := config.Config()

	if got := cfg.App.Name; got != want {
		t.Fatalf("config: expected %s, got: %s", want, got)
	}
}

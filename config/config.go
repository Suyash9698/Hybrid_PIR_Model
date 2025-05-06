package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	N           int
	Mu          float64
	BasePort    int
	DataDir     string
	DBPath      string
	HTTPTimeout time.Duration
	ServerID    int // ⬅️ NEW
}

func Load() (*Config, error) {
	viper.SetDefault("SERVERID", 0)
	viper.SetEnvPrefix("PIR")
	viper.AutomaticEnv()
	viper.SetDefault("N", 6)
	viper.SetDefault("MU", 0.5)
	viper.SetDefault("BASEPORT", 8000)
	viper.SetDefault("DATADIR", "./data")
	viper.SetDefault("DBPATH", "./data/meta.db")
	viper.SetDefault("HTTPTIMEOUT", "3s")

	cfg := &Config{
		N:           viper.GetInt("N"),
		Mu:          viper.GetFloat64("MU"),
		BasePort:    viper.GetInt("BASEPORT"),
		DataDir:     viper.GetString("DATADIR"),
		DBPath:      viper.GetString("DBPATH"),
		HTTPTimeout: viper.GetDuration("HTTPTIMEOUT"),
		ServerID:    viper.GetInt("SERVERID"), // <-- NEW
	}
	return cfg, nil
}

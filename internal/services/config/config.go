package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	App  App  `json:"app"`
	Core Core `json:"core"`
}

type Duration struct {
	Time time.Duration
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Time = dur

	return nil
}

type App struct {
	Addr         string   `json:"addr"`
	ReadTimeout  Duration `json:"read_timeout"`
	WriteTimeout Duration `json:"write_timeout"`
	IdleTimeout  Duration `json:"idle_timeout"`
}

type Core struct {
	GracefulTimeout Duration `json:"graceful_timeout"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("./config/dev.json")
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

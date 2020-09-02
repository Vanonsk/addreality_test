package resources

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/kelseyhightower/envconfig"
)

type Res struct {
	Config Config
	DB     *DB
}

type Config struct {
	MetricsPort int `envconfig:"METRICS_PORT" default:"2112" required:"true"`
	RESTAPIPort int `envconfig:"PORT" default:"8080" required:"true"`
}

func New(logger *zap.SugaredLogger) (*Res, error) {
	conf := Config{}
	err := envconfig.Process("", &conf)
	if err != nil {
		return nil, fmt.Errorf("can't process the config: %w", err)
	}

	db := NewDB(logger)
	return &Res{
		Config: conf,
		DB:     db,
	}, nil
}

func (r *Res) Release() {
	r.DB.Close()
}

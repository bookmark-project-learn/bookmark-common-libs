package rabbitmq

import "github.com/kelseyhightower/envconfig"

// config represents the configuration of RabbitMQ, no exchange
type RabbitMQConfig struct {
	RabbitUser     string `envconfig:"RABBITMQ_DEFAULT_USER" default:"guest"`
	RabbitPassword string `envconfig:"RABBITMQ_DEFAULT_PASS" default:"guest"`
	Host           string `envconfig:"RABBITMQ_HOST" default:"localhost"`
	Port           string `envconfig:"RABBITMQ_PORT" default:"5672"`
}

// newConfig creates a new config struct based on the environment variables
func newConfig(envPrefix string) (*RabbitMQConfig, error) {
	cfg := &RabbitMQConfig{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *RabbitMQConfig) GetBrokerURI() string {
	return "amqp://" + c.RabbitUser + ":" + c.RabbitPassword + "@" + c.Host + ":5672/"
}

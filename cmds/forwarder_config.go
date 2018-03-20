package cmds

import (
	"fmt"
	"os"
	"strings"

	"github.com/BaritoLog/barito-flow/river"
	"github.com/sirupsen/logrus"
)

const (
	EnvForwarderKafkaBrokers         = "BARITO_FORWARDER_KAFKA_BROKERS"
	EnvForwarderKafkaConsumerGroupId = "BARITO_FORWARDER_KAFKA_CONSUMER_GROUP_ID"
	EnvForwarderKafkaConsumerTopic   = "BARITO_FORWARDER_KAFKA_CONSUMER_TOPIC"
	EnvForwarderElasticsearchUrl     = "BARITO_FORWARDER_ELASTICSEARCH_URL"
)

type ForwarderConfig struct {
	KafkaBrokers         []string
	KafkaConsumerGroupId string
	KafkaConsumerTopic   []string
	ElasticsearchUrl     string
}

// NewReceiverConfigByEnv
func NewForwarderConfigByEnv() (*ForwarderConfig, error) {

	kafkaBrokersList := os.Getenv(EnvForwarderKafkaBrokers)
	kafkaBrokers := []string{}
	if kafkaBrokersList == "" {
		kafkaBrokers = []string{"localhost:9092"}
	} else {
		kafkaBrokers = strings.Split(kafkaBrokersList, ",")
	}

	kafkaConsumerGroupId := os.Getenv(EnvForwarderKafkaConsumerGroupId)
	if kafkaConsumerGroupId == "" {
		err := fmt.Errorf("%s %s", EnvForwarderKafkaConsumerGroupId, "is empty")
		return nil, err
	}

	kafkaConsumerTopic := os.Getenv(EnvForwarderKafkaConsumerTopic)
	if kafkaConsumerTopic == "" {
		err := fmt.Errorf("%s %s", EnvForwarderKafkaConsumerTopic, "is empty")
		return nil, err
	}

	elasticsearchUrl := os.Getenv(EnvForwarderElasticsearchUrl)
	if elasticsearchUrl == "" {
		elasticsearchUrl = "http://localhost:9200"
	}

	config := &ForwarderConfig{
		KafkaBrokers:         kafkaBrokers,
		KafkaConsumerGroupId: kafkaConsumerGroupId,
		KafkaConsumerTopic:   []string{kafkaConsumerTopic},
		ElasticsearchUrl:     elasticsearchUrl,
	}

	return config, nil
}

// KafkaUpstreamConfig
func (c ForwarderConfig) KafkaUpstream() (river.Upstream, error) {
	return river.NewKafkaUpstream(river.KafkaUpstreamConfig{
		Brokers:         c.KafkaBrokers,
		ConsumerGroupId: c.KafkaConsumerGroupId,
		ConsumerTopic:   c.KafkaConsumerTopic,
	})
}

// ElasticsearchDownstreamConfig
func (c ForwarderConfig) ElasticsearchDownstream() (river.Downstream, error) {
	return river.NewElasticsearchDownstream(river.ElasticsearchDownstreamConfig{
		Urls: c.ElasticsearchUrl,
	})
}

func (c ForwarderConfig) Info(log *logrus.Logger) {
	log.Infof("%s=%s", EnvForwarderKafkaBrokers, c.KafkaBrokers)
	log.Infof("%s=%s", EnvForwarderKafkaConsumerGroupId, c.KafkaConsumerGroupId)
	log.Infof("%s=%s", EnvForwarderKafkaConsumerTopic, c.KafkaConsumerTopic)
	log.Infof("%s=%s", EnvForwarderElasticsearchUrl, c.ElasticsearchUrl)
}
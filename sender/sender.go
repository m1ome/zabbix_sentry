package sender

import (
	"strconv"

	"github.com/blacked/go-zabbix"
	"github.com/pkg/errors"
)

type (
	Sender struct {
		sender *zabbix.Sender
	}

	Metric struct {
		Host   string
		Metric string
		Value  int64
	}

	Options struct {
		Host string
		Port int
	}
)

func New(opts Options) Sender {
	return Sender{
		sender: zabbix.NewSender(opts.Host, opts.Port),
	}
}

func (s Sender) Send(metrics []Metric) error {
	zabbixMetrics := make([]*zabbix.Metric, len(metrics))
	for i := 0; i < len(metrics); i++ {
		v := strconv.FormatInt(metrics[i].Value, 10)
		zabbixMetrics[i] = zabbix.NewMetric(metrics[i].Host, metrics[i].Metric, v)
	}

	p := zabbix.NewPacket(zabbixMetrics)
	if _, err := s.sender.Send(p); err != nil {
		return errors.Wrap(err, "error sending metrics to zabbix")
	}

	return nil
}

package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/m1ome/zabbix_sentry/producer"
	"github.com/m1ome/zabbix_sentry/sender"
)

var (
	host             string
	sentryEntrypoint string
	sentryApiKey     string
	projects         string
	zabbixHost       string
	zabbixPort       int
	verbose          bool
)

func init() {
	flag.StringVar(&host, "host", "zabbix-sentry", "hostname will be sent to zabbix")
	flag.StringVar(&sentryEntrypoint, "sentry-url", "http://sentry.io/api/0/", "Sentry custom entrypoint")
	flag.StringVar(&sentryApiKey, "sentry-api-key", "", "Sentry api key")
	flag.StringVar(&projects, "projects", "", "projects to be filtered with, comma separated strings")
	flag.StringVar(&zabbixHost, "zabbix-host", "127.0.0.1", "Zabbix host")
	flag.IntVar(&zabbixPort, "zabbix-port", 10051, "Zabbix port")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")

	flag.Parse()
}

func main() {
	p, err := producer.New(producer.Options{
		ApiKey:   sentryApiKey,
		Endpoint: sentryEntrypoint,
	})
	if err != nil {
		log.Fatalf("error creating sentry client: %v", err)
	}

	s := sender.New(sender.Options{
		Host: zabbixHost,
		Port: zabbixPort,
	})

	for {
		stats, err := p.ProjectStats(producer.ProjectStatsQuery{
			Projects: strings.Split(projects, ", "),
		})
		if err != nil {
			log.Fatalf("error getting projects stats: %v", err)
		}

		metrics := make([]sender.Metric, 0)
		for project, stat := range stats {
			metrics = append(metrics, sender.Metric{
				Host:   host,
				Metric: project,
				Value:  int64(stat),
			})
		}

		if err := s.Send(metrics); err != nil {
			log.Fatalf("error sending metrics to zabbix: %v", err)
		}

		if verbose {
			log.Printf("send %d metrics to zabbix", len(metrics))
		}

		time.Sleep(time.Minute)
	}
}

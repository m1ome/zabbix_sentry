package producer

import (
	"fmt"
	"time"

	"github.com/atlassian/go-sentry-api"
	"github.com/pkg/errors"
)

type (
	Producer struct {
		sentryClient *sentry.Client
		resolution   string
	}

	ProjectStatsQuery struct {
		Projects []string
	}

	Options struct {
		ApiKey   string
		Endpoint string
		Timeout  int
	}
)

const (
	defaultTimeout    = 10
	defaultResolution = "1h"
)

var (
	statTypes = []sentry.StatQuery{sentry.StatReceived, sentry.StatRejected, sentry.StatBlacklisted}
)

func New(opts Options) (*Producer, error) {
	if opts.Timeout == 0 {
		opts.Timeout = defaultTimeout
	}

	client, err := sentry.NewClient(opts.ApiKey, &opts.Endpoint, &opts.Timeout)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating sentry client: %v", err)
	}

	return &Producer{
		sentryClient: client,
		resolution:   defaultResolution,
	}, err
}

func (s Producer) ProjectStats(opts ProjectStatsQuery) (map[string]float64, error) {
	projects, err := s.listProjects()
	if err != nil {
		return nil, err
	}

	stats := make(map[string]float64, 0)
	for _, project := range projects {
		if !s.isProjectFit(project, opts.Projects) {
			continue
		}

		for _, statType := range statTypes {
			stat, err := s.sentryClient.GetProjectStats(
				*project.Organization,
				project,
				statType,
				time.Now().Add(-1*time.Minute).Unix(),
				time.Now().Unix(),
				&s.resolution,
			)
			if err != nil {
				return nil, errors.Wrapf(
					err,
					"error getting project %s stats from sentry: %v",
					*project.Slug,
					err,
				)
			}

			project := fmt.Sprintf(
				"sentry.events.[%s,%s,%s]",
				*project.Organization.Slug,
				*project.Slug,
				statType,
			)

			sum := float64(0)
			for _, s := range stat {
				sum += s[1]
			}
			stats[project] = sum
		}
	}

	return stats, nil
}

func (s Producer) isProjectFit(project sentry.Project, projects []string) bool {
	if len(projects) == 0 {
		return true
	}

	for _, p := range projects {
		if *project.Slug == p {
			return true
		}
	}

	return false
}

func (s Producer) listProjects() ([]sentry.Project, error) {
	records, _, err := s.sentryClient.GetProjects()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting sentry projects: %v", err)
	}

	res := make([]sentry.Project, 0)
	res = append(res, records...)

	return res, nil
}

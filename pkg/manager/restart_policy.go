package manager

import (
	"time"

	"github.com/90poe/connectctl/pkg/client/connect"
)

type runtimeRestartPolicy map[string]Policy

const (
	// period to use between restart attempts
	// 10 seconds chosen at random
	defaultRestartPeriod = time.Second * 10
)

func runtimePolicyFromConnectors(connectors []connect.Connector, config *Config) runtimeRestartPolicy {
	// create restart policy here, overriding with any supplied values (if any)
	policy := runtimeRestartPolicy{}

	for _, c := range connectors {
		p := Policy{
			MaxConnectorRestarts:   1,
			ConnectorRestartPeriod: defaultRestartPeriod,
			MaxTaskRestarts:        1,
			TaskRestartPeriod:      defaultRestartPeriod,
		}
		if config != nil {
			// apply globals (if any)
			if config.GlobalMaxConnectorRestarts != 0 {
				p.MaxConnectorRestarts = config.GlobalMaxConnectorRestarts
			}
			if config.GlobalMaxTaskRestarts != 0 {
				p.MaxTaskRestarts = config.GlobalMaxTaskRestarts
			}
			if config.GlobalConnectorRestartPeriod != 0 {
				p.ConnectorRestartPeriod = config.GlobalConnectorRestartPeriod
			}
			if config.GlobalTaskRestartPeriod != 0 {
				p.TaskRestartPeriod = config.GlobalTaskRestartPeriod
			}
		}
		policy[c.Name] = p
	}

	// apply overrides (if any)
	if config != nil && config.RestartOverrides != nil {
		for k, v := range config.RestartOverrides.Connectors {
			p := policy[k]

			if v.MaxConnectorRestarts != 0 {
				p.MaxConnectorRestarts = v.MaxConnectorRestarts
			}
			if v.ConnectorRestartPeriod != 0 {
				p.ConnectorRestartPeriod = v.ConnectorRestartPeriod
			}
			if v.MaxTaskRestarts != 0 {
				p.MaxTaskRestarts = v.MaxTaskRestarts
			}
			if v.TaskRestartPeriod != 0 {
				p.TaskRestartPeriod = v.TaskRestartPeriod
			}

			policy[k] = p
		}
	}
	return policy
}

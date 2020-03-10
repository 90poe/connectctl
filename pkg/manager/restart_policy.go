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
			ConnectorRestartsMax:   1,
			ConnectorRestartPeriod: defaultRestartPeriod,
			TaskRestartsMax:        1,
			TaskRestartPeriod:      defaultRestartPeriod,
		}
		if config != nil {
			// apply globals (if any)
			if config.GlobalConnectorRestartsMax != 0 {
				p.ConnectorRestartsMax = config.GlobalConnectorRestartsMax
			}
			if config.GlobalTaskRestartsMax != 0 {
				p.TaskRestartsMax = config.GlobalTaskRestartsMax
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

			if v.ConnectorRestartsMax != 0 {
				p.ConnectorRestartsMax = v.ConnectorRestartsMax
			}
			if v.ConnectorRestartPeriod != 0 {
				p.ConnectorRestartPeriod = v.ConnectorRestartPeriod
			}
			if v.TaskRestartsMax != 0 {
				p.TaskRestartsMax = v.TaskRestartsMax
			}
			if v.TaskRestartPeriod != 0 {
				p.TaskRestartPeriod = v.TaskRestartPeriod
			}

			policy[k] = p
		}
	}
	return policy
}

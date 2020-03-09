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

func runtimePolicyFromConnectors(connectors []connect.Connector, overrides *RestartPolicy) runtimeRestartPolicy {
	// create restart policy here, overriding with any supplied values (if any)
	policy := runtimeRestartPolicy{}

	for _, c := range connectors {
		policy[c.Name] = Policy{
			MaxConnectorRestarts:   1,
			ConnectorRestartPeriod: defaultRestartPeriod,
			MaxTaskRestarts:        1,
			TaskRestartPeriod:      defaultRestartPeriod,
		}
	}

	// apply overrides
	if overrides != nil {
		for k, v := range overrides.Connectors {
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

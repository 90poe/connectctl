package connectors

import (
	"github.com/90poe/connectctl/internal/version"
	"github.com/90poe/connectctl/pkg/manager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type resumeConnectorsCmdParams struct {
	ClusterURL string
	Connectors []string
}

func resumeConnectorsCmd() *cobra.Command {
	params := &resumeConnectorsCmdParams{}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume connectors in a cluster",
		Long:  "",
		Run: func(cmd *cobra.Command, _ []string) {
			doResumeConnectors(cmd, params)
		},
	}

	addCommonConnectorsFlags(resumeCmd, &params.ClusterURL)
	addConnectorNamesFlags(resumeCmd, &params.Connectors)

	return resumeCmd
}

func doResumeConnectors(_ *cobra.Command, params *resumeConnectorsCmdParams) {
	clusterLogger := log.WithField("cluster", params.ClusterURL)
	clusterLogger.Infof("resuming onnectots: %s", params.Connectors)

	config := &manager.Config{
		ClusterURL: params.ClusterURL,
		Version:    version.Version,
	}
	clusterLogger.WithField("config", config).Trace("resume connectors configuration")

	mngr, err := manager.NewConnectorsManager(config)
	if err != nil {
		clusterLogger.WithError(err).Fatalln("error creating connectors manager")
	}

	err = mngr.Resume(params.Connectors)
	if err != nil {
		clusterLogger.WithError(err).Fatal("error resuming connectors")
	}

	clusterLogger.Info("connectors resumed successfully")
}

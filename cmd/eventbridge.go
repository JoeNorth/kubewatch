/*
Copyright 2018 Bitnami

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// eventbridgeConfigCmd represents the eventbridge subcommand
var eventbridgeConfigCmd = &cobra.Command{
	Use:   "eventbridge",
	Short: "specific eventbridge configuration",
	Long:  `specific eventbridge configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		endpointId, err := cmd.Flags().GetString("endpointid")
		if err == nil {
			if len(endpointId) > 0 {
				conf.Handler.EventBridge.EndpointId = endpointId
			}
		} else {
			logrus.Fatal(err)
		}

		clusterArn, err := cmd.Flags().GetString("clusterArn")
		if err == nil {
			if len(clusterArn) > 0 {
				conf.Handler.EventBridge.ClusterArn = clusterArn
			}
		} else {
			logrus.Fatal(err)
		}

		eventBusName, err := cmd.Flags().GetString("eventBusName")
		if err == nil {
			if len(eventBusName) > 0 {
				conf.Handler.EventBridge.EventBusName = eventBusName
			}
		} else {
			logrus.Fatal(err)
		}

		if err = conf.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	eventbridgeConfigCmd.Flags().StringP("endpointid", "e", "", "Specify EventBridge endpoint id (optional)")
	eventbridgeConfigCmd.Flags().StringP("clusterArn", "c", "", "Specify EKS cluster ARN  (optional)")
	eventbridgeConfigCmd.Flags().StringP("eventBusName", "b", "", "Specify EventBridge event bus name (optional)")
}

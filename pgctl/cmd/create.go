/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a PostgreSQL cluster",
	Long:  `Creates a PostgreSQL cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		// var kubeconfig *string
		// if home := homeDir(); home != "" {
		// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		// } else {
		// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		// }
		// flag.Parse()

		// use the current context in kubeconfig
		// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		// if err != nil {
		// 	panic(err.Error())
		// }

		// create the clientset
		// clientset, err := kubernetes.NewForConfig(config)
		// if err != nil {
		// 	panic(err.Error())
		// }

		// clusterdefinition := &ext.CustomResourceDefinition{
		// 	TypeMeta: v1.TypeMeta{
		// 		APIVersion: "postgresql.iboware.com/v1alpha1",
		// 		Kind:       "PostgreSQL",
		// 	},
		// 	ObjectMeta: v1.ObjectMeta{
		// 		Name: "",
		// 	},
		// }

		// 	apiVersion: postgresql.iboware.com/v1alpha1
		// kind: PostgreSQL
		// metadata:
		//   name: iboware
		//   namespace: default
		// spec:
		//   # Add fields here
		//   replicas: 3
		//   disksize: "4Gi"
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().StringP("name", "n", "postgresql", "Name of the cluster")
	createCmd.Flags().IntP("replicas", "r", 3, "Amount of Replicas")
	createCmd.MarkFlagRequired("name")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

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
	"context"
	"fmt"
	"log"
	"os"

	v1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var replicas int32

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a PostgreSQL cluster",
	Long:  `Creates a PostgreSQL cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		scheme := runtime.NewScheme()
		_ = clientgoscheme.AddToScheme(scheme)
		_ = v1alpha1.AddToScheme(scheme)

		kubeconfig := ctrl.GetConfigOrDie()
		kubeclient, err := client.New(kubeconfig, client.Options{Scheme: scheme})
		if err != nil {
			log.Fatal(err)
		}
		name, _ := cmd.Flags().GetString("name")
		replicas, _ := cmd.Flags().GetInt32("replicas")
		storage, _ := cmd.Flags().GetString("storage")
		namespace, _ := cmd.Flags().GetString("namespace")

		error := kubeclient.Create(context.Background(), &v1alpha1.PostgreSQL{
			ObjectMeta: v1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: v1alpha1.PostgreSQLSpec{Replicas: replicas, DiskSize: storage},
		})
		if error != nil {
			log.Panic(error)
		}
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
	createCmd.Flags().StringP("name", "n", "", "Name of the Cluster")
	createCmd.Flags().Int32P("replicas", "r", 3, "Amount of Replicas")
	createCmd.Flags().StringP("disksize", "d", "8Gi", "Disk Size per Replica")
	createCmd.Flags().StringP("namespace", "ns", "default", "Namespace")

	createCmd.MarkFlagRequired("name")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

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

	v1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [cluster name]",
	Short: "Deletes an existing PostgreSQL cluster",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("delete needs a name for the cluster")
		}
		name := args[0]
		namespace, _ := cmd.Flags().GetString("namespace")

		scheme := runtime.NewScheme()
		_ = clientgoscheme.AddToScheme(scheme)
		_ = v1alpha1.AddToScheme(scheme)

		kubeconfig := ctrl.GetConfigOrDie()
		kubeclient, err := client.New(kubeconfig, client.Options{Scheme: scheme})
		if err != nil {
			log.Fatal(err)
		}

		err2 := kubeclient.Delete(context.Background(), &v1alpha1.PostgreSQL{
			ObjectMeta: v1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			}})
		if err2 != nil {
			log.Fatal(err2)
		} else {
			fmt.Printf("Cluster %v under Namespace:%v has been succesfully deleted.\n", name, namespace)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("namespace", "n", "default", "Namespace")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

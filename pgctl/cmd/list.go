/*
Copyright © 2020 IBRAHIM VAROL <iboware@gmail.com>

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
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists PostgreSQL Clusters",
	Long:  `Lists PostgreSQL Clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		scheme := runtime.NewScheme()
		_ = clientgoscheme.AddToScheme(scheme)
		_ = v1alpha1.AddToScheme(scheme)

		kubeconfig := ctrl.GetConfigOrDie()
		kubeclient, err := client.New(kubeconfig, client.Options{Scheme: scheme})
		if err != nil {
			log.Fatal(err)
		}
		list := v1alpha1.PostgreSQLList{}
		err2 := kubeclient.List(context.Background(), &list)
		if err2 != nil {
			log.Fatal(err2)
		}

		fmt.Printf("%-20v %-10v %-20v\n", "Name", "Replicas", "Namespace")

		for _, item := range list.Items {
			fmt.Printf("%-20v %-10v %-20v\n", item.Name, item.Spec.Replicas, item.Namespace)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

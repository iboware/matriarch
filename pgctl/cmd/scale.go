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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command{
	Use:   "scale [cluster name]",
	Short: "Scales a cluster",
	Long:  `Scales a cluster with given parameters up or down`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("scale needs a name for the cluster")
		}

		name := args[0]
		replicas, _ := cmd.Flags().GetInt32("replicas")
		namespace, _ := cmd.Flags().GetString("namespace")

		scheme := runtime.NewScheme()
		_ = clientgoscheme.AddToScheme(scheme)
		_ = v1alpha1.AddToScheme(scheme)

		kubeconfig := ctrl.GetConfigOrDie()
		kubeclient, err := client.New(kubeconfig, client.Options{Scheme: scheme})
		if err != nil {
			log.Fatal(err)
		}

		cluster := v1alpha1.PostgreSQL{
			ObjectMeta: v1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}

		errGet := kubeclient.Get(context.Background(),
			types.NamespacedName{Name: name, Namespace: namespace},
			&cluster)

		if errGet != nil {
			log.Fatal(errGet)
		}
		cluster.Spec.Replicas = replicas

		errUpd := kubeclient.Update(context.Background(),
			&cluster)

		if errUpd != nil {
			log.Fatal(errUpd)
		} else {
			fmt.Printf("Cluster %v under Namespace:%v has been succesfully scaled to %d replicas.\n", name, namespace, replicas)
		}
	},
}

func init() {
	rootCmd.AddCommand(scaleCmd)
	scaleCmd.Flags().StringP("namespace", "n", "default", "Namespace")
	scaleCmd.Flags().Int32P("replicas", "r", 3, "Replica amount")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scaleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scaleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

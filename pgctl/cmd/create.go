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
	"strings"

	v1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	"github.com/sethvargo/go-password/password"
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
	Use:   "create [cluster name]",
	Short: "Creates a PostgreSQL cluster",
	Long:  `Creates a PostgreSQL cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		if len(args) < 1 {
			log.Fatal("create needs a name for the cluster")
		}
		var generatedPass bool = false
		scheme := runtime.NewScheme()
		_ = clientgoscheme.AddToScheme(scheme)
		_ = v1alpha1.AddToScheme(scheme)

		kubeconfig := ctrl.GetConfigOrDie()
		kubeclient, err := client.New(kubeconfig, client.Options{Scheme: scheme})
		if err != nil {
			log.Fatal(err)
		}
		name := args[0]
		replicas, _ := cmd.Flags().GetInt32("replicas")
		storage, _ := cmd.Flags().GetString("disksize")
		namespace, _ := cmd.Flags().GetString("namespace")
		postgrespassword, _ := cmd.Flags().GetString("postgrespassword")
		repmgrspassword, _ := cmd.Flags().GetString("repmgrpassword")

		if len(strings.TrimSpace(postgrespassword)) == 0 {
			postgrespassword, _ = password.Generate(32, 10, 10, false, false)
			generatedPass = true
		}

		if len(strings.TrimSpace(repmgrspassword)) == 0 {
			repmgrspassword = postgrespassword
		}

		error := kubeclient.Create(context.Background(), &v1alpha1.PostgreSQL{
			ObjectMeta: v1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: v1alpha1.PostgreSQLSpec{
				Replicas:         replicas,
				DiskSize:         storage,
				PostgresPassword: postgrespassword,
				RepMGRPassword:   repmgrspassword,
			},
		})
		if error != nil {
			log.Panic(error)
		}

		if generatedPass {
			fmt.Printf("PostgreSQL auto-generated password: %v/n", postgrespassword)
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
	createCmd.Flags().Int32P("replicas", "r", 3, "Amount of Replicas")
	createCmd.Flags().StringP("disksize", "d", "8Gi", "Disk Size per Replica")
	createCmd.Flags().StringP("namespace", "n", "default", "Namespace")
	createCmd.Flags().StringP("postgrespassword", "p", "", "Password for default PostgreSql user")
	createCmd.Flags().StringP("repmgrpassword", "m", "", "Password for RepMGR. If not specified PostgreSQL password will be used.")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

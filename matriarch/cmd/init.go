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
	"strings"

	"io/ioutil"

	v1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	"github.com/iboware/postgresql-operator/matriarch/utils"
	"github.com/spf13/cobra"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Deploys operator into the active cluster.",
	Long:  `Deploys operator into the active cluster, which is specified in your local kubeconfig file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deploying the operator...")

		err := utils.DownloadFile("https://github.com/iboware/postgresql-operator/releases/download/v0.3.7/postgresql-operator.crd.yaml", "/var/tmp/postgresql-operator.crd.yaml")

		if err != nil {
			log.Fatal(err)
			return
		}

		operatorFile, _ := ioutil.ReadFile("/var/tmp/postgresql-operator.crd.yaml")
		files := strings.Split(string(operatorFile), "---")
		sch := runtime.NewScheme()
		_ = clientgoscheme.AddToScheme(sch)
		_ = apiextv1beta1.AddToScheme(sch)
		_ = v1alpha1.AddToScheme(sch)

		kubeconfig := ctrl.GetConfigOrDie()
		kubeclient, err := client.New(kubeconfig, client.Options{Scheme: sch})

		if err != nil {
			log.Fatal(err)
			return
		}

		decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode
		for _, file := range files {

			obj, _, err := decode([]byte(file), nil, nil)
			if err != nil {
				log.Fatal(err)
				return
			}

			err = kubeclient.Create(context.Background(), obj)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		fmt.Println("Operator Deployed Successfully!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

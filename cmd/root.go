// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	"github.com/runyontr/kaudit/pkg/discovery"
	//Talk to GKE
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
)

var cfgFile string
var kubeconfig *string
var clientset *kubernetes.Clientset

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kaudit",
	Short: "Audit resources in Kubernetes against jsonspec",
	Long: `There are suggested labels and annotations for applications and resources
that can be suggested and required by either the k8s community, or by an organization


This tool looks to be able to consume a json spec and validate either all instances
of all resources adhere to the json spec, or be able to validate a specific
resource type adheres to the provided spec.


e.g.

# Audit all resources (TODO, what does "all" mean?  Lets start with just things in the Workloads API)
kaudit --spec allspec.json

#Just audit pods
kaudit pods --spec podsspec.json


#Limit to just v1 and betav2 apis
kaudit --spec allspec.json --version v1,betav2


`,
// Uncomment the following line if your bare application
// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		//we'll use this in the future to find all resource types
		_, err := discovery.GetResourceTypes(clientset)
		if err != nil{
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}


		//Load the schema to validate against
		//TODO (@trunyon) load spec from command line flag
		schemaLoader := gojsonschema.NewReferenceLoader(os.ExpandEnv("file://$PWD/spec.json"))
		schema, err := gojsonschema.NewSchema(schemaLoader)
		if err != nil{
			panic(err)
		}

		//Hard coded deployment search as first look at validating objects
		//TODO (@trunyon) load namespace from command line flag
		deps, err := clientset.AppsV1beta2().Deployments("default").List(v1.ListOptions{})
		if err != nil{
			fmt.Printf("Error getting deployments: %v\n", err)
			os.Exit(1)
		}

		for _, d := range deps.Items{
			b, _ := json.Marshal(d)
			//check each against spec.json
			documentLoader :=  gojsonschema.NewStringLoader(string(b))

			result, err := schema.Validate(documentLoader)
			if err != nil {
				panic(err.Error())
			}

			if result.Valid() {
				fmt.Printf("The document is valid\n")
			} else {
				fmt.Printf("Deployment: %v", d.Name)
				fmt.Printf("The document is not valid. see errors :\n")
				for _, desc := range result.Errors() {
					fmt.Printf("- %s\n", desc)
				}
			}

		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kaudit.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().StringP("namespace", "n", "default", "Limit search to provided namespace")
	RootCmd.Flags().StringP("spec", "s","spec.json","JSON Spec to use")
	if home := homeDir(); home != "" {
		kubeconfig = RootCmd.Flags().String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = RootCmd.Flags().String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".kaudit") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}


	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
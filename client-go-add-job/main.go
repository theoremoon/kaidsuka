package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	batchv1 "k8s.io/api/batch/v1"

	/*
		v1 "k8s.io/api/core/v1"
		metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	*/
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const NameSpaceName = "default"

func run() error {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <Job yaml>\n", os.Args[0])
		return nil
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return err
	}
	var job batchv1.Job
	if err := yaml.Unmarshal(data, &job); err != nil {
		return err
	}

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	/*
		ns := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: NameSpaceName},
		}
		if _, err := clientset.CoreV1().Namespaces().Create(&ns); err != nil {
			return err
		}
	*/

	jobsApi := clientset.BatchV1().Jobs(NameSpaceName)
	if _, err := jobsApi.Create(&job); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

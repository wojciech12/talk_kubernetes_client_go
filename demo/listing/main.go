package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeConfig *string
	home := os.Getenv("HOME")
	if home == "" {
		panic("HOME not found")
	}

	kubeConfig = flag.String("kubeconfig",
		filepath.Join(home, ".kube", "config"),
		"(optional) absolute path to the kubeconfig file")

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		check(clientSet, "kube-system", "etcd-minikube")
		time.Sleep(10 * time.Second)
	}
}

func check(clientSet kubernetes.Interface, namespace string, pod string) bool {
	pods, err := clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	_, err = clientSet.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		return true
	}
	return false

}

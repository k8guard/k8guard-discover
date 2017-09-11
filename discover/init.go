package discover

import (
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/k8s"
	"k8s.io/client-go/kubernetes"
)

// Clientset talks to kubernetes API
var Clientset *kubernetes.Clientset
var err error

func init() {
	Clientset, err = k8s.LoadClientset()

	if err != nil {
		lib.Log.Error("error loading ClientSet ", err)
		panic(err)
	}
}

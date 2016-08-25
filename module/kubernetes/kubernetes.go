package kubernetes

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	namespaceLock *sync.RWMutex
	k8sClient     *unversioned.Client
)

const (
	// K8sClientErr Kubernetes client error
	K8sClientErr = "Kubernetes client is not start"
)

func init() {
	namespaceLock = new(sync.RWMutex)
}

func getClient() error {
	if k8sClient == nil {
		k8sClient = models.K8sClient
	}
	if k8sClient == nil {
		return k8sClientErr()
	}
	return nil
}

func k8sClientErr() error {
	return errors.New(K8sClientErr)
}

func formatResult(result string, detail string) *models.K8SRes {
	log.Println(fmt.Sprintf("Stage in k8s result is %v, detail is %v", result, detail))
	return &models.K8SRes{
		Result: result,
		Detail: detail,
	}
}

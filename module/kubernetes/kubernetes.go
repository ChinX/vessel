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
	k8s           *unversioned.Client
)

const (
	// K8sClientErr Kubernetes client not connected error
	K8sClientErr = "Kubernetes client is not connected"
)

func init() {
	namespaceLock = new(sync.RWMutex)
}

func checkClient() error {
	if models.K8S == nil {
		if err := models.InitK8S(); err != nil {
			return err
		}
	}
	k8s = models.K8S
	if k8s == nil {
		return errors.New(K8sClientErr)
	}
	return nil
}

func formatResult(result string, detail string) *models.K8SRes {
	log.Println(fmt.Sprintf("Stage in k8s result is %v, detail is %v", result, detail))
	return &models.K8SRes{
		Result: result,
		Detail: detail,
	}
}

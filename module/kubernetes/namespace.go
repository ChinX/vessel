package kubernetes

import (
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
)

func createNamespace(stage *models.Stage) error {
	if err := checkClient(); err != nil {
		return err
	}
	k8sNamespace := k8s.Namespaces()
	namespaceLock.RLock()
	if _, err := k8sNamespace.Get(stage.Namespace); err != nil {
		namespaceLock.RUnlock()
		namespaceLock.Lock()
		if _, err := k8sNamespace.Get(stage.Namespace); err != nil {
			namespaceObj := &api.Namespace{
				ObjectMeta: api.ObjectMeta{
					Name:   stage.Namespace,
					Labels: map[string]string{},
				},
			}
			namespaceObj.SetLabels(map[string]string{models.LabelKey: stage.PipelineName})

			if _, err := k8sNamespace.Create(namespaceObj); err != nil {
				namespaceLock.Unlock()
				return err
			}
		}
		namespaceLock.Unlock()
	} else {
		namespaceLock.RUnlock()
	}
	return nil
}

func deleteNamespace(stage *models.Stage) error {
	if err := checkClient(); err != nil {
		return err
	}
	k8sNamespace := k8s.Namespaces()
	namespaceLock.RLock()
	if _, err := k8sNamespace.Get(stage.Namespace); err == nil {
		namespaceLock.RUnlock()
		namespaceLock.Lock()
		if _, err := k8sNamespace.Get(stage.Namespace); err == nil {
			if err := k8sNamespace.Delete(stage.Namespace); err != nil {
				namespaceLock.Unlock()
				return err
			}
		}
		namespaceLock.Unlock()
	} else {
		namespaceLock.RUnlock()
	}
	return nil
}

func watchDeleteNamespace(stage *models.Stage, namespaceCh chan error) {
	k8sNamespace := k8s.Namespaces()
	timeChan := time.After(time.Duration(stage.Hourglass.GetLeftNanoseconds()))
	running := true
	for running {
		if _, err := k8sNamespace.Get(stage.Namespace); err != nil {
			namespaceCh <- nil
			return
		}
		select {
		case <-time.After(time.Duration(1) * time.Second):
		case <-timeChan:
			running = false
		}

	}
	namespaceCh <- fmt.Errorf("Unexpected err when watch namespace : name = %v", stage.Namespace)
	return
}

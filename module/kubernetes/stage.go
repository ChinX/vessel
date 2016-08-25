package kubernetes

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
)

// CreateStage from kubernetes
func CreateStage(stage *models.Stage) (res *models.K8SRes) {
	log.Printf("Create statge in k8s stage name = %v", stage.Name)
	if stage.Hourglass.GetLeftNanoseconds() < 0 {
		return formatResult(models.ResultTimeout, "Start stage in kubernetes timeout")
	}
	has, err := checkRC(stage)
	if err != nil {
		return formatResult(models.ResultFailed, err.Error())
	} else if has {
		return formatResult(models.ResultFailed, fmt.Sprintf("Replication controller %v already exist", stage.Name))
	}
	if err := createNamespace(stage); err != nil {
		return formatResult(models.ResultFailed, err.Error())
	}

	ch := make(chan *models.K8SRes)
	go watchServiceStatus(stage, models.WatchAdded, ch)
	go watchRCStatus(stage, models.WatchAdded, ch)
	go watchPodStatus(stage, models.WatchAdded, ch)

	if err := createService(stage); err != nil {
		return formatResult(models.ResultFailed, err.Error())
	}
	if err := createRC(stage); err != nil {
		return formatResult(models.ResultFailed, err.Error())
	}

	for count := 3; count > 0; count-- {
		select {
		case res = <-ch:
			log.Println("Watch res :", res)
			if res.Result != models.ResultSuccess {
				return res
			}
		}
	}
	return res
}

// DeleteStage from kubernetes
func DeleteStage(stage *models.Stage) (res *models.K8SRes) {
	log.Printf("Delete statge in k8s stage name = %v", stage.Name)
	has, err := checkRC(stage)
	if err != nil {
		return formatResult(models.ResultFailed, err.Error())
	} else if !has {
		return formatResult(models.ResultFailed, fmt.Sprintf("Replication controller %v not start", stage.Name))
	}

	ch := make(chan *models.K8SRes)
	go watchPodStatus(stage, models.WatchDeleted, ch)

	if err := deleteService(stage); err != nil {
		return formatResult(models.ResultFailed, err.Error())
	}
	if err := deleteRC(stage); err != nil {
		return formatResult(models.ResultFailed, err.Error())
	}
	if err := deletePods(stage); err != nil {
		return formatResult(models.ResultFailed, err.Error())
	}

	if count, err := getRCCount(stage); err == nil && count == 0 {
		namespaceCh := make(chan error)
		go watchDeleteNamespace(stage, namespaceCh)
		deleteNamespace(stage)
		<-namespaceCh
	}
	return <-ch
}

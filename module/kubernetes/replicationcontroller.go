package kubernetes

import (
	"fmt"
	"log"
	"time"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

func createRC(stage *models.Stage) error {
	if err := getClient(); err != nil {
		return err
	}
	rc := &api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Labels: map[string]string{},
		},
		Spec: api.ReplicationControllerSpec{
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: map[string]string{},
				},
			},
			Selector: map[string]string{},
		},
	}
	rc.Spec.Template.Spec.Containers = make([]api.Container, 1)
	rc.SetName(stage.Name)
	rc.SetNamespace(stage.Namespace)
	rc.Labels[models.LabelKey] = stage.Name
	rc.Spec.Replicas = int32(stage.Replicas)
	rc.Spec.Template.SetName(stage.Name)
	rc.Spec.Template.Labels[models.LabelKey] = stage.Name
	rc.Spec.Template.Spec.Containers[0] = api.Container{
		Ports: []api.ContainerPort{
			api.ContainerPort{
				Name:          stage.Name,
				ContainerPort: int32(stage.Port),
			},
		},
		Name:            stage.Name,
		Image:           stage.Image,
		ImagePullPolicy: "IfNotPresent",
	}
	if stage.EnvName != "" && stage.EnvValue != "" {
		rc.Spec.Template.Spec.Containers[0].Env = []api.EnvVar{
			api.EnvVar{
				Name:  stage.EnvName,
				Value: stage.EnvValue,
			},
		}
	}
	rc.Spec.Selector[models.LabelKey] = stage.Name
	if _, err := k8sClient.ReplicationControllers(stage.Namespace).Create(rc); err != nil {
		log.Println("Create rc err :", err)
		return err
	}
	return nil
}

func deleteRC(stage *models.Stage) error {
	if err := getClient(); err != nil {
		return err
	}
	return k8sClient.ReplicationControllers(stage.Namespace).Delete(stage.Name)
}

func checkRC(stage *models.Stage) (bool, error) {
	if err := getClient(); err != nil {
		return false, err
	}
	rcs, err := k8sClient.ReplicationControllers(stage.Namespace).List(api.ListOptions{})
	if err != nil {
		return false, nil
	}
	for _, rc := range rcs.Items {
		if rc.Name == stage.Name {
			return true, nil
		}
	}
	return false, nil
}

func watchRCStatus(stage *models.Stage, checkOp string, ch chan *models.K8SRes) {
	if err := getClient(); err != nil {
		ch <- formatResult(models.ResultFailed, err.Error())
		return
	}
	if checkOp != string(watch.Added) && checkOp != string(watch.Deleted) {
		ch <- formatResult(models.ResultFailed, fmt.Sprintf("Unexpected err when watch RC : name = %v", stage.Name))
		return
	}
	if stage.Hourglass.GetLeftNanoseconds() <= 0 {
		ch <- formatResult(models.ResultTimeout, fmt.Sprintf("Watch RC insterface timeout when name = %v", stage.Name))
		return
	}

	opts := api.ListOptions{LabelSelector: labels.Set{models.LabelKey: stage.Name}.AsSelector()}
	w, err := k8sClient.ReplicationControllers(stage.Namespace).Watch(opts)
	if err != nil {
		ch <- formatResult(models.ResultFailed, err.Error())
		return
	}
	timeChan := time.After(time.Duration(stage.Hourglass.GetLeftNanoseconds()))
	select {
	case event, ok := <-w.ResultChan():
		if !ok {
			ch <- formatResult(models.ResultFailed, fmt.Sprintf("Unexpected err when watch RC : name = %v", stage.Name))
			w.Stop()
			return
		}
		if string(event.Type) == checkOp {
			ch <- formatResult(models.ResultSuccess, "")
			w.Stop()
			return
		}
	case <-timeChan:
		ch <- formatResult(models.ResultTimeout, fmt.Sprintf("Watch RC insterface timeout when name = %v", stage.Name))
		w.Stop()
		return
	}
}

func getRCCount(stage *models.Stage) (int, error) {
	if err := getClient(); err != nil {
		return 0, err
	}
	rcs, err := k8sClient.ReplicationControllers(stage.Namespace).List(api.ListOptions{})
	if err != nil {
		return 0, err
	}
	return len(rcs.Items), nil
}

package kubernetes

import (
	"fmt"
	"log"
	"time"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/watch"
)

func createService(stage *models.Stage) error {
	if err := getClient(); err != nil {
		return err
	}
	service := &api.Service{
		ObjectMeta: api.ObjectMeta{
			Labels: map[string]string{},
		},
		Spec: api.ServiceSpec{
			Selector: map[string]string{},
		},
	}
	service.Spec.Ports = make([]api.ServicePort, 1)
	service.ObjectMeta.SetName(stage.Name)
	service.ObjectMeta.SetNamespace(stage.Namespace)
	service.ObjectMeta.Labels[models.LabelKey] = stage.Name
	service.Spec.Ports[0] = api.ServicePort{Port: int32(stage.Port), TargetPort: intstr.FromString(stage.Name)}
	service.Spec.Selector[models.LabelKey] = stage.Name
	if _, err := k8sClient.Services(stage.Namespace).Create(service); err != nil {
		log.Println("Create service err :", err)
		return err
	}
	return nil
}

func deleteService(stage *models.Stage) error {
	if err := getClient(); err != nil {
		return err
	}
	return k8sClient.Services(stage.Namespace).Delete(stage.Name)
}

func watchServiceStatus(stage *models.Stage, checkOp string, ch chan *models.K8SRes) {
	if err := getClient(); err != nil {
		ch <- formatResult(models.ResultFailed, err.Error())
		return
	}
	if checkOp != string(watch.Added) && checkOp != string(watch.Deleted) {
		ch <- formatResult(models.ResultFailed, fmt.Sprintf("Unexpected err when watch service : name = %v", stage.Name))
		return
	}
	if stage.Hourglass.GetLeftNanoseconds() <= 0 {
		ch <- formatResult(models.ResultTimeout, fmt.Sprintf("Watch service insterface timeout when name = %v", stage.Name))
		return
	}

	opts := api.ListOptions{LabelSelector: labels.Set{models.LabelKey: stage.Name}.AsSelector()}
	w, err := k8sClient.Services(stage.Namespace).Watch(opts)
	if err != nil {
		ch <- formatResult(models.ResultFailed, err.Error())
		w.Stop()
		return
	}
	timeChan := time.After(time.Duration(stage.Hourglass.GetLeftNanoseconds()))
	select {
	case event, ok := <-w.ResultChan():
		if !ok {
			ch <- formatResult(models.ResultFailed, fmt.Sprintf("Unexpected err when watch service : name = %v", stage.Name))
			w.Stop()
			return
		}
		if string(event.Type) == checkOp {
			ch <- formatResult(models.ResultSuccess, "")
			w.Stop()
			return
		}
	case <-timeChan:
		ch <- formatResult(models.ResultTimeout, fmt.Sprintf("Watch service insterface timeout when name = %v", stage.Name))
		w.Stop()
		return
	}
}

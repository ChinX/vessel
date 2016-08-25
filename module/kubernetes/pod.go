package kubernetes

import (
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

func watchPodStatus(stage *models.Stage, checkOp string, ch chan *models.K8SRes) {
	if err := getClient(); err != nil {
		ch <- formatResult(models.ResultFailed, err.Error())
		return
	}
	if checkOp != string(watch.Added) && checkOp != string(watch.Deleted) {
		ch <- formatResult(models.ResultFailed, fmt.Sprintf("Unexpected err when watch pod : name = %v", stage.Name))
		return
	}
	if stage.Hourglass.GetLeftNanoseconds() <= 0 {
		ch <- formatResult(models.ResultTimeout, fmt.Sprintf("Watch pod insterface timeout when name = %v", stage.Name))
		return
	}

	opts := api.ListOptions{LabelSelector: labels.Set{models.LabelKey: stage.Name}.AsSelector()}
	w, err := k8sClient.Pods(stage.Namespace).Watch(opts)
	if err != nil {
		ch <- formatResult(models.ResultFailed, err.Error())
		return
	}

	if checkOp == string(watch.Added) {
		checkOp = string(watch.Modified)
	}

	timeChan := time.After(time.Duration(stage.Hourglass.GetLeftNanoseconds()))
	sum := int(stage.Replicas)
	for count := 0; count < sum; {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				ch <- formatResult(models.ResultFailed, fmt.Sprintf("Unexpected err when watch pod : name = %v", stage.Name))
				w.Stop()
				return
			}
			if string(event.Type) != checkOp || event.Object.(*api.Pod).Status.Phase != api.PodRunning {
				continue
			}
			count++
		case <-timeChan:
			ch <- formatResult(models.ResultTimeout, fmt.Sprintf("Watch pod insterface timeout when name = %v", stage.Name))
			w.Stop()
			return
		}
	}
	ch <- formatResult(models.ResultSuccess, "")
	w.Stop()
}

func deletePods(stage *models.Stage) error {
	if err := getClient(); err != nil {
		return err
	}
	opts := api.ListOptions{LabelSelector: labels.Set{models.LabelKey: stage.Name}.AsSelector()}
	pods, err := k8sClient.Pods(stage.Namespace).List(opts)
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		delErr := k8sClient.Pods(stage.Namespace).Delete(pod.ObjectMeta.Name, &api.DeleteOptions{})
		if delErr != nil {
			err = delErr
		}
	}
	return err
}

func getPodIPList(stage *models.Stage) ([]string, error) {
	if err := getClient(); err != nil {
		return nil, err
	}
	opts := api.ListOptions{LabelSelector: labels.Set{models.LabelKey: stage.Name}.AsSelector()}
	pods, err := k8sClient.Pods(stage.Namespace).List(opts)
	if err != nil {
		return nil, err
	}
	ipList := make([]string, 0, len(pods.Items))
	for _, pod := range pods.Items {
		ipList = append(ipList, pod.Status.PodIP)
	}
	return ipList, nil
}

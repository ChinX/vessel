package scheduler

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/point"
)

type schedulerHand func(interface{}, map[string]bool, chan *models.ExecutedResult) bool

// StartStage start point on scheduler
func StartPoint(executorMap map[string]interface{}, startMark string) []*models.ExecutedResult {
	return execute(executorMap, startMark, point.Start)
}

// StopPoint stop point on scheduler
func StopPoint(executorMap map[string]interface{}, startMark string) []*models.ExecutedResult {
	return execute(executorMap, startMark, point.Stop)
}

func execute(executorMap map[string]interface{}, startMark string, handler schedulerHand) []*models.ExecutedResult {
	count := len(executorMap)
	readyMap := map[string]bool{startMark: true}
	finishChan := make(chan *models.ExecutedResult, count)
	resultList := make([]*models.ExecutedResult, 0, count)
	running := true
	for running {
		for name, executor := range executorMap {
			if _, ok := readyMap[name]; ok {
				continue
			}
			if !handler(executor, readyMap, finishChan) {
				continue
			}
			readyMap[name] = false
		}
		result := <-finishChan
		resultList = append(resultList, result)
		if result.Status != models.ResultSuccess {
			running = false
		} else {
			readyMap[result.Name] = true
			resultLen := len(resultList)
			running = resultLen != count
			log.Println(fmt.Sprintf("scheduler StartStage name = %v and finish num = %d", result.Name, resultLen))
		}
	}
	return resultList
}

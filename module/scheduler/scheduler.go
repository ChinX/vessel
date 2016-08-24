package scheduler

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/point"
)

type schedulerHand func(info interface{}, finishChan chan *models.ExecutedResult)

// StartStage start stage on scheduler
func StartPoint(executorMap map[string]*models.Point) []*models.ExecutedResult {
	readyMap := map[string]bool{"": true}

	count := len(executorMap)
	finishChan := make(chan *models.ExecutedResult, count)
	resultList := make([]*models.ExecutedResult, 0, count)

	running := true
	for running {
		go execProgress(executorMap, readyMap, finishChan, point.Start)
		result := <-finishChan
		resultList = append(resultList, result)
		if result.Status != models.ResultSuccess {
			running = false
		} else {
			readyMap[result.Name] = true
			running = len(resultList) != count
			log.Println(fmt.Sprintf("scheduler StartPoint name = %v and finish num = %d", result.Name, len(resultList)))
		}
	}
	return resultList
}

// StopStage stop stage on scheduler
func StopPoint(executorMap map[string]*models.Point) []*models.ExecutedResult {
	readyMap := map[string]bool{"": true}

	count := len(executorMap)
	finishChan := make(chan *models.ExecutedResult, count)
	resultList := make([]*models.ExecutedResult, 0, count)

	running := true
	for running {
		go execProgress(executorMap, readyMap, finishChan, point.Stop)
		result := <-finishChan
		resultList = append(resultList, result)
		if result.Status != models.ResultSuccess {
			running = false
		} else {
			readyMap[result.Name] = true
			running = len(resultList) != count
			log.Println(fmt.Sprintf("scheduler StopPoint name = %v and finish num = %d", result.Name, len(resultList)))
		}
	}
	return resultList
}

func execProgress(executorMap map[string]*models.Point, readyMap map[string]bool, finishChan chan *models.ExecutedResult, handler schedulerHand) {
	log.Println("Scheduler ready map is ", readyMap)
	for name, executor := range executorMap {
		if _, ok := readyMap[name]; ok {
			continue
		}
		isReady := true
		for _, from := range executor.From {
			if isReady, _ = readyMap[from]; !isReady {
				break
			}
		}
		if !isReady {
			continue
		}
		readyMap[name] = false
		go handler(executor.Info, finishChan)
	}
}

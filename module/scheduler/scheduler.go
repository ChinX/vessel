package scheduler

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
)

type schedulerHand func(executor models.Executor, readyMap map[string]bool, finishChan chan *models.ExecutedResult) bool

func Start(executorMap map[string]models.Executor, startMark string) []*models.ExecutedResult {
	return execute(executorMap, startMark, startProgress)
}

func Stop(executorMap map[string]models.Executor, startMark string) []*models.ExecutedResult {
	return execute(executorMap, startMark, stopProgress)
}

func execute(executorMap map[string]models.Executor, startMark string, handler schedulerHand) []*models.ExecutedResult {
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

func startProgress(executor models.Executor, readyMap map[string]bool, finishChan chan *models.ExecutedResult) bool {
	return executor.Start(readyMap, finishChan)
}

func stopProgress(executor models.Executor, readyMap map[string]bool, finishChan chan *models.ExecutedResult) bool {
	return executor.Stop(readyMap, finishChan)
}

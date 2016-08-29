package scheduler

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/stage"
)

type schedulerHand func(interface{}, map[string]bool, chan *models.ExecutedResult)

// StartStage start point on scheduler
func StartPoint(executorList []interface{}, startMark string) []*models.ExecutedResult {
	return execute(executorList, startMark, stage.Start)
}

// StopPoint stop point on scheduler
func StopPoint(executorList []interface{}, startMark string) []*models.ExecutedResult {
	return execute(executorList, startMark, stage.Stop)
}

func execute(executorList []interface{}, startMark string, handler schedulerHand) []*models.ExecutedResult {
	count := len(executorList)
	readyMap := map[string]bool{startMark: true}
	finishChan := make(chan *models.ExecutedResult, count)
	resultList := make([]*models.ExecutedResult, 0, count)
	running := true
	for running {
		for _, executor := range executorList {
			go handler(executor, readyMap, finishChan)
		}
		result := <-finishChan
		resultList = append(resultList, result)
		if result.Status != models.ResultSuccess {
			running = false
		} else {
			resultLen := len(resultList)
			running = resultLen != count
			log.Println(fmt.Sprintf("scheduler StartStage name = %v and finish num = %d", result.Name, resultLen))
		}
	}
	return resultList
}

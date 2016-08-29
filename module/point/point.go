package point

import (
	"github.com/containerops/vessel/module/stage"
	"github.com/containerops/vessel/models"
)

func Start(info interface{}, readyMap map[string]bool, finishChan chan *models.ExecutedResult) bool {
	pointVsn := info.(*models.PointVersion)
	for _, condition := range pointVsn.Conditions {
		if isReady, _ := readyMap[condition]; !isReady {
			return false
		}
	}
	if pointVsn.MateDate.Type == models.EndPoint{
		finishChan <- stage.FillSchedulingResult(models.EndPoint,models.ResultSuccess, "")
	}else{
		go stage.Start(pointVsn.StageVersion,finishChan)
	}
	return true
}

func Stop(info interface{}, readyMap map[string]bool, finishChan chan *models.ExecutedResult) bool {
	pointVsn := info.(*models.PointVersion)
	go stage.Stop(pointVsn.StageVersion,finishChan)
	return true
}
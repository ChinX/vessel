package stage

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
	kubeclt "github.com/containerops/vessel/module/kubernetes"
	"github.com/containerops/vessel/module/point"
)

// Start stage
func Start(info interface{}, readyMap map[string]bool, finishChan chan *models.ExecutedResult) {
	stageVsn := info.(*models.StageVersion)
	metaData := stageVsn.MetaData
	if stageVsn.State != "" {
		return
	}

	meet, ended := point.CheckCondition(stageVsn.PointVersion, readyMap)
	if ended {
		finishChan <- FillSchedulingResult(models.EndPointMark, models.ResultSuccess, "")
		return
	}
	if !meet {
		return
	}
	readyMap[metaData.Name] = true
	//TODO:Save stageVersion
	stageVsn.State = models.StateReady

	//res := kubeclt.CreateStage(metaData)
	//if res.Result != models.ResultSuccess {
	//	finishChan <- FillSchedulingResult(metaData.Name, res.Result, res.Detail)
	//	return
	//}

	//TODO:Update stageVersion
	stageVsn.State = models.StateRunning
	finishChan <- FillSchedulingResult(stageVsn.MetaData.Name, models.ResultSuccess, "")
}

// Stop stage
func Stop(info interface{}, readyMap map[string]bool, finishChan chan *models.ExecutedResult) {
	stageVsn := info.(*models.StageVersion)
	metaData := stageVsn.MetaData
	if stageVsn.State != models.StateReady || stageVsn.State != models.StateRunning {
		return
	}
	readyMap[metaData.Name] = true

	res := kubeclt.DeleteStage(stageVsn.MetaData)
	//TODO:Update stageVersion
	stageVsn.State = models.StateDeleted
	finishChan <- FillSchedulingResult(stageVsn.MetaData.Name, res.Result, res.Detail)
}

func FillSchedulingResult(stageName, result string, detail string) *models.ExecutedResult {
	log.Println(fmt.Sprintf("Stage name = %v result is %v, detail is %v", stageName, result, detail))
	return &models.ExecutedResult{
		Name:   stageName,
		Status: result,
		Detail: detail,
	}
}

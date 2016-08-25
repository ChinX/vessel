package stage

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
	kubeclt "github.com/containerops/vessel/module/kubernetes"
	"github.com/containerops/vessel/utils"
)

// StartStage start stage workflow
func StartStage(stage *models.Stage, finishChan chan *models.ExecutedResult) {
	if stage.Name == models.EndPointMark {
		finishChan <- fillSchedulingResult(stage, models.ResultSuccess, "")
		return
	}
	//TODO:Save stageVersion
	//stageVersion :=
	//stageVersion.Status = models.StateReady

	res := kubeclt.CreateStage(stage)
	if res.Result != models.ResultSuccess {
		finishChan <- fillSchedulingResult(stage, res.Result, res.Detail)
		return
	}

	//TODO:Update stageVersion
	//stageVersion.Status = models.StateRunning

	finishChan <- fillSchedulingResult(stage, models.ResultSuccess, "")
}

// StopStage stop stage workflow
func StopStage(stage *models.Stage, finishChan chan *models.ExecutedResult) {
	res := kubeclt.DeleteStage(stage)

	//TODO:Update stageVersion
	//stageVersion.Status = models.StateDeleted
	finishChan <- fillSchedulingResult(stage, res.Result, res.Detail)
}

func fillSchedulingResult(stage *models.Stage, result string, detail string) *models.ExecutedResult {
	log.Println(fmt.Sprintf("Stage name = %v result is %v, detail is %v", stage.Name, result, detail))
	stageName := ""
	namespace := ""
	if stage != nil {
		stageName = stage.Name
		namespace = stage.Namespace
	}
	return &models.ExecutedResult{
		Name:   stageName,
		Status: result,
		Result: &models.StageResult{
			ID:        utils.UUID(),
			Namespace: namespace,
			Name:      stageName,
			Result:    result,
			Detail:    detail,
		},
	}
}

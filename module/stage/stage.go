package stage

import (
	"fmt"
	"log"

	"github.com/containerops/vessel/models"
	kubeclt "github.com/containerops/vessel/module/kubernetes"
)

// Start stage
func Start(stageVsn *models.StageVersion, finishChan chan *models.ExecutedResult) {
	//TODO:Save stageVersion
	//stageVersion :=
	//stageVersion.Status = models.StateReady

	res := kubeclt.CreateStage(stageVsn.MateDate)
	if res.Result != models.ResultSuccess {
		finishChan <- FillSchedulingResult(stageVsn.MateDate.Name, res.Result, res.Detail)
		return
	}

	//TODO:Update stageVersion
	//stageVersion.Status = models.StateRunning

	finishChan <- FillSchedulingResult(stageVsn.MateDate.Name, models.ResultSuccess, "")
}

// Stop stage
func Stop(stageVsn *models.StageVersion, finishChan chan *models.ExecutedResult) {
	res := kubeclt.DeleteStage(stageVsn.MateDate)

	//TODO:Update stageVersion
	//stageVersion.Status = models.StateDeleted
	finishChan <- FillSchedulingResult(stageVsn.MateDate.Name, res.Result, res.Detail)
}

func FillSchedulingResult(stageName, result string, detail string) *models.ExecutedResult {
	log.Println(fmt.Sprintf("Stage name = %v result is %v, detail is %v", stageName, result, detail))
	return &models.ExecutedResult{
		Name:   stageName,
		Status: result,
		Detail: detail,
	}
}
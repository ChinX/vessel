package pipeline

import (
	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/dependence"
	"github.com/containerops/vessel/module/scheduler"
	"encoding/json"
)

// CreatePipeline new pipeline with PipelineTemplate
func CreatePipeline(pipelineTemplate *models.PipelineTemplate) []byte {
	log.Println("Create pipeline")
	pipeline := pipelineTemplate.MetaData

	// TODO:check pipeline already exist
	//if namespace name pipeline in db {
	//	detail := fmt.Sprintf("Pipeline = %v in namespane = %v is already exist", pipeline.Name, pipeline.Namespace)
	//	bytes, _ := formatOutputBytes(pipelineTemplate, pipeline, nil, detail)
	//	return bytes
	//}

	// TODO:check stage already exist
	//for _, stage := range pipeline.Stages {
	//if namespace name stage in db {
	//	detail := fmt.Sprintf("Stage = %v in namespane = %v is already exist", stage.Name, stage.Namespace)
	//	bytes, _ := formatOutputBytes(pipelineTemplate, pipeline, nil, detail)
	//	return bytes
	//}
	//}

	if err := dependence.CheckPipeline(pipeline); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}

	// TODO:save all pipeline
	log.Printf("Create pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Create job is done")
	bytes, _ := outputResult(pipeline, 0, nil, "")
	return bytes
}

func StartPipeline(pID uint64) []byte {
	log.Println("Start pipeline")
	// TODO:Get pipeline form db
	pipeline := &models.Pipeline{
		ID: pID,
	}

	executorMap, err := dependence.ParsePipeline(pipeline)
	if err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}

	pipelineVersion := &models.PipelineVersion{
		PID:pipeline.ID,
		Status:models.StateReady,
	}
	// TODO:save pipeline version

	schedulingRes := scheduler.Start(executorMap, models.StartPointMark)
	bytes, success := outputResult(pipeline, pipelineVersion.ID, schedulingRes, "")
	if success {
		pipelineVersion.Status = models.StateRunning
		// TODO:update pipeline version status
	} else {
		//rollback by pipeline failed
		go removePipeline(executorMap, pipelineVersion, "run pipeline error")
	}
	log.Printf("Start pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Start job is done")
	return bytes
}

func StopPipeline(pID uint64, pvID uint64) []byte {
	log.Println("Stop pipeline")
	// TODO: Get pipeline form db
	pipeline := &models.Pipeline{
		ID: pID,
	}

	executorMap, err := dependence.ParsePipeline(pipeline)
	if err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}

	// TODO: Get pipeline version form db
	pipelineVersion := &models.PipelineVersion{
		ID : pvID,
	}

	schedulingRes := removePipeline(executorMap, pipelineVersion, "")
	bytes, _ := outputResult(pipeline, pipelineVersion.ID, schedulingRes, "")
	log.Printf("Delete pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Delete job is done")
	return bytes
}

func DeletePipeline(pID uint64) []byte {
	log.Println("Delete pipeline")
	// TODO: Get pipeline form db
	pipeline := &models.Pipeline{
		ID: pID,
	}

	// TODO: Get pipeline version list form db with pID when is not delete
	// if len(list) != 0{
		//return has running can not delete
	//}

	// TODO:delete pipeline
	log.Printf("Delete pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Delete job is done")
	bytes, _ := outputResult(pipeline, 0, nil, "")
	return bytes
}

func RenewPipeline(pID uint64, pipelineTemplate *models.PipelineTemplate) []byte {
	log.Println("Renew pipeline")
	pipeline := pipelineTemplate.MetaData

	// TODO:check pipeline already exist
	//if namespace name pipeline not in db {
	//	detail := fmt.Sprintf("Pipeline = %v in namespane = %v is not already exist", pipeline.Name, pipeline.Namespace)
	//	bytes, _ := formatOutputBytes(pipelineTemplate, pipeline, nil, detail)
	//	return bytes
	//}

	// TODO:check stage already exist
	//for _, stage := range pipeline.Stages {
	//if namespace name stage not in db {
	//	detail := fmt.Sprintf("Stage = %v in namespane = %v is not already exist", stage.Name, stage.Namespace)
	//	bytes, _ := formatOutputBytes(pipelineTemplate, pipeline, nil, detail)
	//	return bytes
	//}
	//}

	if err := dependence.CheckPipeline(pipeline); err != nil {
		bytes, _ := outputResult(pipeline, 0, nil, err.Error())
		return bytes
	}

	// TODO:update all pipeline with pID
	log.Printf("Renew pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Renew job is done")
	bytes, _ := outputResult(pipeline, 0, nil, "")
	return bytes
	return nil
}

func GetPipeline(pID uint64) []byte {
	log.Println("Renew pipeline")
	// TODO:Get pipeline
	return nil
}

func removePipeline(executorMap map[string]models.Executor, pipelineVersion *models.PipelineVersion, detail string) []*models.ExecutedResult {
	schedulingRes := scheduler.Stop(executorMap, models.StartPointMark)
	pipelineVersion.Status = models.StateDeleted
	pipelineVersion.Detail = detail
	// TODO: delete pipeline version

	return schedulingRes
}

func outputResult(pipeline *models.Pipeline, pvID uint64, schedulingRes []*models.ExecutedResult, detail string) ([]byte, bool) {
	log.Println("Pipeline result :", schedulingRes)
	status := models.ResultFailed
	if detail == "" {
		status = models.ResultSuccess
		if schedulingRes != nil{
			for _, result := range schedulingRes {
				if status != result.Status {
					status = result.Status
					detail = result.Detail
					break
				}
			}
		}
	}
	output := &models.PipelineResult{
		PID:pipeline.ID,
		Name:pipeline.Name,
		Namespace:pipeline.Namespace,
		Status:status,
		Detail:detail,
	}
	if pvID != 0 {
		output.PvID = pvID
	}
	bytes, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Pipeline result is %v", string(bytes))
	return bytes, status == models.ResultSuccess
}

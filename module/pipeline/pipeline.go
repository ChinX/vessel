package pipeline

import (
	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/dependence"
	"github.com/containerops/vessel/module/scheduler"
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
		bytes, _ := formatOutputBytes()
		return bytes
	}

	// TODO:save pipeline
	log.Printf("Create pipeline name = %v in namespace '%v' is over", pipeline.Namespace, pipeline.Name)
	log.Print("Create job is done")
	return nil
}

func StartPipeline(pID uint64) []byte {
	log.Println("Start pipeline")
	// Get pipeline form db
	pipeline := &models.Pipeline{}
	pointMap, err := dependence.ParsePipeline(pipeline)
	if err != nil {
		bytes, _ := formatOutputBytes()
		return bytes
	}
	schedulingRes := scheduler.Start(pointMap, models.StartPointMark)
	log.Println(schedulingRes)
	return nil
}

func StopPipeline(pID uint64, pvID uint64) []byte {
	log.Println("Stop pipeline")
	// Get pipeline form db
	pipeline := &models.Pipeline{}
	pointMap, err := dependence.ParsePipeline(pipeline)
	if err != nil {
		bytes, _ := formatOutputBytes()
		return bytes
	}
	schedulingRes := scheduler.Stop(pointMap, models.StartPointMark)
	log.Println(schedulingRes)
	return nil
}

func RemovePipeline(pID uint64) []byte {
	return nil
}

func RenewPipeline(pID uint64, pipelineTemplate *models.PipelineTemplate) []byte {
	return nil
}

func GetPipeline(pID uint64) []byte {
	return nil
}

func formatOutputBytes() ([]byte, bool) {
	return nil, false
}

package dependence

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"encoding/json"
	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/utils"
	"github.com/containerops/vessel/utils/timer"
)

// ParsePipelineVersion parse executor map
func ParsePipelineVersion(pipelineVsn *models.PipelineVersion) []interface{} {
	pipe := pipelineVsn.MetaData
	hourglass := timer.InitHourglass(time.Duration(pipe.Timeout) * time.Second)
	executorList := make([]interface{}, 0, len(pipe.Stages)+1)

	//parse point
	pointVsnMap := make(map[string]*models.PointVersion, 0)
	var endPoint *models.PointVersion
	for _, point := range pipe.Points {
		triggers := utils.JSONStrToSlice(point.Triggers)
		pointVsn := &models.PointVersion{
			PvID:       pipelineVsn.ID,
			PointID:    point.ID,
			Conditions: utils.JSONStrToSlice(point.Conditions),
			MetaData:   point,
			Kind:       point.Type,
		}
		if point.Type == models.EndPoint {
			endPoint = pointVsn
			continue
		}
		for _, trigger := range triggers {
			pointVsnMap[trigger] = pointVsn
		}
	}

	//parse stage
	for _, stage := range pipe.Stages {
		stageVsn := &models.StageVersion{
			PvID:     pipelineVsn.ID,
			SID:      stage.ID,
			MetaData: stage,
		}
		stage.Hourglass = hourglass
		stage.PipelineName = pipe.Name
		pointVsn, ok := pointVsnMap[stage.Name]
		if !ok {
			pointVsn = &models.PointVersion{
				Kind:       models.TemporaryPoint,
				Conditions: utils.JSONStrToSlice(stage.Dependencies),
			}
			pointVsnMap[stage.Name] = pointVsn
		}
		stageVsn.PointVersion = pointVsn
		executorList = append(executorList, stageVsn)
	}

	return append(executorList, &models.StageVersion{
		PointVersion: endPoint,
	})
}

// CheckDependence check pipeline dependence
func CheckDependence(pipeline *models.Pipeline) error {
	conditionMap, err := checkPoints(pipeline.Points)
	if err != nil {
		return err
	}
	return checkStages(pipeline.Stages, conditionMap)
}

func checkPoints(points []*models.Point) (map[string][]string, error) {
	conditionMap := make(map[string][]string, 0)
	startPointCount := 0
	endPointCount := 0
	for _, pointInfo := range points {
		triggers := utils.JSONStrToSlice(pointInfo.Triggers)
		conditions := utils.JSONStrToSlice(pointInfo.Conditions)

		if pointInfo.Type == models.StartPoint {
			if conditions[0] != "" {
				return nil, fmt.Errorf("%v point condition must be empty", pointInfo.Type)
			}
			startPointCount++
		} else if conditions[0] == "" {
			return nil, fmt.Errorf("%v point condition must be not empty", pointInfo.Type)
		}

		if pointInfo.Type == models.EndPoint {
			if triggers[0] != "" {
				return nil, fmt.Errorf("%v point trigger must be empty", pointInfo.Type)
			}
			endPointCount++
		} else if triggers[0] == "" {
			return nil, fmt.Errorf("%v point trigger must be not empty", pointInfo.Type)
		}

		for _, trigger := range triggers {
			if _, ok := conditionMap[trigger]; ok {
				return nil, fmt.Errorf("Point trigger: %v is already exist", trigger)
			}
			conditionMap[trigger] = conditions
		}
	}
	if startPointCount < 1 {
		return nil, errors.New("Start point count must be greater than 0")
	}
	if endPointCount != 1 {
		return nil, errors.New("End point count must be 1")
	}
	return conditionMap, nil
}

func checkStages(stages []*models.Stage, conditionMap map[string][]string) error {
	stageMap := make(map[string]*models.Stage, 0)
	stageListMap := make(map[string][]string, 0)
	for _, stage := range stages {
		if stage.Name == "" {
			return errors.New("Stage has an empty name")
		}
		if _, ok := stageMap[stage.Name]; ok {
			return fmt.Errorf("Stage name: %v already exist", stage.Name)
		}
		bytes, _ := json.Marshal(stage)
		log.Println(string(bytes))
		dependencies := strings.Split(stage.Dependencies, ",")
		if conditions, ok := conditionMap[stage.Name]; ok {
			dependencies = conditions
			delete(conditionMap, stage.Name)
		} else {
			if _, ok := stageMap[stage.Name]; ok {
				return fmt.Errorf("Stage name: %v already exist", stage.Name)
			}
			if dependencies[0] == "" {
				return fmt.Errorf("No point stage: '%v' dependencies must be not empty", stage.Name)
			}
		}
		stageMap[stage.Name] = stage
		for _, dependence := range dependencies {
			stageList, ok := stageListMap[dependence]
			if !ok {
				stageList = make([]string, 0, 10)
			}
			stageList = append(stageList, stage.Name)
			stageListMap[dependence] = stageList
		}
	}
	return checkDependenceValidity(stageListMap, stageMap)
}

func checkDependenceValidity(stageListMap map[string][]string, stageMap map[string]*models.Stage) error {
	if len(stageListMap[""]) == 0 {
		return errors.New("The first start stage list is empty")
	}

	//Check dependence stage name is exist
	for dependenceName := range stageListMap {
		if dependenceName == "" {
			continue
		}
		_, ok := stageMap[dependenceName]
		if !ok {
			return fmt.Errorf("Dependence stage name: %v is not exist", dependenceName)
		}
	}

	//Check dependence directed acyclic graph
	return checkEndlessChain(stageListMap, make([]string, 0, 10), "")

}

func checkEndlessChain(stageListMap map[string][]string, chain []string, checkName string) error {
	if checkName != "" {
		for _, chainItem := range chain {
			if chainItem == checkName {
				return fmt.Errorf("Dependence chain [%v,%v] is endless chain", strings.Join(chain, ","), checkName)
			}
		}
	}
	stageList, ok := stageListMap[checkName]
	if ok {
		for _, nextStage := range stageList {
			chain = append(chain, checkName)
			err := checkEndlessChain(stageListMap, chain, nextStage)
			if err != nil {
				return err
			}
			chain = chain[0 : len(chain)-1]
		}
	}
	return nil
}

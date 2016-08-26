package dependence

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/point"
	"github.com/containerops/vessel/utils"
	"github.com/containerops/vessel/utils/timer"
	"k8s.io/kubernetes/pkg/util/json"
)

// CheckPipeline check point dependence
func CheckPipeline(pipeline *models.Pipeline) error {
	if err := checkUserPoint(pipeline.Points); err != nil {
		return err
	}

	executorMap, err := ParsePipeline(pipeline)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(executorMap)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(bytes))
	}
	return checkValidity(executorMap, models.StartPointMark)
}

// ParsePipeline parse point map from pipeline
func ParsePipeline(pipeline *models.Pipeline) (map[string]models.Executor, error) {
	hourglass := timer.InitHourglass(time.Duration(pipeline.Timeout) * time.Second)
	//parse user point
	pointMap, err := parsePoints(pipeline.Points)
	if err != nil {
		return nil, err
	}

	//endPoint
	execPoint := pointMap[models.EndPointMark]
	executorMap := map[string]models.Executor{models.EndPointMark: execPoint}
	execPoint.Info = &models.Stage{Name: models.EndPointMark}
	delete(pointMap, models.EndPointMark)

	//parse point from stage
	for _, stage := range pipeline.Stages {
		if stage.Name == "" {
			return nil, errors.New("Stage name must be not empty")
		}
		if execPoint, ok := pointMap[stage.Name]; ok {
			execPoint.Info = stage
			executorMap[stage.Name] = execPoint
			delete(pointMap, stage.Name)
		} else {
			if _, ok := executorMap[stage.Name]; ok {
				return nil, fmt.Errorf("Stage name: %v already exist", stage.Name)
			}
			if stage.Dependencies == "" {
				return nil, fmt.Errorf("No point stage: '%v' dependencies must be not empty", stage.Name)
			}
			executorMap[stage.Name] = &point.Point{
				Info: stage,
				From: utils.JSONStrToSlice(stage.Dependencies),
			}
		}
		stage.Hourglass = hourglass
		stage.Namespace = pipeline.Namespace
		stage.PipelineName = pipeline.Name
	}
	if len(pointMap) > 0 {
		byteArr, _ := json.Marshal(pointMap)
		return nil, fmt.Errorf("Pipeline points has unused point %v", string(byteArr))
	}
	return executorMap, nil
}

func parsePoints(points []*models.Point) (map[string]*point.Point, error) {
	pointMap := make(map[string]*point.Point, 0)
	for _, pointInfo := range points {
		triggers := utils.JSONStrToSlice(pointInfo.Triggers)
		for _, trigger := range triggers {
			if trigger == "" {
				trigger = models.EndPointMark
			}
			if _, ok := pointMap[trigger]; ok {
				return nil, fmt.Errorf("Point trigger :%v is already exist", trigger)
			}
			conditions := utils.JSONStrToSlice(pointInfo.Conditions)
			if conditions[0] == "" {
				conditions[0] = models.StartPointMark
			}
			pointMap[trigger] = &point.Point{From: conditions}
		}
	}
	return pointMap, nil
}

func checkUserPoint(points []*models.Point) error {
	starPointCount := 0
	endPointCount := 0
	for _, pointInfo := range points {
		triggers := utils.JSONStrToSlice(pointInfo.Triggers)
		conditions := utils.JSONStrToSlice(pointInfo.Conditions)
		switch pointInfo.Type {
		case models.StarPoint:
			if conditions[0] != "" {
				return fmt.Errorf("%v point condition must be empty", pointInfo.Type)
			}
			if triggers[0] == "" {
				return fmt.Errorf("%v point trigger must be not empty", pointInfo.Type)
			}
			starPointCount++
		case models.CheckPoint:
			if conditions[0] == "" {
				return fmt.Errorf("%v point condition must be not empty", pointInfo.Type)
			}
			if triggers[0] == "" {
				return fmt.Errorf("%v point trigger must be not empty", pointInfo.Type)
			}
		case models.EndPoint:
			if conditions[0] == "" {
				return fmt.Errorf("%v point condition must be not empty", pointInfo.Type)
			}
			if triggers[0] != "" {
				return fmt.Errorf("%v point trigger must be empty", pointInfo.Type)
			}
			endPointCount++
		}
	}
	if starPointCount < 1 {
		return errors.New("Start point count must be greater than 0")
	}
	if endPointCount != 1 {
		return errors.New("End point count must be 1")
	}
	return nil
}

func checkValidity(executorMap map[string]models.Executor, startName string) error {
	pointListMap := make(map[string][]string, 0)
	for name, execPoint := range executorMap {
		pointFrom := execPoint.GetFrom()
		for _, from := range pointFrom {
			//Check stage name is exist
			if from == models.StartPointMark {
				continue
			}
			if _, ok := executorMap[from]; !ok {
				return fmt.Errorf("Stage name: %v is not exist", from)
			}

			pointList, ok := pointListMap[from]
			if !ok {
				pointList = make([]string, 0, 10)
			}
			pointList = append(pointList, name)
			pointListMap[from] = pointList

		}
	}

	//Check dependence directed acyclic graph
	return checkEndlessChain(pointListMap, make([]string, 0, 10), startName)

}

func checkEndlessChain(pointListMap map[string][]string, chain []string, checkName string) error {
	if checkName != "" {
		for _, chainItem := range chain {
			if chainItem == checkName {
				return fmt.Errorf("Dependence chain [%v,%v] is endless chain", strings.Join(chain, ","), checkName)
			}
		}
	}
	stageList, ok := pointListMap[checkName]
	if ok {
		for _, nextStage := range stageList {
			chain = append(chain, checkName)
			err := checkEndlessChain(pointListMap, chain, nextStage)
			if err != nil {
				return err
			}
			chain = chain[0 : len(chain)-1]
		}
	}
	return nil
}

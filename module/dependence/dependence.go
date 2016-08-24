package dependence

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/utils"
)

// ParsePipeline parse point map from pipeline
func ParsePipeline(pipeline *models.Pipeline) (map[string]*models.Point, error) {
	endPoint, err := checkPoint(pipeline.Points)
	if err != nil {
		return nil, err
	}
	pointMap := make(map[string]*models.Point, 0)
	pointListMap := make(map[string][]string, 0)
	for _, stage := range pipeline.Stages {
		if stage.Name == "" {
			return nil, errors.New("Stage name must be not empty")
		}
		_, ok := pointMap[stage.Name]
		if ok {
			return nil, fmt.Errorf("Stage name: %v already exist", stage.Name)
		}

		point, err := getPoint(stage, pipeline.Points)
		if err != nil {
			return nil, err
		}
		pointMap[stage.Name] = point
		for _, from := range point.From {
			executorList, ok := pointListMap[from]
			if !ok {
				executorList = make([]string, 0, 10)
			}
			executorList = append(executorList, stage.Name)
			pointListMap[from] = executorList
		}
	}
	pointMap["$endPoint$"] = endPoint
	return pointMap, checkDependenceValidity(pointListMap, pointMap)
}

func checkPoint(points []*models.Point) (*models.Point, error) {
	var endPoint *models.Point
	starPointCount := 0
	endPointCount := 0
	for _, point := range points {
		point.To = utils.JsonStrToSlice(point.Triggers)
		point.From = utils.JsonStrToSlice(point.Conditions)
		switch point.Type {
		case models.StarPoint:
			if point.From[0] != "" {
				return nil, fmt.Errorf("%v point condition must be empty", point.Type)
			}
			if point.To[0] == "" {
				return nil, fmt.Errorf("%v point trigger must be not empty", point.Type)
			}
			starPointCount++
		case models.CheckPoint:
			if point.From[0] == "" {
				return nil, fmt.Errorf("%v point condition must be not empty", point.Type)
			}
			if point.To[0] == "" {
				return nil, fmt.Errorf("%v point trigger must be not empty", point.Type)
			}
		case models.EndPoint:
			if point.From[0] == "" {
				return nil, fmt.Errorf("%v point condition must be not empty", point.Type)
			}
			if point.To[0] != "" {
				return nil, fmt.Errorf("%v point trigger must be empty", point.Type)
			}
			endPointCount++
			endPoint = point
		}

	}
	if starPointCount < 1 {
		return nil, errors.New("Start point count must be greater than 0")
	}
	if endPointCount != 1 {
		return nil, errors.New("End point count must be 1")
	}
	return endPoint, nil
}

func getPoint(stage *models.Stage, points []*models.Point) (*models.Point, error) {
	dependencies := utils.JsonStrToSlice(stage.Dependencies)
	for _, point := range points {
		for _, trigger := range point.To {
			if stage.Name == trigger {
				if dependencies[0] != "" {
					return nil, errors.New("Point trigger stage must be no dependencies")
				}
				return point, nil
			}

		}
	}
	return &models.Point{Stage: stage, From: dependencies}, nil

}

func checkDependenceValidity(pointListMap map[string][]string, pointMap map[string]*models.Point) error {
	if len(pointListMap[""]) == 0 {
		return errors.New("The first start stage list is empty")
	}

	//Check dependence stage name is exist
	for dependenceName := range pointListMap {
		if dependenceName == "" {
			continue
		}
		_, ok := pointMap[dependenceName]
		if !ok {
			return fmt.Errorf("Dependence stage name: %v is not exist", dependenceName)
		}
	}

	//Check dependence directed acyclic graph
	return checkEndlessChain(pointListMap, make([]string, 0, 10), "")

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

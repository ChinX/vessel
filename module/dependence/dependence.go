package dependence

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/utils"
	"github.com/chinx/vessel/utils/timer"
	"time"
)

// CheckPipeline check point dependence
func CheckPipeline(pipeline *models.Pipeline) error {
	if err := checkUserPoint(pipeline.Points); err != nil {
		return err
	}

	pointMap, err := ParsePipeline(pipeline)
	if err != nil {
		return err
	}

	return checkValidity(pointMap)
}

// ParsePipeline parse point map from pipeline
func ParsePipeline(pipeline *models.Pipeline) (map[string]*models.Point, error) {
	pointMap := make(map[string]*models.Point, 0)
	hourglass := timer.InitHourglass(time.Duration(pipeline.Timeout)*time.Second)

	//parse user point
	for _, point := range pipeline.Points {
		if err := parsePoint(point, pointMap); err != nil {
			return nil, err
		}
		if point.Type == models.EndPoint {
			pointMap["$endPoint$"] = point
		}
	}

	//parse point from stage
	for _, stage := range pipeline.Stages {
		if stage.Name == "" {
			return nil, errors.New("Stage name must be not empty")
		}
		dependencies := utils.JSONStrToSlice(stage.Dependencies)
		if point, ok := pointMap[stage.Name]; ok {
			if point.Info != nil {
				return nil, fmt.Errorf("Stage name: %v already exist", stage.Name)
			}
			if dependencies != "" {
				return nil, fmt.Errorf("Point stage '%v' dependencies must be empty", stage.Name)
			}
			point.Info = stage
			point.From = utils.JSONStrToSlice(point.Conditions)
		} else {
			if dependencies == "" {
				return nil, fmt.Errorf("No point stage '%v' dependencies must be not empty", stage.Name)
			}
			pointMap[stage.Name] = &models.Point{
				Info: stage,
				From:  dependencies,
			}
		}
		stage.Hourglass = hourglass
	}
	return pointMap, nil
}

func parsePoint(point *models.Point, pointMap map[string]*models.Point) error {
	triggers := utils.JSONStrToSlice(point.Triggers)
	for _, trigger := range triggers {
		if _, ok := pointMap[trigger]; ok {
			return fmt.Errorf("Point trigger :%v is already exist", trigger)
		}
		pointMap[trigger] = point
	}
	return nil
}

func checkUserPoint(points []*models.Point) error {
	starPointCount := 0
	endPointCount := 0
	for _, point := range points {
		triggers := utils.JSONStrToSlice(point.Triggers)
		conditions := utils.JSONStrToSlice(point.Conditions)
		switch point.Type {
		case models.StarPoint:
			if conditions[0] != "" {
				return fmt.Errorf("%v point condition must be empty", point.Type)
			}
			if triggers[0] == "" {
				return fmt.Errorf("%v point trigger must be not empty", point.Type)
			}
			starPointCount++
		case models.CheckPoint:
			if conditions[0] == "" {
				return fmt.Errorf("%v point condition must be not empty", point.Type)
			}
			if triggers[0] == "" {
				return fmt.Errorf("%v point trigger must be not empty", point.Type)
			}
		case models.EndPoint:
			if conditions[0] == "" {
				return fmt.Errorf("%v point condition must be not empty", point.Type)
			}
			if triggers[0] != "" {
				return fmt.Errorf("%v point trigger must be empty", point.Type)
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

func checkValidity(pointMap map[string]*models.Point) error {
	pointListMap := make(map[string][]string, 0)
	for _, point := range pointMap {
		for name, from := range point.From {
			//Check stage name is exist
			if _, ok := pointMap[from]; !ok {
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

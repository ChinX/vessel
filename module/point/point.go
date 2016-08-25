package point

import (
	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/stage"
)

type Point struct {
	Info *models.Stage
	From []string
}

func (p Point) Start(readyMap map[string]bool, finishChan chan *models.ExecutedResult) bool {
	for _, from := range p.From {
		if isReady, _ := readyMap[from]; !isReady {
			return false
		}
	}
	go stage.StartStage(p.Info, finishChan)
	return true
}

func (p Point) Stop(readyMap map[string]bool, finishChan chan *models.ExecutedResult) bool {
	go stage.StopStage(p.Info, finishChan)
	return true
}

func (p Point) GetFrom() []string {
	return p.From
}

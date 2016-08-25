package point

import (
	"log"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/stage"
)

type Point struct {
	Info *models.Stage
	From []string
}

func (p Point) Start(readyMap map[string]bool, finishChan chan *models.ExecutedResult) bool {
	if p.Info == nil {
		log.Println("EndPoint triggers")
		return true
	}
	for _, from := range p.From {
		if isReady, _ := readyMap[from]; !isReady {
			return false
		}
	}
	go stage.StartStage(p.Info, finishChan)
	return true
}

func (p Point) Stop(readyMap map[string]bool, finishChan chan *models.ExecutedResult) bool {
	if p.Info == nil {
		log.Println("EndPoint triggers")
		return true
	}
	go stage.StopStage(p.Info, finishChan)
	return true
}

func (p Point) GetFrom() []string {
	return p.From
}

func (p Point) SetInfo(info *models.Stage) {
	p.Info = info
}

func (p Point) HasInfo() bool {
	return p.Info != nil
}

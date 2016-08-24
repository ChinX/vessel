package point

import "github.com/containerops/vessel/models"

func Start(info interface{}, finishChan chan *models.ExecutedResult) bool {
	return true
}

func Stop(info interface{}, finishChan chan *models.ExecutedResult) bool {
	return true
}
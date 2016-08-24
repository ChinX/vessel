package dependence

import (
	"encoding/json"
	"testing"

	"github.com/containerops/vessel/models"
)

func Test_ParsePipelineTemp(t *testing.T) {
	//str := jsonStr()
	//pipelineTemp := &models.PipelineTemplate{}
	//err := json.Unmarshal([]byte(str), pipelineTemp)
	//if err != nil {
	//	t.Log(err)
	//	return
	//}
	pipelineTemp := getTemp()
	stageMap, err := ParsePipeline(pipelineTemp.MetaData)
	if err != nil {
		t.Error(err)
		return
	}
	bytes, err := json.Marshal(stageMap)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(string(bytes))
	}
}

func getTemp() *models.PipelineTemplate {
	return &models.PipelineTemplate{
		Kind:"CCloud",
		APIVersion:"v1",
		MetaData: &models.Pipeline{
			Namespace:"guestbook",
			Name:"guestbook",
			Timeout:60,
			Stages:[]*models.Stage{
				&models.Stage{
					Name: "stageA",
					Type: models.StageContainer,
					Dependencies:"",
				},
				&models.Stage{
					Name: "stageB",
					Type: models.StageContainer,
					Dependencies:"",
				},
				&models.Stage{
					Name: "stageC",
					Type: models.StageContainer,
					Dependencies:"",
				},
				&models.Stage{
					Name: "stageD",
					Type: models.StageContainer,
					Dependencies:"",
				},
				&models.Stage{
					Name: "stageE",
					Type: models.StageContainer,
					Dependencies:"",
				},
				&models.Stage{
					Name: "stageF",
					Type: models.StageContainer,
					Dependencies:"stageA",
				},
				&models.Stage{
					Name: "stageG",
					Type: models.StageContainer,
					Dependencies:"stageA,stageB",
				},
			},
			Points:[]*models.Point{
				&models.Point{
					Type:models.StarPoint,
					Triggers:"stageA,stageB",
					Conditions:"",
				},
				&models.Point{
					Type:models.CheckPoint,
					Triggers:"stageC,stageD",
					Conditions:"stageA,stageB",
				},
				&models.Point{
					Type:models.CheckPoint,
					Triggers:"stageE",
					Conditions:"stageC,stageD",
				},
				&models.Point{
					Type:models.EndPoint,
					Conditions:"stageE",
				},
			},
		},
	}
}

func jsonStr() string {
	return `{
	    "kind": "CCloud",
	    "apiVersion": "v1",
	    "status": "",
	    "apiServerUrl": "",
	    "apiServerAuth": "",
	    "metadata": {
		"name": "guestbook",
		"namespace": "guestbook",
		"selfLink": "",
		"uid": "",
		"creationTimestamp": "",
		"deletionTimestamp": "",
		"labels": {
		    "app": "zenlin"
		},
		"annotations": {},
		"timeoutDuration": 60
	    },
	    "spec": [{
		"name": "redis-master",
		"replicas": 1,
		"dependence": "",
		"kind": "value",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "gcr.io/google_containers/redis:e2e",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "redis-slave",
		"replicas": 2,
		"dependence": "redis-master",
		"kind": "value",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "gcr.io/google_samples/gb-redisslave:v1",
		"port": 6379,
		"envName": "",
		"envValue": ""
	    }, {
		"name": "frontend",
		"replicas": 3,
		"dependence": "redis-slave",
		"kind": "value",
		"statusCheckLink": "/health",
		"statusCheckInterval": 0,
		"statusCheckCount": 0,
		"image": "gcr.io/google_samples/gb-frontend:v3",
		"port": 80,
		"envName": "GET_HOSTS_FROM",
		"envValue": "dns"
	    }]
	}`
}

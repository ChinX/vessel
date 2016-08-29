package dependence

import (
	"encoding/json"
	"testing"

	"github.com/containerops/vessel/models"
)

func Test_CheckPipeline(t *testing.T) {
	str := jsonStr()
	pipelineTemp := &models.PipelineTemplate{}
	err := json.Unmarshal([]byte(str), pipelineTemp)
	if err != nil {
		t.Log(err)
		return
	}

	if err := CheckDependence(pipelineTemp.MetaData); err != nil {
		t.Error(err)
		return
	}
}

func jsonStr() string {
	return `{
    "kind": "CCloud",
    "apiVersion": "v1",
    "metadata": {
        "name": "guestbook",
        "namespace": "guestbook",
        "timeout": 60,
        "stages": [
            {
                "name": "redis-master",
                "replicas": 1,
                "dependencies": "",
                "kind": "value",
                "statusCheckLink": "/health",
                "statusCheckInterval": 0,
                "statusCheckCount": 0,
                "image": "gcr.io/google_containers/redis:e2e",
                "port": 6379,
                "envName": "",
                "envValue": ""
            },
            {
                "name": "redis-slave",
                "replicas": 2,
                "dependencies": "",
                "kind": "value",
                "statusCheckLink": "/health",
                "statusCheckInterval": 0,
                "statusCheckCount": 0,
                "image": "gcr.io/google_samples/gb-redisslave:v1",
                "port": 6379,
                "envName": "",
                "envValue": ""
            },
            {
                "name": "frontend",
                "replicas": 3,
                "dependencies": "",
                "kind": "value",
                "statusCheckLink": "/health",
                "statusCheckInterval": 0,
                "statusCheckCount": 0,
                "image": "gcr.io/google_samples/gb-frontend:v3",
                "port": 80,
                "envName": "GET_HOSTS_FROM",
                "envValue": "dns"
            },
            {
                "name": "frontend1",
                "replicas": 3,
                "dependencies": "",
                "kind": "value",
                "statusCheckLink": "/health",
                "statusCheckInterval": 0,
                "statusCheckCount": 0,
                "image": "gcr.io/google_samples/gb-frontend:v3",
                "port": 80,
                "envName": "GET_HOSTS_FROM",
                "envValue": "dns"
            },
            {
                "name": "frontend2",
                "replicas": 3,
                "dependencies": "frontend1",
                "kind": "value",
                "statusCheckLink": "/health",
                "statusCheckInterval": 0,
                "statusCheckCount": 0,
                "image": "gcr.io/google_samples/gb-frontend:v3",
                "port": 80,
                "envName": "GET_HOSTS_FROM",
                "envValue": "dns"
            }
        ],
        "points": [
            {
                "type": "Start",
                "triggers": "redis-master",
                "conditions": ""
            },
            {
                "type": "Check",
                "triggers": "redis-slave,frontend",
                "conditions": "redis-master"
            },
            {
                "type": "Check",
                "triggers": "frontend1",
                "conditions": "redis-slave,frontend"
            },
            {
                "type": "End",
                "triggers": "",
                "conditions": "frontend2"
            }
        ]
    }
}`
}

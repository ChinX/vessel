package etcd

import (
	"fmt"
	"testing"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/coreos/etcd/client"
)

func TestCreateClient(t *testing.T) {
	if err := clientEtcd(); err != nil {
		t.Errorf("Create client err : %v", err.Error())
	}
}

func clientEtcd() (err error) {
	settingPoints := []map[string]string{
		map[string]string{
			"host": "127.0.0.1",
			"port": "4001",
		},
		map[string]string{
			"host": "localhost",
			"port": "4001",
		},
	}
	if etcd != nil {
		return nil
	}
	endPoints := []string{}
	for _, value := range settingPoints {
		endPoints = append(endPoints, fmt.Sprintf(models.EtcdConnectPath, value["host"], value["port"]))
	}

	cfg := client.Config{
		Endpoints:               endPoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	etcd, err = client.New(cfg)
	if err != nil {
		return err
	}
	return nil
}

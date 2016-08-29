package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var (
	etcd client.Client
)

const (
	// EtcdClientErr etcd client not connected error
	EtcdClientErr = "ETCD client is not connected"
)

func checkClient() error {
	if models.ETCD == nil {
		if err := models.InitEtcd(); err != nil {
			return err
		}
	}
	etcd = models.ETCD
	if etcd == nil {
		return errors.New(EtcdClientErr)
	}
	return nil
}

// Get get data from etcd
func Get(key string) (*client.Response, error) {
	if err := checkClient(); err != nil {
		return nil, err
	}
	return client.NewKeysAPI(etcd).Get(context.Background(), key, nil)
}

// Set save data to etcd
func Set(key string, value string, opts *client.SetOptions) error {
	if err := checkClient(); err != nil {
		return err
	}
	_, err := client.NewKeysAPI(etcd).Set(context.Background(), key, value, opts)
	return err
}

// GetValue get data from etcd as string
func GetValue(key string) (string, error) {
	resp, err := Get(key)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// SetValue save data to etcd as string
func SetValue(key string, value string) error {
	return Set(key, value, nil)
}

// GetJSON get data from etcd as JSON
func GetJSON(key string, v interface{}) error {
	jsonStr, err := GetValue(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonStr), v)
}

// SetJSON save data to etcd as JSON
func SetJSON(key string, value interface{}) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return SetValue(key, string(jsonBytes))
}

// SetDir create dir for etcd
func SetDir(key string) error {
	return Set(key, "", &client.SetOptions{Dir: true, PrevExist: client.PrevExist})
}

// GetDir get data from etcd dir
func GetDir(key string) (client.Nodes, error) {
	resp, err := Get(key)
	if err != nil {
		return nil, err
	}
	if !resp.Node.Dir {
		return nil, fmt.Errorf("%v is not dir in ETCD", key)
	}
	return resp.Node.Nodes, nil
}

// SetValueTTL set data TTL to etcd
func SetValueTTL(key string, value string, timeLife uint64) error {
	return Set(key, value, &client.SetOptions{TTL: time.Duration(timeLife)})
}

// SetDirTTL set dir TTl to etcd
func SetDirTTL(key string, timeLife uint64) error {
	return Set(key, "", &client.SetOptions{TTL: time.Duration(timeLife), Dir: true, PrevExist: client.PrevExist})
}

// Watch on etcd
func Watch(key string) (client.Watcher, error) {
	if err := checkClient(); err != nil {
		return nil, err
	}
	return client.NewKeysAPI(etcd).Watcher(key, &client.WatcherOptions{
		Recursive: true,
	}), nil
}

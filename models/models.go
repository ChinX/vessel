package models

import (
	"time"
	"fmt"
	"log"

	"github.com/containerops/vessel/setting"
	"github.com/coreos/etcd/client"
	"github.com/jinzhu/gorm"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// EtcdClient etcd client
	ETCD client.Client
	// K8sClient k8s client
	K8S *unversioned.Client

	db *gorm.DB

	//DataValidStatus the sql data is valid (general)
	DataValidStatus uint = 0
	//DataInValidStatus the sql data is invalid (delete)
	DataInValidStatus uint = 1
)

const (
	// EtcdConnectPath connect path for etcd
	EtcdConnectPath = "http://%s:%s"
	// K8sConnectPath connect path for k8s
	K8sConnectPath = "%v:%v"
	// DBConnectPath connect path for DB
	DBConnectPath = "%s:%s@%s(%s:%s)/%s?charset=%s&loc=%s&parseTime=%s"
)

// InitEtcd for etcd init
func InitEtcd() error {
	if ETCD == nil {
		var etcdEndPoints []string
		for _, value := range setting.RunTime.Etcd.Endpoints {
			etcdEndPoints = append(etcdEndPoints, fmt.Sprintf(EtcdConnectPath, value["host"], value["port"]))
		}

		cfg := client.Config{
			Endpoints: etcdEndPoints,
			Transport: client.DefaultTransport,
			// Set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
		var err error
		ETCD, err = client.New(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

// InitK8S for K8S init
func InitK8S() error {
	if K8S == nil {
		clientConfig := restclient.Config{}
		host := fmt.Sprintf(K8sConnectPath, setting.RunTime.K8s.Host, setting.RunTime.K8s.Port)
		clientConfig.Host = host
		// ClientConfig.Host = setting.RunTime.Database.Host
		client, err := unversioned.New(&clientConfig)
		if err != nil {
			log.Printf("New unversioned client err: %v!\n", err.Error())
			return err
		}
		K8S = client
	}
	return nil
}

// InitDatabase for mysql init
func InitDatabase() error {
	var err error
	if db == nil {
		dbArgs := fmt.Sprintf(DBConnectPath, setting.RunTime.Database.Username, setting.RunTime.Database.Password,
			setting.RunTime.Database.Protocol, setting.RunTime.Database.Host, setting.RunTime.Database.Port,
			setting.RunTime.Database.Schema, setting.RunTime.Database.Param["charset"],
			setting.RunTime.Database.Param["loc"], setting.RunTime.Database.Param["parseTime"])
		if db, err = gorm.Open("mysql", dbArgs); err != nil {
			return err
		}
		db.LogMode(true)
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		db.SingularTable(true)
		if err = Sync(); err != nil {
			return err
		}
	}
	return nil
}

//Sync database structs
func Sync() error {
	log.Println("Sync database structs ")
	db.AutoMigrate(&Pipeline{})
	db.AutoMigrate(&PipelineVersion{})
	db.AutoMigrate(&Point{})
	db.AutoMigrate(&Stage{})
	db.AutoMigrate(&StageVersion{})
	return nil
}

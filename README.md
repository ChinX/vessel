# Vessel 0001

test json
```

{
   "kind":"TestGroupServices",
   "apiVersion":"1",
   "metadata":{
      "name":"TestPipeline1",
      "namespace":"TestPipelineNS1",
      "selfLink":"CI REST API URI",
      "uid":"CI Key",
      "creationTimestamp":"backup",
      "deletionTimestamp":"backup",
      "timeoutDuration":7200,
      "labels":{"desc":"backup"},
      "annotations":{"ann":"backup"}
   },

   "spec":[
      {
         "name":"TestMasterServices",
         "dependence":"",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3,
         "replicsa":1,
         "image":"redis-master",
         "port":6379
      },
      {
         "name":"TestSlaveServices1",
         "dependence":"TestMasterServices",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3,
         "replicsa":2,
         "image":"redis-slave",
         "port":6379
      },
      {
         "name":"TestSlaveServices2",
         "dependence":"TestMasterServices",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3,
         "replicsa":1,
         "image":"mysql",
         "port":3306
      },
      {
         "name":"BaseServices0",
         "dependence":"",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3,
         "replicsa":1,
         "image":"restapi",
         "port":8080
      },
      {
         "name":"BaseServices1",
         "dependence":"",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3,
         "replicsa":1,
         "image":"xmlapi",
         "port":9090
      },
      {
         "name":"BaseServices2",
         "dependence":"",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3,
         "replicsa":1,
         "image":"yamlapi",
         "port":9090
      },
      {
         "name":"BaseServices3",
         "dependence":"BaseServices0,BaseServices1",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3,
         "replicsa":1,
         "image":"haproxy",
         "port":10000
      },
      {
         "name":"BaseServices4",
         "dependence":"BaseServices1,BaseServices2,BaseServices3",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3,
         "replicsa":1,
         "image":"oos",
         "port":20000
      },
      {
         "name":"TestServices4",
         "dependence":"BaseServices4",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3,
         "replicsa":1,
         "image":"unittest",
         "port":30000
      }
   ]
}

```

runtime.yaml

```

---
run:
  runMode: dev
  logPath: log/vessel
http:
  listenMode: http
  httpsCertFile: cert/containerops/containerops.crt
  httpsKeyFile: cert/containerops/containerops.key
  host: 0.0.0.0
  port: 8080
database:
  username: vessel
  password: vessel
  protocol: tcp
  host: localhost
  port: 3306
  schema: vessel
  param:
    charset: utf8
    parseTime: True
    loc: Local
etcd:
  endpoints:
    - host: localhost
      port: 2379
    - host: 127.0.0.1
      port: 2379
  username: etcd
  password: etcd

```

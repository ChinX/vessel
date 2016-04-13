# Vessel

test json
'''
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
      "labels":"backup",
      "annotations":"backup"
   },
   "spec":[
      {
         "name":"TestMasterServices",
         "dependence":"",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3
      },
      {
         "name":"TestSlaveServices1",
         "dependence":"TestMasterServices",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3
      },
      {
         "name":"TestSlaveServices2",
         "dependence":"TestMasterServices",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3
      },
      {
         "name":"BaseServices0",
         "dependence":"",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3
      },
      {
         "name":"BaseServices1",
         "dependence":"",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3
      },
      {
         "name":"BaseServices2",
         "dependence":"",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3
      },
      {
         "name":"BaseServices3",
         "dependence":"BaseServices0,BaseServices1",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3
      },
      {
         "name":"BaseServices4",
         "dependence":"BaseServices1,BaseServices2,BaseServices3",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3
      },
      {
         "name":"TestServices4",
         "dependence":"BaseServices4",
         "kind":"backup",
         "status_check_url":"",
         "status_check_interval":30,
         "status_check_count":3
      }
   ]
}

'''

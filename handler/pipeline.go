package handler

import (
	"net/http"

	"github.com/containerops/vessel/models"
	"gopkg.in/macaron.v1"
)

// POSTPipeline new a pipeline
func POSTPipeline(ctx *macaron.Context, reqData models.PipelineTemplate) (int, []byte) {
	return http.StatusOK, []byte("")
}

// DELETEPipeline delete the pipeline with pid
func DELETEPipelinePID(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// PUTPipelinePID update the pipeline with pid
func PUTPipelinePID(ctx *macaron.Context, reqData models.PipelineTemplate) (int, []byte) {
	return http.StatusOK, []byte("")
}

// GETPipelinePID get the pipeline with pid
func GETPipelinePID(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// POSTPipelinePID exec the pipeline with pid
func POSTPipelinePID(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// DELETEPipeline delete the pipeline with pid
func DELETEPipelinePIDPvID(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// GETPipelinePIDPvID get the pipeline result with pid and pvid
func GETPipelinePIDPvID(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// GETPipelinePIDPvIDLogs get system logs with pid and pvid
func GETPipelinePIDPvIDLogs(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// GETPipelineNamespaceName get pipeline list with pipeline name and namespace
func GETPipelineNamespaceName(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

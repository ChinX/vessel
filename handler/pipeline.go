package handler

import (
	"net/http"
	"strconv"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/pipeline"
	"gopkg.in/macaron.v1"
)

// POSTPipeline new a pipeline
func POSTPipeline(ctx *macaron.Context, reqData models.PipelineTemplate) (int, []byte) {
	return http.StatusOK, pipeline.CreatePipeline(&reqData)
}

// DELETEPipelinePID delete the pipeline with pid
func DELETEPipelinePID(ctx *macaron.Context) (int, []byte) {
	pID, err := strconv.ParseUint(ctx.Params("pid"), 10, 10)
	if err != nil {
		return requestErrBytes([]string{"pid"}, err)
	}
	return http.StatusOK, pipeline.RemovePipeline(pID)
}

// PUTPipelinePID update the pipeline with pid
func PUTPipelinePID(ctx *macaron.Context, reqData models.PipelineTemplate) (int, []byte) {
	pID, err := strconv.ParseUint(ctx.Params("pid"), 10, 10)
	if err != nil {
		return requestErrBytes([]string{"pid"}, err)
	}
	return http.StatusOK, pipeline.RenewPipeline(pID, &reqData)
}

// GETPipelinePID get the pipeline with pid
func GETPipelinePID(ctx *macaron.Context) (int, []byte) {
	pID, err := strconv.ParseUint(ctx.Params("pid"), 10, 10)
	if err != nil {
		return requestErrBytes([]string{"pid"}, err)
	}
	return http.StatusOK, pipeline.GetPipeline(pID)
}

// POSTPipelinePID exec the pipeline with pid
func POSTPipelinePID(ctx *macaron.Context) (int, []byte) {
	pID, err := strconv.ParseUint(ctx.Params("pid"), 10, 10)
	if err != nil {
		return requestErrBytes([]string{"pid"}, err)
	}
	return http.StatusOK, pipeline.StartPipeline(pID)
}

// DELETEPipelinePIDPvID stop the pipeline with pid and pvid
func DELETEPipelinePIDPvID(ctx *macaron.Context) (int, []byte) {
	pID, err := strconv.ParseUint(ctx.Params("pid"), 10, 10)
	if err != nil {
		return requestErrBytes([]string{"pid"}, err)
	}
	pvID, err := strconv.ParseUint(ctx.Params("pvid"), 10, 10)
	if err != nil {
		return requestErrBytes([]string{"pvid"}, err)
	}
	return http.StatusOK, pipeline.StopPipeline(pID, pvID)
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

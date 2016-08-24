package handler

import (
	"net/http"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/pipeline"
	"gopkg.in/macaron.v1"
)

// V1POSTPipeline handler for HTTP POST
func V1POSTPipeline(ctx *macaron.Context, reqData models.PipelineSpecTemplate) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTPipeline handler for HTTP PUT
func V1PUTPipeline(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

// V1GETPipeline handler for HTTP GET
func V1GETPipeline(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

// V1DELETEPipeline handler for HTTP DELETE
func V1DELETEPipeline(ctx *macaron.Context, reqData models.PipelineSpecTemplate) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1RunPipeline handler of run pipeline
func V1RunPipeline(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

// V1GETPipelineStatus handler of pipeline status for HTTP GET
func V1GETPipelineStatus(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

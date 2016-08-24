package router

import (
	"gopkg.in/macaron.v1"
)

// SetRouters set routers to macaron
func SetRouters(m *macaron.Macaron) {
	//m.Get("/test", handler.V1GETHealth)
	//m.Group("/vessel", func() {
	//	m.Group("/v1", func() {
	//		m.Group("/pipeline", func() {
	//
	//			m.Post("/", binding.Bind(models.PipelineSpecTemplate{}), handler.V1POSTPipeline)
	//			m.Put("/:pipeline", handler.V1PUTPipeline)
	//			m.Get("/:pipeline", handler.V1GETPipeline)
	//			m.Delete("/", binding.Bind(models.PipelineSpecTemplate{}), handler.V1DELETEPipeline)
	//
	//			m.Put("/:pipeline/run", handler.V1RunPipeline)
	//
	//			m.Get("/:pipeline/status", handler.V1GETPipelineStatus)
	//		})
	//	})
	//})
}

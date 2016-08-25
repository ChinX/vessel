package router

import (
	"github.com/containerops/vessel/handler"
	"github.com/containerops/vessel/models"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

// SetRouters set routers to macaron
func SetRouters(m *macaron.Macaron) {
	m.Group("/vessel", func() {
		m.Post("/", binding.Bind(models.PipelineTemplate{}), handler.POSTPipeline)
		m.Group("/:pid", func() {
			m.Post("/", handler.POSTPipelinePID)
			m.Put("/", handler.PUTPipelinePID)
			m.Get("/", handler.GETPipelinePID)
			m.Delete("/", handler.DELETEPipelinePID)
			m.Group("/:pvid", func() {
				m.Get("/", handler.GETPipelinePIDPvID)
				m.Delete("/", handler.DELETEPipelinePIDPvID)
				m.Get("/logs", handler.GETPipelinePIDPvIDLogs)
			})
		})
		m.Get("/:namespace/:name", handler.GETPipelineNamespaceName)
	})
}

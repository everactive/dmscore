package web

import (
	"context"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/everactive/dmscore/iot-management/web"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thejerf/suture/v4"
	"gorm.io/gorm"
	"time"
)

type HandlerService struct {
	manage    manage.Manage
	legacyWeb *web.Service
	engine    *gin.Engine
	db        *gorm.DB
}

func New(srv manage.Manage, db *gorm.DB) *suture.Supervisor {
	sup := suture.NewSimple("webhandlers")

	engine := gin.Default()

	engine.Use(gin.Logger())

	legacyWeb := web.NewService(srv, engine)
	hs := &HandlerService{srv, legacyWeb, engine, db}

	sup.Add(hs)
	sup.Add(legacyWeb)

	group := engine.Group("/v1")
	group.Use(web.AuthMiddleWare)
	group.POST("/:orgid/models/:model/required", hs.AddRequiredModelSnap)
	group.DELETE("/:orgid/models/:model/required", hs.DeleteRequiredModelSnap)
	group.GET("/:orgid/models/:model/required", hs.RequiredModelSnaps)

	return sup
}

func (h *HandlerService) Serve(ctx context.Context) error {
	intervalTicker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Errorf("We're done: %s", ctx.Err())
			return nil
		case <-intervalTicker.C:
			log.Infof("Still ticking")
		}
	}

	return nil
}

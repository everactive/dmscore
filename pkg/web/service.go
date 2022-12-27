package web

import (
	"context"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/everactive/dmscore/iot-management/web"
	log "github.com/sirupsen/logrus"
	"github.com/thejerf/suture/v4"
	"gorm.io/gorm"
	"time"
)

type HandlerService struct {
	db        *gorm.DB
}

func New(srv manage.Manage, db *gorm.DB) *suture.Supervisor {
	sup := suture.NewSimple("webhandlers")

	legacyWeb := web.NewService(srv)
        hs := &HandlerService{db}

        sup.Add(hs)
	sup.Add(legacyWeb)

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

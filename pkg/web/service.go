package web

import (
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/everactive/dmscore/iot-management/web"
	"github.com/thejerf/suture/v4"
)

func New(srv manage.Manage) *suture.Supervisor {
	sup := suture.NewSimple("webhandlers")

	legacyWeb := web.NewService(srv)

	sup.Add(legacyWeb)

	return sup
}
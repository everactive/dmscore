// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Management Service
 * Copyright 2019 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU Affero General Public License version 3, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY,
 * SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package web

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/everactive/dmscore/iot-management/config"
	"github.com/everactive/dmscore/iot-management/config/configkey"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// StoreSearchHandler fetches the available snaps from the store
func (wb Service) StoreSearchHandler(c *gin.Context) {
	w := c.Writer

	user, err := getUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		formatStandardResponse("UserAuth", "", c)
		return
	}
	w.Header().Set("Content-Type", JSONHeader)

	storeURL := viper.GetString(configkey.StoreURL)
	req, err := http.NewRequest("GET", storeURL+"snaps/search?q="+c.Param("snapName"), nil)
	if err != nil {
		log.Error(err)
		fmt.Fprint(w, "{}")
		return
	}
	req.Header.Add("X-Ubuntu-Series", "16")

	storeID := viper.GetString(fmt.Sprintf(config.ModelKeyTemplate, c.Param("model")))
	if storeID == "" {
		log.Errorf("unrecognized device model: %s", c.Param("model"))
		_, errInt := fmt.Fprint(w, "{}")
		if errInt != nil {
			log.Error(errInt)
		}
		return
	}

	req.Header.Add("X-Ubuntu-Store", storeID)

	client := &http.Client{}
	resp, err2 := client.Do(req)
	defer func() {
		errInt := resp.Body.Close()
		if errInt != nil {
			log.Error(errInt)
		}
	}()

	if err2 != nil {
		fmt.Fprint(w, "{}")
		return
	}

	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Fprint(w, "{}")
		return
	}

	fmt.Fprint(w, string(body))
}

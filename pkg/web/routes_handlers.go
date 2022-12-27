package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everactive/dmscore/api"
	"github.com/everactive/dmscore/iot-management/datastore"
	"github.com/everactive/dmscore/iot-management/service/manage"
	"github.com/everactive/dmscore/iot-management/web"
	"github.com/everactive/dmscore/pkg/messages"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

var ErrUserInvalidOrNotAuthorized = errors.New("user is invalid or not authorized")
var ErrDecodingBody = errors.New("error decoding body of request")

func getModelRequiredSnap(c *gin.Context) (*messages.ModelRequiredSnap, error) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, ErrDecodingBody
	}

	var modelSnap messages.ModelRequiredSnap
	err = json.Unmarshal(bodyBytes, &modelSnap)

	if err != nil {
		return nil, ErrDecodingBody
	}

	return &modelSnap, nil
}

func (h *HandlerService) AddRequiredModelSnap(c *gin.Context) {
	user, err := web.GetUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		response := api.StandardResponse{Code: "UserAuth", Message: fmt.Sprintf("AddRequiredModelSnap: %+v", err)}
		c.JSON(http.StatusUnauthorized, &response)
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response := api.StandardResponse{Code: "AddRequiredModelSnap"}
		c.JSON(http.StatusInternalServerError, &response)
		return
	}

	var modelSnap messages.ModelRequiredSnap
	err = json.Unmarshal(bodyBytes, &modelSnap)

	if err != nil {
		response := api.StandardResponse{Code: "AddRequiredModelSnap"}
		c.JSON(http.StatusInternalServerError, &response)
		return
	}

	requiredSnap, err := h.manage.AddModelRequiredSnap(c.Param("orgid"), user.Username, c.Param("model"), modelSnap.Snap, user.Role)

	if err != nil {
		if err == manage.NotAuthorizedErr {
			response := api.StandardResponse{Code: "UserAuth"}
			c.JSON(http.StatusUnauthorized, &response)
			return
		}

		response := api.StandardResponse{Code: "Error", Message: err.Error()}
		c.JSON(http.StatusInternalServerError, &response)
		return
	}

	c.JSON(http.StatusOK, &requiredSnap)
	return
}

func (h *HandlerService) DeleteRequiredModelSnap(c *gin.Context) {
	user, err := web.GetUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		response := api.StandardResponse{Code: "UserAuth", Message: ErrUserInvalidOrNotAuthorized.Error()}
		c.JSON(http.StatusUnauthorized, &response)
		return
	}

	requiredSnap, err := getModelRequiredSnap(c)
	if err != nil {
		if err == ErrDecodingBody {
			response := api.StandardResponse{Code: "Error", Message: err.Error()}
			c.JSON(http.StatusInternalServerError, &response)
			return
		} else {
			response := api.StandardResponse{Code: "UserAuth"}
			c.JSON(http.StatusUnauthorized, &response)
			return
		}
	}

	err = h.manage.DeleteModelRequiredSnap(c.Param("orgid"), user.Username, c.Param("model"), requiredSnap.Snap, user.Role)

	if err != nil {
		response := api.StandardResponse{Code: "Error", Message: err.Error()}
		c.JSON(http.StatusInternalServerError, &response)
		return
	}

	response := api.StandardResponse{}
	c.JSON(http.StatusOK, &response)
}

// RequiredModelSnaps gets the snaps currently required for a given model
func (h *HandlerService) RequiredModelSnaps(c *gin.Context) {
	user, err := web.GetUserFromContextAndCheckPermissions(c, datastore.Standard)
	if user == nil || err != nil {
		response := api.StandardResponse{Code: "UserAuth", Message: fmt.Sprintf("AddRequiredModelSnap: %+v", err)}
		c.JSON(http.StatusUnauthorized, &response)
		return
	}

	device, err := h.manage.GetModelRequiredSnaps(c.Param("orgid"), user.Username, c.Param("model"), user.Role)
	if err != nil {
		response := api.StandardResponse{Code: "Error", Message: err.Error()}
		c.JSON(http.StatusInternalServerError, &response)
		return
	}

	c.JSON(http.StatusOK, &device)
	return
}

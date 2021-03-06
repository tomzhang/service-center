//Copyright 2017 Huawei Technologies Co., Ltd
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
package v4

import (
	"encoding/json"
	"github.com/ServiceComb/service-center/pkg/rest"
	"github.com/ServiceComb/service-center/pkg/util"
	"github.com/ServiceComb/service-center/server/core"
	pb "github.com/ServiceComb/service-center/server/core/proto"
	scerr "github.com/ServiceComb/service-center/server/error"
	"github.com/ServiceComb/service-center/server/rest/controller"
	"io/ioutil"
	"net/http"
	"strings"
)

type MicroServiceService struct {
	//
}

func (this *MicroServiceService) URLPatterns() []rest.Route {
	return []rest.Route{
		{rest.HTTP_METHOD_GET, "/v4/:domain/registry/existence", this.GetExistence},
		{rest.HTTP_METHOD_GET, "/v4/:domain/registry/microservices", this.GetServices},
		{rest.HTTP_METHOD_GET, "/v4/:domain/registry/microservices/:serviceId", this.GetServiceOne},
		{rest.HTTP_METHOD_POST, "/v4/:domain/registry/microservices", this.Register},
		{rest.HTTP_METHOD_PUT, "/v4/:domain/registry/microservices/:serviceId/properties", this.Update},
		{rest.HTTP_METHOD_DELETE, "/v4/:domain/registry/microservices/:serviceId", this.Unregister},
		{rest.HTTP_METHOD_DELETE, "/v4/:domain/registry/microservices", this.UnregisterServices},
	}
}

func (this *MicroServiceService) Register(w http.ResponseWriter, r *http.Request) {
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Logger().Error("body err", err)
		controller.WriteError(w, scerr.ErrInvalidParams, err.Error())
		return
	}
	var request pb.CreateServiceRequest
	err = json.Unmarshal(message, &request)
	if err != nil {
		util.Logger().Error("Unmarshal error", err)
		controller.WriteError(w, scerr.ErrInvalidParams, err.Error())
		return
	}
	resp, err := core.ServiceAPI.Create(r.Context(), &request)
	respInternal := resp.Response
	resp.Response = nil
	controller.WriteResponse(w, respInternal, resp)
}

func (this *MicroServiceService) Update(w http.ResponseWriter, r *http.Request) {
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Logger().Error("body err", err)
		controller.WriteError(w, scerr.ErrInvalidParams, err.Error())
		return
	}
	request := &pb.UpdateServicePropsRequest{
		ServiceId: r.URL.Query().Get(":serviceId"),
	}
	err = json.Unmarshal(message, request)
	if err != nil {
		util.Logger().Error("Unmarshal error", err)
		controller.WriteError(w, scerr.ErrInvalidParams, err.Error())
		return
	}
	resp, err := core.ServiceAPI.UpdateProperties(r.Context(), request)
	controller.WriteResponse(w, resp.GetResponse(), nil)
}

func (this *MicroServiceService) Unregister(w http.ResponseWriter, r *http.Request) {
	force := r.URL.Query().Get("force")
	serviceId := r.URL.Query().Get(":serviceId")
	util.Logger().Warnf(nil, "Service %s unregists, force is %s.", serviceId, force)
	if force != "0" && force != "1" && strings.TrimSpace(force) != "" {
		controller.WriteError(w, scerr.ErrInvalidParams, "parameter force must be 1 or 0")
		return
	}
	request := &pb.DeleteServiceRequest{
		ServiceId: serviceId,
		Force:     force == "1",
	}
	resp, _ := core.ServiceAPI.Delete(r.Context(), request)
	controller.WriteResponse(w, resp.GetResponse(), nil)
}

func (this *MicroServiceService) GetServices(w http.ResponseWriter, r *http.Request) {
	request := &pb.GetServicesRequest{}
	util.Logger().Debugf("domain is %s", util.ParseDomain(r.Context()))
	resp, _ := core.ServiceAPI.GetServices(r.Context(), request)
	respInternal := resp.Response
	resp.Response = nil
	controller.WriteResponse(w, respInternal, resp)
}

func (this *MicroServiceService) GetExistence(w http.ResponseWriter, r *http.Request) {
	request := &pb.GetExistenceRequest{
		Type:        r.URL.Query().Get("type"),
		AppId:       r.URL.Query().Get("appId"),
		ServiceName: r.URL.Query().Get("serviceName"),
		Version:     r.URL.Query().Get("version"),
		ServiceId:   r.URL.Query().Get("serviceId"),
		SchemaId:    r.URL.Query().Get("schemaId"),
	}
	resp, _ := core.ServiceAPI.Exist(r.Context(), request)
	w.Header().Add("X-Schema-Summary", resp.Summary)
	respInternal := resp.Response
	resp.Response = nil
	resp.Summary = ""
	controller.WriteResponse(w, respInternal, resp)
}

func (this *MicroServiceService) GetServiceOne(w http.ResponseWriter, r *http.Request) {
	request := &pb.GetServiceRequest{
		ServiceId: r.URL.Query().Get(":serviceId"),
	}
	resp, _ := core.ServiceAPI.GetOne(r.Context(), request)
	respInternal := resp.Response
	resp.Response = nil
	controller.WriteResponse(w, respInternal, resp)
}

func (this *MicroServiceService) UnregisterServices(w http.ResponseWriter, r *http.Request) {
	request_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Logger().Error("body ,err", err)
		controller.WriteError(w, scerr.ErrInvalidParams, err.Error())
		return
	}

	request := &pb.DelServicesRequest{}

	err = json.Unmarshal(request_body, request)
	if err != nil {
		util.Logger().Error("unmarshal ,err ", err)
		controller.WriteError(w, scerr.ErrInvalidParams, err.Error())
		return
	}

	resp, err := core.ServiceAPI.DeleteServices(r.Context(), request)
	respInternal := resp.Response
	resp.Response = nil
	controller.WriteResponse(w, respInternal, resp)
}

/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package http

import (
	"context"
	"encoding/json"

	"github.com/gorilla/mux"

	"vitess.io/vitess/go/vt/vtadmin/errors"

	vtadminpb "vitess.io/vitess/go/vt/proto/vtadmin"
	vtctldatapb "vitess.io/vitess/go/vt/proto/vtctldata"
)

// CreateKeyspace implements the http wrapper for POST /keyspace/{cluster_id}.
func CreateKeyspace(ctx context.Context, r Request, api *API) *JSONResponse {
	vars := mux.Vars(r.Request)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var req vtctldatapb.CreateKeyspaceRequest
	if err := decoder.Decode(&req); err != nil {
		return NewJSONResponse(nil, &errors.BadRequest{
			Err: err,
		})
	}

	resp, err := api.server.CreateKeyspace(ctx, &vtadminpb.CreateKeyspaceRequest{
		ClusterId: vars["cluster_id"],
		Options:   &req,
	})
	return NewJSONResponse(resp, err)
}

// DeleteKeyspace implements the http wrapper for DELETE /keyspace/{cluster_id}/{name}[?recursive=].
func DeleteKeyspace(ctx context.Context, r Request, api *API) *JSONResponse {
	vars := mux.Vars(r.Request)
	recursive, err := r.ParseQueryParamAsBool("recursive", false)
	if err != nil {
		return NewJSONResponse(nil, err)
	}

	resp, err := api.server.DeleteKeyspace(ctx, &vtadminpb.DeleteKeyspaceRequest{
		ClusterId: vars["cluster_id"],
		Options: &vtctldatapb.DeleteKeyspaceRequest{
			Keyspace:  vars["name"],
			Recursive: recursive,
		},
	})
	return NewJSONResponse(resp, err)
}

// GetKeyspace implements the http wrapper for /keyspace/{cluster_id}/{name}.
func GetKeyspace(ctx context.Context, r Request, api *API) *JSONResponse {
	vars := mux.Vars(r.Request)
	keyspace, err := api.server.GetKeyspace(ctx, &vtadminpb.GetKeyspaceRequest{
		ClusterId: vars["cluster_id"],
		Keyspace:  vars["name"],
	})

	return NewJSONResponse(keyspace, err)
}

// GetKeyspaces implements the http wrapper for /keyspaces[?cluster=[&cluster=]].
func GetKeyspaces(ctx context.Context, r Request, api *API) *JSONResponse {
	keyspaces, err := api.server.GetKeyspaces(ctx, &vtadminpb.GetKeyspacesRequest{
		ClusterIds: r.URL.Query()["cluster"],
	})

	return NewJSONResponse(keyspaces, err)
}

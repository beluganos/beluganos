// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ribssrv

import (
	"fabricflow/ribs/api/ribsapi"
	"fabricflow/ribs/pkgs/ribsdbm"
	"fabricflow/util/gobgp/apiutil"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/pkg/packet/bgp"
	"golang.org/x/net/context"
)

//
// RibUpdate is rib update mesage.
//
type RibUpdate struct {
	Path *api.Path
	Rt   string
}

//
// NewRibUpdate returns new RibUpdate.
//
func NewRibUpdate(path *api.Path, rt string) *RibUpdate {
	return &RibUpdate{
		Path: path,
		Rt:   rt,
	}
}

//
// ToAPI returns RibUpdate of RibsAPI.
//
func (m *RibUpdate) ToAPI() (*ribsapi.RibUpdate, error) {
	return NewRibUpdateAPIFromGoBGPPath(m.Path, m.Rt)
}

//
// NewRibUpdateFromAPI returns RibUpdate.
//
func NewRibUpdateFromAPI(msg *ribsapi.RibUpdate) (*RibUpdate, error) {
	path, rt, err := NewGoBGPPathFromAPI(msg)
	if err != nil {
		return nil, err
	}

	return &RibUpdate{
		Path: path,
		Rt:   rt,
	}, nil
}

//
// NewRibUpdateAPIFromGoBGPPath returns RibUpdate(RibsApi)
//
func NewRibUpdateAPIFromGoBGPPath(path *api.Path, rt string) (*ribsapi.RibUpdate, error) {
	b, err := proto.Marshal(path)
	if err != nil {
		return nil, err
	}

	return &ribsapi.RibUpdate{
		Rt:   rt,
		Path: b,
	}, nil
}

//
// NewGoBGPPathFromAPI returns Path(GoBGP API)
//
func NewGoBGPPathFromAPI(msg *ribsapi.RibUpdate) (*api.Path, string, error) {
	path := &api.Path{}
	if err := proto.Unmarshal(msg.Path, path); err != nil {
		return nil, "", err
	}

	return path, msg.Rt, nil
}

func getExtendedCommunityRouteTarget(pattrs []bgp.PathAttributeInterface, rt string) (bgp.ExtendedCommunityInterface, bool) {
	exRT, ok := apiutil.GetNativeExtCommunityAttribute(bgp.EC_SUBTYPE_ROUTE_TARGET, pattrs)
	if !ok {
		return nil, false
	}

	if rt == ribsdbm.RTany {
		return exRT, true
	}

	return exRT, (exRT.String() == rt)
}

func listGoBGPPath(client api.GobgpApiClient, family *api.Family, f func(*api.Path) error) error {
	req := api.ListPathRequest{
		TableType: api.TableType_GLOBAL,
		Family:    family,
	}
	stream, err := client.ListPath(context.Background(), &req)
	if err != nil {
		return fmt.Errorf("ListPath error. %s", err)
	}

FOR_LOOP:
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break FOR_LOOP
		}
		if err != nil {
			return err
		}

		paths := msg.GetDestination().GetPaths()
		if paths == nil {
			continue FOR_LOOP
		}

		for _, path := range paths {
			if err := f(path); err != nil {
				return err
			}
		}
	}

	return nil
}

func addGoBGPPath(client api.GobgpApiClient, path *api.Path) error {
	req := api.AddPathRequest{
		TableType: api.TableType_GLOBAL,
		VrfId:     "",
		Path:      path,
	}

	_, err := client.AddPath(context.Background(), &req)
	return err
}

func deleteGoBGPPath(client api.GobgpApiClient, path *api.Path) error {
	req := api.DeletePathRequest{
		TableType: api.TableType_GLOBAL,
		VrfId:     "",
		Path:      path,
	}

	_, err := client.DeletePath(context.Background(), &req)
	return err
}

func modGoBGPPath(client api.GobgpApiClient, path *api.Path) error {
	if path.IsWithdraw {
		return deleteGoBGPPath(client, path)
	}

	return addGoBGPPath(client, path)
}

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
	"fabricflow/ribs/pkgs/ribscfg"
	"fabricflow/util/gobgp/apiutil"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	gobgpapi "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/pkg/packet/bgp"
	"github.com/sirupsen/logrus"
)

//
// LogLogger is interface of Logging.
//
type LogLogger interface {
	Logf(logrus.Level, string, ...interface{})
}

func logNotOutput(level logrus.Level) bool {
	return (logrus.GetLevel() < level)
}

func strBgpAddrPrefix(prefix bgp.AddrPrefixInterface) string {
	switch p := prefix.(type) {
	case *bgp.LabeledVPNIPAddrPrefix:
		return fmt.Sprintf("%s %s", p, &p.Labels)

	case *bgp.LabeledVPNIPv6AddrPrefix:
		return fmt.Sprintf("%s %s", p, &p.Labels)

	default:
		return prefix.String()
	}
}

//
// LogBgpPath output gobgpapi.Path.
//
func LogBgpPath(logger LogLogger, level logrus.Level, path *gobgpapi.Path) {
	if logNotOutput(level) {
		return
	}

	nlri, err := apiutil.GetNativeNlri(path)
	if err != nil {
		logger.Logf(level, "UnmarshalNLRI( error. %s", err)
		return
	}

	logger.Logf(level, "NLRI      : %s", strBgpAddrPrefix(nlri))
	logger.Logf(level, "Family    : %s/%s", path.Family.Afi, path.Family.Safi)

	pattrs, err := apiutil.UnmarshalPathAttributes(path.Pattrs)
	if err != nil {
		logger.Logf(level, "UnmarshalPathAttributes error. %s", err)
		return
	}

	for _, pattr := range pattrs {
		logger.Logf(level, "path attr : %s", pattr.GetType())
		logger.Logf(level, "path attr : %s", pattr)
	}

	if age, err := ptypes.Timestamp(path.Age); err != nil {
		logger.Logf(level, "age       : %s", age)
	}

	logger.Logf(level, "Withdraw  : %t", path.IsWithdraw)
	logger.Logf(level, "SourceASN : %d", path.SourceAsn)
	logger.Logf(level, "SourceID  : %s", path.SourceId)
	logger.Logf(level, "NeighborIP: %s", path.NeighborIp)
	logger.Logf(level, "NH invalid: %t", path.IsNexthopInvalid)
}

//
// LogConfig output Config.
//
func LogConfig(logger LogLogger, level logrus.Level, c *ribscfg.Config) {
	if logNotOutput(level) {
		return
	}

	if c.Ribs.Disable {
		logger.Logf(level, "Ribs,Disable     : %t", c.Ribs.Disable)
		return
	}

	logger.Logf(level, "Node.NId         : %d", c.Node.NId)
	logger.Logf(level, "Node.Label       : %d", c.Node.Label)
	logger.Logf(level, "Node.Ifname      : '%s'", c.Node.NIdIfname)
	logger.Logf(level, "NLA.API          : '%s'", c.NLA.API)
	logger.Logf(level, "Ribs,Disable     : %t", c.Ribs.Disable)
	logger.Logf(level, "Ribs,Core        : '%s'", c.Ribs.Core)
	logger.Logf(level, "Ribs.Api         : '%s'", c.Ribs.API)
	logger.Logf(level, "Ribs,SyncTime    : %d", c.Ribs.GetSyncTime())
	logger.Logf(level, "Ribs,Nexthop.Mode: '%s'", c.Ribs.Nexthops.Mode)
	logger.Logf(level, "Ribs,Nexthop.Args: '%s'", c.Ribs.Nexthops.Args)
	logger.Logf(level, "Ribs,Bgp.Addr    : '%s'", c.Ribs.Bgp.Addr)
	logger.Logf(level, "Ribs,BGP.Port    : %d", c.Ribs.Bgp.Port)
	logger.Logf(level, "Ribs,BGP.Port    : %d", c.Ribs.Bgp.Port)
	logger.Logf(level, "Ribs,BGP.Family  : '%s'", c.Ribs.Bgp.RouteFamily)
	logger.Logf(level, "Ribs,VRF.Iface   : '%s'", c.Ribs.Vrf.Iface)
	logger.Logf(level, "Ribs,VRF.RT      : '%s'", c.Ribs.Vrf.Rt)
	logger.Logf(level, "Ribs,VRF.RD      : '%s'", c.Ribs.Vrf.Rd)

}

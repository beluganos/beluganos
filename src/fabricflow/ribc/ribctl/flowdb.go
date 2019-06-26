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

package ribctl

import (
	fibcapi "fabricflow/fibc/api"

	"github.com/spf13/viper"
)

type PolicyACLMatchConfig struct {
	EthDst  string `mapstructure:"eth_dst"`
	EthType uint16 `mapstructure:"eth_type"`
	IpProto uint16 `mapstructure:"ip_proto"`
	IpDst   string `mapstructure:"ip_dst"`
	TpSrc   uint16 `mapstructure:"tp_src"`
	TpDst   uint16 `mapstructure:"tp_dst"`
	InPort  uint32 `mapstructure:"in_port"`
	Vrf     uint8  `mapstructure:"vrf"`
}

func (c *PolicyACLMatchConfig) ToAPI() *fibcapi.PolicyACLFlow_Match {
	if len(c.IpDst) != 0 {
		if _, _, err := fibcapi.ParseMaskedIP(c.IpDst); err != nil {
			return nil
		}
	}

	if len(c.EthDst) != 0 {
		if _, _, err := fibcapi.ParseMaskedMAC(c.EthDst); err != nil {
			return nil
		}
	}

	return &fibcapi.PolicyACLFlow_Match{
		IpDst:   c.IpDst,
		Vrf:     uint32(c.Vrf),
		EthType: uint32(c.EthType),
		IpProto: uint32(c.IpProto),
		TpSrc:   uint32(c.TpSrc),
		TpDst:   uint32(c.TpDst),
		EthDst:  c.EthDst,
		InPort:  c.InPort,
	}
}

type PolicyACLActionConfig struct {
	Name  string `mapstructure:"name"`
	Value uint32 `mapstructure:"value"`
}

func (a *PolicyACLActionConfig) ToAPI() *fibcapi.PolicyACLFlow_Action {
	name, err := fibcapi.ParsePolicyACLFlowActionName(a.Name)
	if err != nil {
		return nil
	}

	return &fibcapi.PolicyACLFlow_Action{
		Name:  name,
		Value: a.Value,
	}
}

type PolicyACLFlowConfig struct {
	Match  PolicyACLMatchConfig  `mapstructure:"match"`
	Action PolicyACLActionConfig `mapstructure:"action"`
}

func (f *PolicyACLFlowConfig) ToAPI() *fibcapi.PolicyACLFlow {
	m := f.Match.ToAPI()
	if m == nil {
		return nil
	}

	a := f.Action.ToAPI()
	if a == nil {
		return nil
	}

	return &fibcapi.PolicyACLFlow{Match: m, Action: a}
}

func (c *PolicyACLFlowConfig) Clone() *PolicyACLFlowConfig {
	pc := *c
	return &pc
}

func (c *PolicyACLFlowConfig) SetPort(vrf uint8, inPort uint32) *PolicyACLFlowConfig {
	c.Match.InPort = inPort
	c.Match.Vrf = vrf
	return c
}

type FlowConfig struct {
	PolicyACL []*PolicyACLFlowConfig `mapstructure:"policy_acl"`
}

const FLOWDB_BUILTIN_CONFIG = "_builtin_"

type FlowDB struct {
	Flows map[string]*FlowConfig `mapstructure:"flows"`
	viper *viper.Viper
}

func NewFlowDB() *FlowDB {
	return &FlowDB{
		Flows: map[string]*FlowConfig{},
		viper: viper.New(),
	}
}

func (c *FlowDB) SetConfigFile(path, format string) *FlowDB {
	c.viper.SetConfigFile(path)
	c.viper.SetConfigType(format)
	return c
}

func (c *FlowDB) Load() error {
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	return c.viper.Unmarshal(c)
}

func (c *FlowDB) Config(name string) *FlowConfig {
	if name == FLOWDB_BUILTIN_CONFIG {
		return NewBuiltinFlowConfig()
	}

	if cfg, ok := c.Flows[name]; ok {
		return cfg
	}

	return nil
}

func NewBuiltinFlowConfig() *FlowConfig {
	policyACL := []*PolicyACLFlowConfig{
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthType: fibcapi.ETHTYPE_LACP,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthType: fibcapi.ETHTYPE_ARP,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthType: fibcapi.ETHTYPE_IPV4,
				IpDst:   fibcapi.MCADDR_ALLROUTERS,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthType: fibcapi.ETHTYPE_IPV4,
				IpDst:   fibcapi.MCADDR_OSPF_HELLO,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthType: fibcapi.ETHTYPE_IPV4,
				IpDst:   fibcapi.MCADDR_OSPF_ALLDR,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthType: fibcapi.ETHTYPE_IPV6,
				IpDst:   fibcapi.MCADDR6_I_LOCAL,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthType: fibcapi.ETHTYPE_IPV6,
				IpDst:   fibcapi.MCADDR6_L_LOCAL,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthType: fibcapi.ETHTYPE_IPV6,
				IpDst:   fibcapi.MCADDR6_S_LOCAL,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthType: fibcapi.ETHTYPE_IPV6,
				IpDst:   fibcapi.UCADDR6_L_LOCAL,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthDst: fibcapi.HWADDR_ISIS_LEVEL1,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
		&PolicyACLFlowConfig{
			Match: PolicyACLMatchConfig{
				EthDst: fibcapi.HWADDR_ISIS_LEVEL2,
			},
			Action: PolicyACLActionConfig{
				Name: "OUTPUT",
			},
		},
	}

	return &FlowConfig{
		PolicyACL: policyACL,
	}
}

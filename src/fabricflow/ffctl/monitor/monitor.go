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

package monitor

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

const (
	monIfConfigFile  = "/etc/beluganos/ffctl_ifmon.yaml"
	monIfConfigType  = "yaml"
	monIfConfigStats = "default"
	monIfConfigDPath = "default"
	monIfInterval    = 5 * time.Second

	fibcAddr = "localhost"
	fibcPort = uint16(50081)
)

type MonitorIfaceCmd struct {
	ctrl *IfaceController
}

func NewMonitorIfaceCmd() *MonitorIfaceCmd {
	return &MonitorIfaceCmd{
		ctrl: NewIfaceController(),
	}
}

func (c *MonitorIfaceCmd) setFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringVarP(&c.ctrl.ConfigFile, "config-file", "c", monIfConfigFile, "config file name.")
	cmd.Flags().StringVarP(&c.ctrl.ConfigType, "config-type", "", monIfConfigType, "config file format.")
	cmd.Flags().StringVarP(&c.ctrl.ConfigStats, "config-ifstats", "", monIfConfigStats, "config stats name.")
	cmd.Flags().StringVarP(&c.ctrl.ConfigDPath, "config-dpath", "", monIfConfigDPath, "config ports name.")
	cmd.Flags().DurationVarP(&c.ctrl.Interval, "interval", "i", monIfInterval, "refresh interval time.")
	cmd.Flags().Uint64VarP(&c.ctrl.DpathCfg.DpID, "datapath-id", "d", 0, "datapath id.")
	cmd.Flags().Uint16VarP(&c.ctrl.HistorySize, "history-size", "", 2, "history num of average.")
	cmd.Flags().StringVarP(&c.ctrl.FibcAddr, "fibc-addr", "", fibcAddr, "fibc api address.")
	cmd.Flags().Uint16VarP(&c.ctrl.FibcPort, "fibc-port", "", fibcPort, "fibc api port.")

	return cmd
}

func (c *MonitorIfaceCmd) run(args []string) error {
	c.ctrl.DpathCfg.Ifaces = args
	return c.ctrl.Run()
}

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "monitor",
		Aliases: []string{"mon"},
		Short:   "monitoring command.",
	}

	monIf := NewMonitorIfaceCmd()

	rootCmd.AddCommand(monIf.setFlags(
		&cobra.Command{
			Use:     "interface [interface name ...]",
			Aliases: []string{"iface", "intf", "if", "i"},
			Short:   "monitor interface counters.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return monIf.run(args)
			},
		},
	))

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "show stat name list.",
			Run: func(cmd *cobra.Command, args []string) {
				keys := []string{}
				m := map[string]string{}
				for key, val := range statsNames {
					keys = append(keys, val)
					m[val] = key
				}

				sort.Strings(keys)

				fmt.Printf("%-36s | %s\n", "< counter name >", "< OpenNSL define >")
				for _, key := range keys {
					fmt.Printf("%-36s | %s\n", key, m[key])
				}
			},
		},
	)

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "sample-config",
			Short: "show sample config.",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(sampleConfig)
			},
		},
	)

	return rootCmd
}

var sampleConfig = `---

ifstats:
  default:
    - label: "Traffic statistics"
      counters:
        - label: "Input unicast packes"
          name:  ifInUcastPkts
        - label: "Input non-unicast packes"
          name:  ifInNUcastPkts
        - label: "Output unicast packets"
          name:  ifOutUcastPkts
        - label: "Output non-unicast packets"
          name:  ifOutNUcastPkts
    - label: "Error statistics"
      counters:
        - label: "Input drop packes"
          name:  ifInUcastPkts
        - label: "Input discards packets"
          name:  ifOutUcastPkts

dpaths:
  default:
    dpid:  1234
    ifaces:
      - eth1
      - eth2
      - eth3
`

var statsNames = map[string]string{
	"SPLSnmpIfInOctets":                         "ifInOctets",
	"SPLSnmpIfInUcastPkts":                      "ifInUcastPkts",
	"SPLSnmpIfInNUcastPkts":                     "ifInNUcastPkts",
	"SPLSnmpIfInDiscards":                       "ifInDiscards",
	"SPLSnmpIfInErrors":                         "ifInErrors",
	"SPLSnmpIfInUnknownProtos":                  "ifInUnknownProtos",
	"SPLSnmpIfOutOctets":                        "ifOutOctets",
	"SPLSnmpIfOutUcastPkts":                     "ifOutUcastPkts",
	"SPLSnmpIfOutNUcastPkts":                    "ifOutNUcastPkts",
	"SPLSnmpIfOutDiscards":                      "ifOutDiscards",
	"SPLSnmpIfOutErrors":                        "ifOutErrors",
	"SPLSnmpIfOutQLen":                          "ifOutQLen",
	"SPLSnmpIpInReceives":                       "ipInReceives",
	"SPLSnmpIpInHdrErrors":                      "ipInHdrErrors",
	"SPLSnmpIpForwDatagrams":                    "ipForwDatagrams",
	"SPLSnmpIpInDiscards":                       "ipInDiscards",
	"SPLSnmpDot1dBasePortDelayExceededDiscards": "dot1dBasePortDelayExceededDiscards",
	"SPLSnmpDot1dBasePortMtuExceededDiscards":   "dot1dBasePortMtuExceededDiscards",
	"SPLSnmpDot1dTpPortInFrames":                "dot1dTpPortInFrames",
	"SPLSnmpDot1dTpPortOutFrames":               "dot1dTpPortOutFrames",
	"SPLSnmpDot1dPortInDiscards":                "dot1dPortInDiscards",
	"SPLSnmpEtherStatsDropEvents":               "etherStatsDropEvents",
	"SPLSnmpEtherStatsMulticastPkts":            "etherStatsMulticastPkts",
	"SPLSnmpEtherStatsBroadcastPkts":            "etherStatsBroadcastPkts",
	"SPLSnmpEtherStatsUndersizePkts":            "etherStatsUndersizePkts",
	"SPLSnmpEtherStatsFragments":                "etherStatsFragments",
	"SPLSnmpEtherStatsPkts64Octets":             "etherStatsPkts64Octets",
	"SPLSnmpEtherStatsPkts65to127Octets":        "etherStatsPkts65to127Octets",
	"SPLSnmpEtherStatsPkts128to255Octets":       "etherStatsPkts128to255Octets",
	"SPLSnmpEtherStatsPkts256to511Octets":       "etherStatsPkts256to511Octets",
	"SPLSnmpEtherStatsPkts512to1023Octets":      "etherStatsPkts512to1023Octets",
	"SPLSnmpEtherStatsPkts1024to1518Octets":     "etherStatsPkts1024to1518Octets",
	"SPLSnmpEtherStatsOversizePkts":             "etherStatsOversizePkts",
	"SPLSnmpEtherRxOversizePkts":                "etherRxOversizePkts",
	"SPLSnmpEtherTxOversizePkts":                "etherTxOversizePkts",
	"SPLSnmpEtherStatsJabbers":                  "etherStatsJabbers",
	"SPLSnmpEtherStatsOctets":                   "etherStatsOctets",
	"SPLSnmpEtherStatsPkts":                     "etherStatsPkts",
	"SPLSnmpEtherStatsCollisions":               "etherStatsCollisions",
	"SPLSnmpEtherStatsCRCAlignErrors":           "etherStatsCRCAlignErrors",
	"SPLSnmpEtherStatsTXNoErrors":               "etherStatsTXNoErrors",
	"SPLSnmpEtherStatsRXNoErrors":               "etherStatsRXNoErrors",
	"SPLSnmpDot3StatsAlignmentErrors":           "dot3StatsAlignmentErrors",
	"SPLSnmpDot3StatsFCSErrors":                 "dot3StatsFCSErrors",
	"SPLSnmpDot3StatsSingleCollisionFrames":     "dot3StatsSingleCollisionFrames",
	"SPLSnmpDot3StatsMultipleCollisionFrames":   "dot3StatsMultipleCollisionFrames",
	"SPLSnmpDot3StatsSQETTestErrors":            "dot3StatsSQETTestErrors",
	"SPLSnmpDot3StatsDeferredTransmissions":     "dot3StatsDeferredTransmissions",
	"SPLSnmpDot3StatsLateCollisions":            "dot3StatsLateCollisions",
	"SPLSnmpDot3StatsExcessiveCollisions":       "dot3StatsExcessiveCollisions",
	"SPLSnmpDot3StatsInternalMacTransmitErrors": "dot3StatsInternalMacTransmitErrors",
	"SPLSnmpDot3StatsCarrierSenseErrors":        "dot3StatsCarrierSenseErrors",
	"SPLSnmpDot3StatsFrameTooLongs":             "dot3StatsFrameTooLongs",
	"SPLSnmpDot3StatsInternalMacReceiveErrors":  "dot3StatsInternalMacReceiveErrors",
	"SPLSnmpDot3StatsSymbolErrors":              "dot3StatsSymbolErrors",
	"SPLSnmpDot3ControlInUnknownOpcodes":        "dot3ControlInUnknownOpcodes",
	"SPLSnmpDot3InPauseFrames":                  "dot3InPauseFrames",
	"SPLSnmpDot3OutPauseFrames":                 "dot3OutPauseFrames",
	"SPLSnmpIfHCInOctets":                       "ifHCInOctets",
	"SPLSnmpIfHCInUcastPkts":                    "ifHCInUcastPkts",
	"SPLSnmpIfHCInMulticastPkts":                "ifHCInMulticastPkts",
	"SPLSnmpIfHCInBroadcastPkts":                "ifHCInBroadcastPkts",
	"SPLSnmpIfHCOutOctets":                      "ifHCOutOctets",
	"SPLSnmpIfHCOutUcastPkts":                   "ifHCOutUcastPkts",
	"SPLSnmpIfHCOutMulticastPkts":               "ifHCOutMulticastPkts",
	"SPLSnmpIfHCOutBroadcastPckts":              "ifHCOutBroadcastPckts",
	"SPLSnmpIpv6IfStatsInReceives":              "ipv6IfStatsInReceives",
	"SPLSnmpIpv6IfStatsInHdrErrors":             "ipv6IfStatsInHdrErrors",
	"SPLSnmpIpv6IfStatsInAddrErrors":            "ipv6IfStatsInAddrErrors",
	"SPLSnmpIpv6IfStatsInDiscards":              "ipv6IfStatsInDiscards",
	"SPLSnmpIpv6IfStatsOutForwDatagrams":        "ipv6IfStatsOutForwDatagrams",
	"SPLSnmpIpv6IfStatsOutDiscards":             "ipv6IfStatsOutDiscards",
	"SPLSnmpIpv6IfStatsInMcastPkts":             "ipv6IfStatsInMcastPkts",
	"SPLSnmpIpv6IfStatsOutMcastPkts":            "ipv6IfStatsOutMcastPkts",
	"SPLSnmpIfInBroadcastPkts":                  "ifInBroadcastPkts",
	"SPLSnmpIfInMulticastPkts":                  "ifInMulticastPkts",
	"SPLSnmpIfOutBroadcastPkts":                 "ifOutBroadcastPkts",
	"SPLSnmpIfOutMulticastPkts":                 "ifOutMulticastPkts",
	"SPLSnmpIeee8021PfcRequests":                "ieee8021PfcRequests",
	"SPLSnmpIeee8021PfcIndications":             "ieee8021PfcIndications",
	"NSLSnmpReceivedUndersizePkts":              "NSL_receivedUndersizePkts",
	"NSLSnmpTransmittedUndersizePkts":           "NSL_transmittedUndersizePkts",
	"NSLSnmpIPMCBridgedPckts":                   "NSL_iPMCBridgedPckts",
	"NSLSnmpIPMCRoutedPckts":                    "NSL_iPMCRoutedPckts",
	"NSLSnmpIPMCInDroppedPckts":                 "NSL_iPMCInDroppedPckts",
	"NSLSnmpIPMCOutDroppedPckts":                "NSL_iPMCOutDroppedPckts",
	"NSLSnmpEtherStatsPkts1519to1522Octets":     "NSL_etherStatsPkts1519to1522Octets",
	"NSLSnmpEtherStatsPkts1522to2047Octets":     "NSL_etherStatsPkts1522to2047Octets",
	"NSLSnmpEtherStatsPkts2048to4095Octets":     "NSL_etherStatsPkts2048to4095Octets",
	"NSLSnmpEtherStatsPkts4095to9216Octets":     "NSL_etherStatsPkts4095to9216Octets",
	"NSLSnmpReceivedPkts64Octets":               "NSL_receivedPkts64Octets",
	"NSLSnmpReceivedPkts65to127Octets":          "NSL_receivedPkts65to127Octets",
	"NSLSnmpReceivedPkts128to255Octets":         "NSL_receivedPkts128to255Octets",
	"NSLSnmpReceivedPkts256to511Octets":         "NSL_receivedPkts256to511Octets",
	"NSLSnmpReceivedPkts512to1023Octets":        "NSL_receivedPkts512to1023Octets",
	"NSLSnmpReceivedPkts1024to1518Octets":       "NSL_receivedPkts1024to1518Octets",
	"NSLSnmpReceivedPkts1519to2047Octets":       "NSL_receivedPkts1519to2047Octets",
	"NSLSnmpReceivedPkts2048to4095Octets":       "NSL_receivedPkts2048to4095Octets",
	"NSLSnmpReceivedPkts4095to9216Octets":       "NSL_receivedPkts4095to9216Octets",
	"NSLSnmpTransmittedPkts64Octets":            "NSL_transmittedPkts64Octets",
	"NSLSnmpTransmittedPkts65to127Octets":       "NSL_transmittedPkts65to127Octets",
	"NSLSnmpTransmittedPkts128to255Octets":      "NSL_transmittedPkts128to255Octets",
	"NSLSnmpTransmittedPkts256to511Octets":      "NSL_transmittedPkts256to511Octets",
	"NSLSnmpTransmittedPkts512to1023Octets":     "NSL_transmittedPkts512to1023Octets",
	"NSLSnmpTransmittedPkts1024to1518Octets":    "NSL_transmittedPkts1024to1518Octets",
	"NSLSnmpTransmittedPkts1519to2047Octets":    "NSL_transmittedPkts1519to2047Octets",
	"NSLSnmpTransmittedPkts2048to4095Octets":    "NSL_transmittedPkts2048to4095Octets",
	"NSLSnmpTransmittedPkts4095to9216Octets":    "NSL_transmittedPkts4095to9216Octets",
	"NSLSnmpTxControlCells":                     "NSL_txControlCells",
	"NSLSnmpTxDataCells":                        "NSL_txDataCells",
	"NSLSnmpTxDataBytes":                        "NSL_txDataBytes",
	"NSLSnmpRxCrcErrors":                        "NSL_rxCrcErrors",
	"NSLSnmpRxFecCorrectable":                   "NSL_rxFecCorrectable",
	"NSLSnmpRxBecCrcErrors":                     "NSL_rxBecCrcErrors",
	"NSLSnmpRxDisparityErrors":                  "NSL_rxDisparityErrors",
	"NSLSnmpRxControlCells":                     "NSL_rxControlCells",
	"NSLSnmpRxDataCells":                        "NSL_rxDataCells",
	"NSLSnmpRxDataBytes":                        "NSL_rxDataBytes",
	"NSLSnmpRxDroppedRetransmittedControl":      "NSL_rxDroppedRetransmittedControl",
	"NSLSnmpTxBecRetransmit":                    "NSL_txBecRetransmit",
	"NSLSnmpRxBecRetransmit":                    "NSL_txBecRetransmit",
	"NSLSnmpTxAsynFifoRate":                     "NSL_txAsynFifoRate",
	"NSLSnmpRxAsynFifoRate":                     "NSL_rxAsynFifoRate",
	"NSLSnmpRxFecUncorrectable":                 "NSL_rxFecUncorrectable",
	"NSLSnmpRxBecRxFault":                       "NSL_rxBecRxFault",
	"NSLSnmpRxCodeErrors":                       "NSL_rxCodeErrors",
	"NSLSnmpRxRsFecBitError":                    "NSL_rxRsFecBitError",
	"NSLSnmpRxRsFecSymbolError":                 "NSL_rxRsFecSymbolError",
	"NSLSnmpRxLlfcPrimary":                      "NSL_rxLlfcPrimary",
	"NSLSnmpRxLlfcSecondary":                    "NSL_rxLlfcSecondary",
	"NSLSnmpCustomReceive0":                     "NSL_customReceive0",
	"NSLSnmpCustomReceive1":                     "NSL_customReceive1",
	"NSLSnmpCustomReceive2":                     "NSL_customReceive2",
	"NSLSnmpCustomReceive3":                     "NSL_customReceive3",
	"NSLSnmpCustomReceive4":                     "NSL_customReceive4",
	"NSLSnmpCustomReceive5":                     "NSL_customReceive5",
	"NSLSnmpCustomReceive6":                     "NSL_customReceive6",
	"NSLSnmpCustomReceive7":                     "NSL_customReceive7",
	"NSLSnmpCustomReceive8":                     "NSL_customReceive8",
	"NSLSnmpCustomTransmit0":                    "NSL_customTransmit0",
	"NSLSnmpCustomTransmit1":                    "NSL_customTransmit1",
	"NSLSnmpCustomTransmit2":                    "NSL_customTransmit2",
	"NSLSnmpCustomTransmit3":                    "NSL_customTransmit3",
	"NSLSnmpCustomTransmit4":                    "NSL_cstomTransmit4",
	"NSLSnmpCustomTransmit5":                    "NSL_customTransmit5",
	"NSLSnmpCustomTransmit6":                    "NSL_customTransmit6",
	"NSLSnmpCustomTransmit7":                    "NSL_customTransmit7",
	"NSLSnmpCustomTransmit8":                    "NSL_customTransmit8",
	"NSLSnmpCustomTransmit9":                    "NSL_customTransmit9",
	"NSLSnmpCustomTransmit10":                   "NSL_customTransmit10",
	"NSLSnmpCustomTransmit11":                   "NSL_customTransmit11",
	"NSLSnmpCustomTransmit12":                   "NSL_customTransmit12",
	"NSLSnmpCustomTransmit13":                   "NSL_customTransmit13",
	"NSLSnmpCustomTransmit14":                   "NSL_customTransmit14",
	"SPLSnmpDot3StatsInRangeLengthError":        "dot3StatsInRangeLengthError",
	"SPLSnmpDot3OmpEmulationCRC8Errors":         "dot3OmpEmulationCRC8Errors",
	"SPLSnmpDot3MpcpRxGate":                     "dot3MpcpRxGate",
	"SPLSnmpDot3MpcpRxRegister":                 "dot3MpcpRxRegister",
	"SPLSnmpDot3MpcpTxRegRequest":               "dot3MpcpTxRegRequest",
	"SPLSnmpDot3MpcpTxRegAck":                   "dot3MpcpTxRegAck",
	"SPLSnmpDot3MpcpTxReport":                   "dot3MpcpTxReport",
	"SPLSnmpDot3EponFecCorrectedBlocks":         "dot3EponFecCorrectedBlocks",
	"SPLSnmpDot3EponFecUncorrectableBlocks":     "dot3EponFecUncorrectableBlocks",
	"NSLSnmpPonInDroppedOctets":                 "NSL_ponInDroppedOctets",
	"NSLSnmpPonOutDroppedOctets":                "NSL_ponOutDroppedOctets",
	"NSLSnmpPonInDelayedOctets":                 "NSL_pnInDelayedOctets",
	"NSLSnmpPonOutDelayedOctets":                "NSL_ponOutDelayedOctets",
	"NSLSnmpPonInDelayedHundredUs":              "NSL_ponInDelayedHundredUs",
	"NSLSnmpPonOutDelayedHundredUs":             "NSL_ponOutDelayedHundredUs",
	"NSLSnmpPonInFrameErrors":                   "NSL_ponInFrameErrors",
	"NSLSnmpPonInOamFrames":                     "NSL_ponInOamFrames",
	"NSLSnmpPonOutOamFrames":                    "NSL_ponOutOamFrames",
	"NSLSnmpPonOutUnusedOctets":                 "NSL_ponOutUnusedOctets",
	"NSLSnmpEtherStatsPkts9217to16383Octets":    "NSL_etherStatsPkts9217to16383Octets",
	"NSLSnmpReceivedPkts9217to16383Octets":      "NSL_receivedPkts9217to16383Octets",
	"NSLSnmpTransmittedPkts9217to16383Octets":   "NSL_transmittedPkts9217to16383Octets",
	"NSLSnmpRxVlanTagFrame":                     "NSL_rxVlanTagFrame",
	"NSLSnmpRxDoubleVlanTagFrame":               "NSL_rxDoubleVlanTagFrame",
	"NSLSnmpTxVlanTagFrame":                     "NSL_txVlanTagFrame",
	"NSLSnmpTxDoubleVlanTagFrame":               "NSL_TxDoubleVlanTagFrame",
	"NSLSnmpRxPFCControlFrame":                  "NSL_rxPFCControlFrame",
	"NSLSnmpTxPFCControlFrame":                  "NSL_txPFCControlFrame",
	"NSLSnmpRxPFCFrameXonPriority0":             "NSL_rxPFCFrameXonPriority0",
	"NSLSnmpRxPFCFrameXonPriority1":             "NSL_RxPFCFrameXonPriority1",
	"NSLSnmpRxPFCFrameXonPriority2":             "NSL_rxPFCFrameXonPriority2",
	"NSLSnmpRxPFCFrameXonPriority3":             "NSL_rxPFCFrameXonPriority3",
	"NSLSnmpRxPFCFrameXonPriority4":             "NSL_rxPFCFrameXonPriority4",
	"NSLSnmpRxPFCFrameXonPriority5":             "NSL_rxPFCFrameXonPriority5",
	"NSLSnmpRxPFCFrameXonPriority6":             "NSL_rxPFCFrameXonPriority6",
	"NSLSnmpRxPFCFrameXonPriority7":             "NSL_rxPFCFrameXonPriority7",
	"NSLSnmpRxPFCFramePriority0":                "NSL_rxPFCFramePriority0",
	"NSLSnmpRxPFCFramePriority1":                "NSL_rxPFCFramePriority1",
	"NSLSnmpRxPFCFramePriority2":                "NSL_rxPFCFramePriority2",
	"NSLSnmpRxPFCFramePriority3":                "NSL_rxPFCFramePriority3",
	"NSLSnmpRxPFCFramePriority4":                "NSL_rxPFCFramePriority4",
	"NSLSnmpRxPFCFramePriority5":                "NSL_rxPFCFramePriority5",
	"NSLSnmpRxPFCFramePriority6":                "NSL_rxPFCFramePriority6",
	"NSLSnmpRxPFCFramePriority7":                "NSL_rxPFCFramePriority7",
	"NSLSnmpTxPFCFramePriority0":                "NSL_txPFCFramePriority0",
	"NSLSnmpTxPFCFramePriority1":                "NSL_txPFCFramePriority1",
	"NSLSnmpTxPFCFramePriority2":                "NSL_txPFCFramePriority2",
	"NSLSnmpTxPFCFramePriority3":                "NSL_txPFCFramePriority3",
	"NSLSnmpTxPFCFramePriority4":                "NSL_txPFCFramePriority4",
	"NSLSnmpTxPFCFramePriority5":                "NSL_txPFCFramePriority5",
	"NSLSnmpTxPFCFramePriority6":                "NSL_txPFCFramePriority6",
	"NSLSnmpTxPFCFramePriority7":                "NSL_txPFCFramePriority7",
	"SPLSnmpFcmPortClass3RxFrames":              "fcmPortClass3RxFrames",
	"SPLSnmpFcmPortClass3TxFrames":              "fcmPortClass3TxFrames",
	"SPLSnmpFcmPortClass3Discards":              "fcmPortClass3Discards",
	"SPLSnmpFcmPortClass2RxFrames":              "fcmPortClass2RxFrames",
	"SPLSnmpFcmPortClass2TxFrames":              "fcmPortClass2TxFrames",
	"SPLSnmpFcmPortClass2Discards":              "fcmPortClass2Discards",
	"SPLSnmpFcmPortInvalidCRCs":                 "fcmPortInvalidCRCs",
	"SPLSnmpFcmPortDelimiterErrors":             "fcmPortDelimiterErrors",
	"NSLSnmpSampleIngressPkts":                  "NSL_sampleIngressPkts",
	"NSLSnmpSampleIngressSnapshotPkts":          "NSL_sampleIngressSnapshotPkts",
	"NSLSnmpSampleIngressSampledPkts":           "NSL_sampleIngressSampledPkts",
	"NSLSnmpSampleFlexPkts":                     "NSL_sampleFlexPkts",
	"NSLSnmpSampleFlexSnapshotPkts":             "NSL_sampleFlexSnapshotPkts",
	"NSLSnmpSampleFlexSampledPkts":              "NSL_sampleFlexSampledPkts",
	"NSLSnmpEgressProtectionDataDrop":           "NSL_egressProtectionDataDrop",
	"NSLSnmpTxE2ECCControlFrames":               "NSL_txE2ECCControlFrames",
	"NSLSnmpE2EHOLDropPkts":                     "NSL_e2EHOLDropPkts",
	"SPLSnmpEtherStatsTxCRCAlignErrors":         "etherStatsTxCRCAlignErrors",
	"SPLSnmpEtherStatsTxJabbers":                "etherStatsTxJabbers",
	"NSLSnmpMacMergeTxFrag":                     "NSL_macMergeTxFrag",
	"NSLSnmpMacMergeTxVerifyFrame":              "NSL_macMergeTxVerifyFrame",
	"NSLSnmpMacMergeTxReplyFrame":               "NSL_macMergeTxReplyFrame",
	"NSLSnmpMacMergeRxFrameAssErrors":           "NSL_macMergeRxFrameAssErrors",
	"NSLSnmpMacMergeRxFrameSmdErrors":           "NSL_macMergeRxFrameSmdErrors",
	"NSLSnmpMacMergeRxFrameAss":                 "NSL_macMergeRxFrameAss",
	"NSLSnmpMacMergeRxFrag":                     "NSL_macMergeRxFrag",
	"NSLSnmpMacMergeRxVerifyFrame":              "NSL_macMergeRxVerifyFrame",
	"NSLSnmpMacMergeRxReplyFrame":               "NSL_macMergeRxReplyFrame",
	"NSLSnmpMacMergeRxFinalFragSizeError":       "NSL_macMergeRxFinalFragSizeError",
	"NSLSnmpMacMergeRxFragSizeError":            "NSL_macMergeRxFragSizeError",
	"NSLSnmpMacMergeRxDiscard":                  "NSL_macMergeRxDiscard",
	"NSLSnmpMacMergeHoldCount":                  "NSL_macMergeHoldCount",
	"NSLSnmpRxBipErrorCount":                    "NSL_rxBipErrorCount",
}

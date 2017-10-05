// -*- coding: utf-8 -*-

// Copyright (C) 2017 Nippon Telegraph and Telephone Corporation.
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

package nladbm

import (
	"gonla/nlamsg"
	"gonla/nlamsg/nlalink"
	"net"
	"testing"
)

func TestVpn_x1(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10016, nil), 10, 0)

	tbl := newVpnTable()

	// Insert
	if old := tbl.Insert(vpn1); old != nil {
		t.Errorf("VpnTable Insert unmatch. %v", old)
	}
	if v := len(tbl.Vpns); v != 1 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	// Select
	key := NewVpnKey(10, dst1, gw1)
	sel := tbl.Select(key)
	if sel == nil {
		t.Errorf("VpnTable Select unmatch. %v", sel)
	}
	if !sel.NetGw().Equal(gw1) || sel.GetIPNet().String() != dst1.String() {
		t.Errorf("VpnTable Select unmatch. %v", sel)
	}

	// Delete
	del := tbl.Delete(key)
	if del == nil {
		t.Errorf("VpnTable Delete unmatch. %v", del)
	}
	if !del.NetGw().Equal(gw1) || del.GetIPNet().String() != dst1.String() {
		t.Errorf("VpnTable Select unmatch. %v", del)
	}
	if v := len(tbl.Vpns); v != 0 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 0 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 0 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	// Delete(not found)
	del2 := tbl.Delete(key)
	if del2 != nil {
		t.Errorf("VpnTable Delete unmatch. %v", del2)
	}
}

func TestVpn_dst_x2_gw_x2(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	dst2 := &net.IPNet{
		IP:   []byte{1, 1, 2, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	gw2 := net.IP([]byte{10, 0, 1, 2})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10016, nil), 10, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst2, gw2, 10016, nil), 10, 0)

	tbl := newVpnTable()

	// Insert
	tbl.Insert(vpn1)
	if old := tbl.Insert(vpn2); old != nil {
		t.Errorf("VpnTable Insert unmatch. %v", old)
	}
	if v := len(tbl.Vpns); v != 2 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 2 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 2 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	// Select
	key1 := NewVpnKey(10, dst1, gw1)
	sel1 := tbl.Select(key1)
	if sel1 == nil {
		t.Errorf("VpnTable Select unmatch. %v", sel1)
	}
	if !sel1.NetGw().Equal(gw1) || sel1.GetIPNet().String() != dst1.String() {
		t.Errorf("VpnTable Select unmatch. %v", sel1)
	}

	key2 := NewVpnKey(10, dst2, gw2)
	sel2 := tbl.Select(key2)
	if sel2 == nil {
		t.Errorf("VpnTable Select unmatch. %v", sel2)
	}
	if !sel2.NetGw().Equal(gw2) || sel2.GetIPNet().String() != dst2.String() {
		t.Errorf("VpnTable Select unmatch. %v", sel2)
	}

	// Delete
	del1 := tbl.Delete(key1)
	if del1 == nil {
		t.Errorf("VpnTable Delete unmatch. %v", del1)
	}
	if !del1.NetGw().Equal(gw1) || del1.GetIPNet().String() != dst1.String() {
		t.Errorf("VpnTable Select unmatch. %v", del1)
	}
	if v := len(tbl.Vpns); v != 1 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	del2 := tbl.Delete(key2)
	if del2 == nil {
		t.Errorf("VpnTable Delete unmatch. %v", del2)
	}
	if !del2.NetGw().Equal(gw2) || del2.GetIPNet().String() != dst2.String() {
		t.Errorf("VpnTable Select unmatch. %v", del2)
	}
	if v := len(tbl.Vpns); v != 0 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 0 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 0 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	// Delete(not found)
	if del := tbl.Delete(key1); del != nil {
		t.Errorf("VpnTable Delete unmatch. %v", del)
	}
	if del := tbl.Delete(key2); del != nil {
		t.Errorf("VpnTable Delete unmatch. %v", del)
	}
}

func TestVpn_dst_x1_gw_x2(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	dst2 := dst1
	gw1 := net.IP([]byte{10, 0, 1, 1})
	gw2 := net.IP([]byte{10, 0, 1, 2})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10016, nil), 10, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst2, gw2, 10016, nil), 10, 0)

	tbl := newVpnTable()

	// Insert
	tbl.Insert(vpn1)
	if old := tbl.Insert(vpn2); old != nil {
		t.Errorf("VpnTable Insert unmatch. %v", old)
	}
	if v := len(tbl.Vpns); v != 2 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 2 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 2 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	// Select
	key1 := NewVpnKey(10, dst1, gw1)
	sel1 := tbl.Select(key1)
	if sel1 == nil {
		t.Errorf("VpnTable Select unmatch. %v", sel1)
	}
	if !sel1.NetGw().Equal(gw1) || sel1.GetIPNet().String() != dst1.String() {
		t.Errorf("VpnTable Select unmatch. %v", sel1)
	}

	key2 := NewVpnKey(10, dst2, gw2)
	sel2 := tbl.Select(key2)
	if sel2 == nil {
		t.Errorf("VpnTable Select unmatch. %v", sel2)
	}
	if !sel2.NetGw().Equal(gw2) || sel2.GetIPNet().String() != dst2.String() {
		t.Errorf("VpnTable Select unmatch. %v", sel2)
	}

	// Delete
	del1 := tbl.Delete(key1)
	if del1 == nil {
		t.Errorf("VpnTable Delete unmatch. %v", del1)
	}
	if !del1.NetGw().Equal(gw1) || del1.GetIPNet().String() != dst1.String() {
		t.Errorf("VpnTable Select unmatch. %v", del1)
	}
	if v := len(tbl.Vpns); v != 1 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	del2 := tbl.Delete(key2)
	if del2 == nil {
		t.Errorf("VpnTable Delete unmatch. %v", del2)
	}
	if !del2.NetGw().Equal(gw2) || del2.GetIPNet().String() != dst2.String() {
		t.Errorf("VpnTable Select unmatch. %v", del2)
	}
	if v := len(tbl.Vpns); v != 0 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 0 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 0 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	// Delete(not found)
	if del := tbl.Delete(key1); del != nil {
		t.Errorf("VpnTable Delete unmatch. %v", del)
	}
	if del := tbl.Delete(key2); del != nil {
		t.Errorf("VpnTable Delete unmatch. %v", del)
	}
}

func TestVpn_dst_x2_gw_x1(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	dst2 := &net.IPNet{
		IP:   []byte{1, 1, 2, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	gw2 := net.IP([]byte{10, 0, 1, 1})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10016, nil), 10, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst2, gw2, 10016, nil), 10, 0)

	tbl := newVpnTable()

	// Insert
	tbl.Insert(vpn1)
	if old := tbl.Insert(vpn2); old != nil {
		t.Errorf("VpnTable Insert unmatch. %v", old)
	}
	if v := len(tbl.Vpns); v != 2 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	// Select
	key1 := NewVpnKey(10, dst1, gw1)
	sel1 := tbl.Select(key1)
	if sel1 == nil {
		t.Errorf("VpnTable Select unmatch. %v", sel1)
	}
	if !sel1.NetGw().Equal(gw1) || sel1.GetIPNet().String() != dst1.String() {
		t.Errorf("VpnTable Select unmatch. %v", sel1)
	}

	key2 := NewVpnKey(10, dst2, gw2)
	sel2 := tbl.Select(key2)
	if sel2 == nil {
		t.Errorf("VpnTable Select unmatch. %v", sel2)
	}
	if !sel2.NetGw().Equal(gw2) || sel2.GetIPNet().String() != dst2.String() {
		t.Errorf("VpnTable Select unmatch. %v", sel2)
	}

	// Delete
	del1 := tbl.Delete(key1)
	if del1 == nil {
		t.Errorf("VpnTable Delete unmatch. %v", del1)
	}
	if !del1.NetGw().Equal(gw1) || del1.GetIPNet().String() != dst1.String() {
		t.Errorf("VpnTable Select unmatch. %v", del1)
	}
	if v := len(tbl.Vpns); v != 1 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 1 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	del2 := tbl.Delete(key2)
	if del2 == nil {
		t.Errorf("VpnTable Delete unmatch. %v", del2)
	}
	if !del2.NetGw().Equal(gw2) || del2.GetIPNet().String() != dst2.String() {
		t.Errorf("VpnTable Select unmatch. %v", del2)
	}
	if v := len(tbl.Vpns); v != 0 {
		t.Errorf("VpnTable size unmatch. %d", v)
	}
	if v := len(tbl.GwIdx.Entry); v != 0 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}
	if v := len(tbl.VpnGwIdx.Entry); v != 0 {
		t.Errorf("VpnTable.GwIndex size unmatch. %d", v)
	}

	// Delete(not found)
	if del := tbl.Delete(key1); del != nil {
		t.Errorf("VpnTable Delete unmatch. %v", del)
	}
	if del := tbl.Delete(key2); del != nil {
		t.Errorf("VpnTable Delete unmatch. %v", del)
	}
}

func TestVpnGwIndex_Insert(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10016, nil), 10, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10016, nil), 10, 0)

	index := NewVpnGwIndex()

	index.Insert(vpn1)

	if v := len(index.Entry); v != 1 {
		t.Errorf("VpnGwIndex Insert size unmatch. %d", v)
	}

	e, _ := index.Select(gw1)

	if v := len(e.Keys); v != 1 {
		t.Errorf("VpnGwIndex Insert entry size unmatch. %d", v)
	}

	// duplicate entry

	index.Insert(vpn2)

	if v := len(index.Entry); v != 1 {
		t.Errorf("VpnGwIndex Insert size unmatch. %d", v)
	}

	e, _ = index.Select(gw1)

	if v := len(e.Keys); v != 1 {
		t.Errorf("VpnGwIndex Insert entry size unmatch. %d", v)
	}
}

func TestVpnGwIndex_Insert_gw_x1(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	dst2 := &net.IPNet{
		IP:   []byte{1, 1, 2, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10, nil), 0, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst2, gw1, 10, nil), 0, 0)

	index := NewVpnGwIndex()

	index.Insert(vpn1)
	index.Insert(vpn2)

	if v := len(index.Entry); v != 1 {
		t.Errorf("VpnGwIndex Insert size unmatch. %d", v)
	}

	e, _ := index.Select(gw1)

	if v := len(e.Keys); v != 2 {
		t.Errorf("VpnGwIndex Insert entry size unmatch. %d", v)
	}
}

func TestVpnGwIndex_Insert_gw_x2(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	dst2 := &net.IPNet{
		IP:   []byte{1, 1, 2, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	gw2 := net.IP([]byte{10, 0, 1, 2})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10, nil), 0, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst2, gw2, 10, nil), 0, 0)

	index := NewVpnGwIndex()

	index.Insert(vpn1)
	index.Insert(vpn2)

	if v := len(index.Entry); v != 2 {
		t.Errorf("VpnGwIndex Insert size unmatch. %d", v)
	}

	e, _ := index.Select(gw1)

	if v := len(e.Keys); v != 1 {
		t.Errorf("VpnGwIndex Insert entry size unmatch. %d", v)
	}

	e, _ = index.Select(gw2)

	if v := len(e.Keys); v != 1 {
		t.Errorf("VpnGwIndex Insert entry size unmatch. %d", v)
	}
}

func TestVpnGwIndex_Insert_dst_x1_gw_x2(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	gw2 := net.IP([]byte{10, 0, 1, 2})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10, nil), 0, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw2, 10, nil), 0, 0)

	index := NewVpnGwIndex()

	index.Insert(vpn1)
	index.Insert(vpn2)

	if v := len(index.Entry); v != 2 {
		t.Errorf("VpnGwIndex Insert size unmatch. %d", v)
	}

	e, _ := index.Select(gw1)

	if v := len(e.Keys); v != 1 {
		t.Errorf("VpnGwIndex Insert entry size unmatch. %d", v)
	}

	e, _ = index.Select(gw2)

	if v := len(e.Keys); v != 1 {
		t.Errorf("VpnGwIndex Insert entry size unmatch. %d", v)
	}
}

func TestVpnGwIndex_Delete(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10, nil), 0, 0)

	index := NewVpnGwIndex()

	index.Insert(vpn1)
	if v := len(index.Entry); v != 1 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}

	index.Delete(vpn1)
	if v := len(index.Entry); v != 0 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}

	// not found
	index.Delete(vpn1)
	if v := len(index.Entry); v != 0 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}
}

func TestVpnGwIndex_Delete_gw_x1(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	dst2 := &net.IPNet{
		IP:   []byte{1, 1, 2, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10, nil), 0, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst2, gw1, 10, nil), 0, 0)

	index := NewVpnGwIndex()

	index.Insert(vpn1)
	index.Insert(vpn2)
	if v := len(index.Entry); v != 1 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}

	index.Delete(vpn1)
	if v := len(index.Entry); v != 1 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}

	index.Delete(vpn2)
	if v := len(index.Entry); v != 0 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}
}

func TestVpnGwIndex_Delete_gw_x2(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	dst2 := &net.IPNet{
		IP:   []byte{1, 1, 2, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	gw2 := net.IP([]byte{10, 0, 1, 2})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10, nil), 0, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst2, gw2, 10, nil), 0, 0)

	index := NewVpnGwIndex()

	index.Insert(vpn1)
	index.Insert(vpn2)
	if v := len(index.Entry); v != 2 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}

	index.Delete(vpn1)
	if v := len(index.Entry); v != 1 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}

	index.Delete(vpn2)
	if v := len(index.Entry); v != 0 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}
}

func TestVpnGwIndex_Delete_dst_x1_gw_x2(t *testing.T) {
	dst1 := &net.IPNet{
		IP:   []byte{1, 1, 1, 0},
		Mask: []byte{255, 255, 255, 0},
	}
	gw1 := net.IP([]byte{10, 0, 1, 1})
	gw2 := net.IP([]byte{10, 0, 1, 2})
	vpn1 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw1, 10, nil), 0, 0)
	vpn2 := nlamsg.NewVpn(nlalink.NewVpn(dst1, gw2, 10, nil), 0, 0)

	index := NewVpnGwIndex()

	index.Insert(vpn1)
	index.Insert(vpn2)
	if v := len(index.Entry); v != 2 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}

	index.Delete(vpn1)
	if v := len(index.Entry); v != 1 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}

	index.Delete(vpn2)
	if v := len(index.Entry); v != 0 {
		t.Errorf("VpnGwIndex Delete size unmatch. %d", v)
	}
}

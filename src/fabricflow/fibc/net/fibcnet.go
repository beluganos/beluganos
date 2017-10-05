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

package fibcnet

import (
	"bytes"
	"encoding/binary"
	"io"
)

const HEADER_LEN = 8

type Header struct {
	Type   uint16
	Length uint16
	Xid    uint32
}

func ParseHeader(data []byte) (*Header, error) {
	r := bytes.NewReader(data)
	h := &Header{}
	err := binary.Read(r, binary.BigEndian, h)
	return h, err
}

func ReadBytes(r io.Reader, length int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := io.CopyN(buf, r, length)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ReadHeader(r io.Reader) (*Header, error) {
	b, err := ReadBytes(r, HEADER_LEN)
	if err != nil {
		return nil, err
	}
	return ParseHeader(b)
}

func Read(r io.Reader) (*Header, []byte, error) {
	h, err := ReadHeader(r)
	if err != nil {
		return nil, nil, err
	}

	b, err := ReadBytes(r, int64(h.Length))

	return h, b, err
}

type Message interface {
	Bytes() ([]byte, error)
	Type() uint16
}

func WriteMessage(w io.Writer, msg Message, xid uint32) error {
	b, err := msg.Bytes()
	if err != nil {
		return err
	}
	return Write(w, msg.Type(), xid, b)
}

func Write(w io.Writer, t uint16, xid uint32, b []byte) error {
	h := Header{
		Type:   t,
		Length: uint16(len(b)),
		Xid:    xid,
	}
	if err := binary.Write(w, binary.BigEndian, &h); err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, b)
}

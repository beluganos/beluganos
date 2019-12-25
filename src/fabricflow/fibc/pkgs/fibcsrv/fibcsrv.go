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

package fibcsrv

import (
	fibcapi "fabricflow/fibc/api"
	"fabricflow/fibc/pkgs/fibccfg"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

//
// Server is fibc server
//
type Server struct {
	dbctl  *DBCtl
	server *grpc.Server

	logMonitor bool

	log *log.Entry
}

//
// NewServer returns new server.
//
func NewServer() *Server {
	return &Server{
		dbctl:  NewDBCtl(),
		server: grpc.NewServer(),
		log:    log.WithFields(log.Fields{"module": "server"}),
	}
}

//
// SetLogMonitor set log monitoring enable or not.
//
func (s *Server) SetLogMonitor(enable bool) {
	s.logMonitor = enable
}

func (s *Server) newAPAPIServer() *APAPIServer {
	return NewAPAPIServer(NewAPCtl(s.dbctl), s.server)
}

func (s *Server) newVMAPIServer() *VMAPIServer {
	return NewVMAPIServer(NewVMCtl(s.dbctl), s.server)
}

func (s *Server) newVSAPIServer() *VSAPIServer {
	return NewVSAPIServer(NewVSCtl(s.dbctl), s.server)
}

func (s *Server) newDPAPIServer() *DPAPIServer {
	return NewDPAPIServer(NewDPCtl(s.dbctl), s.server)
}

func (s *Server) newLogHook() *FibcLogHook {
	return AddFibcLogHook(s.dbctl)
}

//
// SetConfig add config to db
//
func (s *Server) SetConfig(cfg *fibccfg.Config) {
	s.dbctl.RegisterConfig(cfg)
}

//
// SetNetconfConfig set path, type of netconf config file
//
func (s *Server) SetNetconfConfig(configPath, configType string) {
	if err := s.dbctl.NetconfConfig().SetConfig(configPath, configType); err != nil {
		s.log.Errorf("netconf config initialize error. %s", err)
	}
}

//
// Serve creates servrs.
//
func (s *Server) Serve(lis net.Listener) {
	if s.logMonitor {
		s.newLogHook()
	}
	s.newVMAPIServer()
	s.newDPAPIServer()
	s.newVSAPIServer()
	s.newAPAPIServer()

	s.log.Infof("started.")
	s.server.Serve(lis)
}

//
// FibcLogHook is hook for logrus.
//
type FibcLogHook struct {
	db *DBCtl
}

//
// NewFibcLogHook returns new FibcLogHook
//
func NewFibcLogHook(db *DBCtl) *FibcLogHook {
	return &FibcLogHook{
		db: db,
	}
}

//
// AddFibcLogHook add FibcLogHook to logrus.
//
func AddFibcLogHook(db *DBCtl) *FibcLogHook {
	h := NewFibcLogHook(db)
	log.AddHook(h)

	return h
}

//
// Levels returns levels to hook.
//
func (h *FibcLogHook) Levels() []log.Level {
	return log.AllLevels
}

//
// Fire process log message.
//
func (h *FibcLogHook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		line = entry.Message
	}

	msg := fibcapi.ApMonitorReplyLog{
		Line:  line,
		Level: uint32(entry.Level),
		Time:  entry.Time.UnixNano(),
	}
	h.db.SendAPMonitorReplyLog(&msg)

	return nil
}

//
// GrpcRemoteHostPort returns remote (address, port).
//
func GrpcRemoteHostPort(stream grpc.ServerStream) (string, string) {
	if stream == nil || stream.Context() == nil {
		return "", ""
	}

	if p, ok := peer.FromContext(stream.Context()); ok {
		addr := p.Addr.String()
		if host, port, err := net.SplitHostPort(addr); err == nil {
			return host, port
		}
		return addr, ""
	}

	return "", ""
}

# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: ribtapi.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='ribtapi.proto',
  package='ribtapi',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\rribtapi.proto\x12\x07ribtapi\"S\n\x0bTunnelRoute\x12\x0e\n\x06prefix\x18\x01 \x01(\t\x12\x0f\n\x07nexthop\x18\x02 \x01(\t\x12\x0e\n\x06\x66\x61mily\x18\x03 \x01(\r\x12\x13\n\x0btunnel_type\x18\x04 \x01(\x05\"%\n\x11GetTunnelsRequest\x12\x10\n\x08key_type\x18\x01 \x01(\t\"\xc5\x01\n\x0fGetTunnelsReply\x12\n\n\x02id\x18\x01 \x01(\r\x12\x0c\n\x04type\x18\x02 \x01(\x05\x12\x0e\n\x06remote\x18\x03 \x01(\t\x12\r\n\x05local\x18\x04 \x01(\t\x12\x34\n\x06routes\x18\x05 \x03(\x0b\x32$.ribtapi.GetTunnelsReply.RoutesEntry\x1a\x43\n\x0bRoutesEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12#\n\x05value\x18\x02 \x01(\x0b\x32\x14.ribtapi.TunnelRoute:\x02\x38\x01\x32Q\n\x07RIBTApi\x12\x46\n\nGetTunnels\x12\x1a.ribtapi.GetTunnelsRequest\x1a\x18.ribtapi.GetTunnelsReply\"\x00\x30\x01\x62\x06proto3')
)




_TUNNELROUTE = _descriptor.Descriptor(
  name='TunnelRoute',
  full_name='ribtapi.TunnelRoute',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='prefix', full_name='ribtapi.TunnelRoute.prefix', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='nexthop', full_name='ribtapi.TunnelRoute.nexthop', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='family', full_name='ribtapi.TunnelRoute.family', index=2,
      number=3, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='tunnel_type', full_name='ribtapi.TunnelRoute.tunnel_type', index=3,
      number=4, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=26,
  serialized_end=109,
)


_GETTUNNELSREQUEST = _descriptor.Descriptor(
  name='GetTunnelsRequest',
  full_name='ribtapi.GetTunnelsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key_type', full_name='ribtapi.GetTunnelsRequest.key_type', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=111,
  serialized_end=148,
)


_GETTUNNELSREPLY_ROUTESENTRY = _descriptor.Descriptor(
  name='RoutesEntry',
  full_name='ribtapi.GetTunnelsReply.RoutesEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='ribtapi.GetTunnelsReply.RoutesEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='ribtapi.GetTunnelsReply.RoutesEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=_b('8\001'),
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=281,
  serialized_end=348,
)

_GETTUNNELSREPLY = _descriptor.Descriptor(
  name='GetTunnelsReply',
  full_name='ribtapi.GetTunnelsReply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='ribtapi.GetTunnelsReply.id', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='ribtapi.GetTunnelsReply.type', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='remote', full_name='ribtapi.GetTunnelsReply.remote', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='local', full_name='ribtapi.GetTunnelsReply.local', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='routes', full_name='ribtapi.GetTunnelsReply.routes', index=4,
      number=5, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_GETTUNNELSREPLY_ROUTESENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=151,
  serialized_end=348,
)

_GETTUNNELSREPLY_ROUTESENTRY.fields_by_name['value'].message_type = _TUNNELROUTE
_GETTUNNELSREPLY_ROUTESENTRY.containing_type = _GETTUNNELSREPLY
_GETTUNNELSREPLY.fields_by_name['routes'].message_type = _GETTUNNELSREPLY_ROUTESENTRY
DESCRIPTOR.message_types_by_name['TunnelRoute'] = _TUNNELROUTE
DESCRIPTOR.message_types_by_name['GetTunnelsRequest'] = _GETTUNNELSREQUEST
DESCRIPTOR.message_types_by_name['GetTunnelsReply'] = _GETTUNNELSREPLY
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

TunnelRoute = _reflection.GeneratedProtocolMessageType('TunnelRoute', (_message.Message,), {
  'DESCRIPTOR' : _TUNNELROUTE,
  '__module__' : 'ribtapi_pb2'
  # @@protoc_insertion_point(class_scope:ribtapi.TunnelRoute)
  })
_sym_db.RegisterMessage(TunnelRoute)

GetTunnelsRequest = _reflection.GeneratedProtocolMessageType('GetTunnelsRequest', (_message.Message,), {
  'DESCRIPTOR' : _GETTUNNELSREQUEST,
  '__module__' : 'ribtapi_pb2'
  # @@protoc_insertion_point(class_scope:ribtapi.GetTunnelsRequest)
  })
_sym_db.RegisterMessage(GetTunnelsRequest)

GetTunnelsReply = _reflection.GeneratedProtocolMessageType('GetTunnelsReply', (_message.Message,), {

  'RoutesEntry' : _reflection.GeneratedProtocolMessageType('RoutesEntry', (_message.Message,), {
    'DESCRIPTOR' : _GETTUNNELSREPLY_ROUTESENTRY,
    '__module__' : 'ribtapi_pb2'
    # @@protoc_insertion_point(class_scope:ribtapi.GetTunnelsReply.RoutesEntry)
    })
  ,
  'DESCRIPTOR' : _GETTUNNELSREPLY,
  '__module__' : 'ribtapi_pb2'
  # @@protoc_insertion_point(class_scope:ribtapi.GetTunnelsReply)
  })
_sym_db.RegisterMessage(GetTunnelsReply)
_sym_db.RegisterMessage(GetTunnelsReply.RoutesEntry)


_GETTUNNELSREPLY_ROUTESENTRY._options = None

_RIBTAPI = _descriptor.ServiceDescriptor(
  name='RIBTApi',
  full_name='ribtapi.RIBTApi',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=350,
  serialized_end=431,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetTunnels',
    full_name='ribtapi.RIBTApi.GetTunnels',
    index=0,
    containing_service=None,
    input_type=_GETTUNNELSREQUEST,
    output_type=_GETTUNNELSREPLY,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_RIBTAPI)

DESCRIPTOR.services_by_name['RIBTApi'] = _RIBTAPI

# @@protoc_insertion_point(module_scope)

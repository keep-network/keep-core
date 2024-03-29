// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: pkg/net/gen/pb/handshake.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Envelope contains a marshalled message, as well as a signature over the
// the contents of the message (to ensure an adversary hasn't tampered
// with the contents).
type HandshakeEnvelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The marshalled message.
	Message []byte `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	// Signature of the message.
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	// Peer id of the message creator
	PeerID []byte `protobuf:"bytes,3,opt,name=peerID,proto3" json:"peerID,omitempty"`
}

func (x *HandshakeEnvelope) Reset() {
	*x = HandshakeEnvelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_gen_pb_handshake_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandshakeEnvelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandshakeEnvelope) ProtoMessage() {}

func (x *HandshakeEnvelope) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_gen_pb_handshake_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandshakeEnvelope.ProtoReflect.Descriptor instead.
func (*HandshakeEnvelope) Descriptor() ([]byte, []int) {
	return file_pkg_net_gen_pb_handshake_proto_rawDescGZIP(), []int{0}
}

func (x *HandshakeEnvelope) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *HandshakeEnvelope) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *HandshakeEnvelope) GetPeerID() []byte {
	if x != nil {
		return x.PeerID
	}
	return nil
}

// Act1Message is sent in the first handshake act by the initiator to the
// responder. It contains randomly generated `nonce1`, an 8-byte (64-bit)
// unsigned integer, and the protocol identifier.
//
// Act1Message should be signed with initiator's static private key.
type Act1Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// nonce by initiator; 8-byte (64-bit) nonce as bytes
	Nonce []byte `protobuf:"bytes,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	// the identifier of the protocol the initiator is executing
	Protocol string `protobuf:"bytes,2,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *Act1Message) Reset() {
	*x = Act1Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_gen_pb_handshake_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Act1Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Act1Message) ProtoMessage() {}

func (x *Act1Message) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_gen_pb_handshake_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Act1Message.ProtoReflect.Descriptor instead.
func (*Act1Message) Descriptor() ([]byte, []int) {
	return file_pkg_net_gen_pb_handshake_proto_rawDescGZIP(), []int{1}
}

func (x *Act1Message) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *Act1Message) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

// Act2Message is sent in the second handshake act by the responder to the
// initiator. It contains randomly generated `nonce2`, an 8-byte unsigned
// integer and `challenge` which is a result of SHA256 on the concatenated
// bytes of `nonce1` and `nonce2`, and the protocol identifier.
//
// Act2Message should be signed with responder's static private key.
type Act2Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// nonce from responder; 8-byte (64-bit) nonce as bytes
	Nonce []byte `protobuf:"bytes,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	// bytes of sha256(nonce1||nonce2)
	Challenge []byte `protobuf:"bytes,2,opt,name=challenge,proto3" json:"challenge,omitempty"`
	// the identifier of the protocol the responder is executing
	Protocol string `protobuf:"bytes,3,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *Act2Message) Reset() {
	*x = Act2Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_gen_pb_handshake_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Act2Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Act2Message) ProtoMessage() {}

func (x *Act2Message) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_gen_pb_handshake_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Act2Message.ProtoReflect.Descriptor instead.
func (*Act2Message) Descriptor() ([]byte, []int) {
	return file_pkg_net_gen_pb_handshake_proto_rawDescGZIP(), []int{2}
}

func (x *Act2Message) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *Act2Message) GetChallenge() []byte {
	if x != nil {
		return x.Challenge
	}
	return nil
}

func (x *Act2Message) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

// Act1Message is sent in the first handshake act by the initiator to the
// responder. It contains randomly generated `nonce1`, an 8-byte (64-bit)
// unsigned integer.
//
// Act1Message should be signed with initiator's static private key.
type Act3Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// bytes of sha256(nonce1||nonce2)
	Challenge []byte `protobuf:"bytes,1,opt,name=challenge,proto3" json:"challenge,omitempty"`
}

func (x *Act3Message) Reset() {
	*x = Act3Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_gen_pb_handshake_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Act3Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Act3Message) ProtoMessage() {}

func (x *Act3Message) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_gen_pb_handshake_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Act3Message.ProtoReflect.Descriptor instead.
func (*Act3Message) Descriptor() ([]byte, []int) {
	return file_pkg_net_gen_pb_handshake_proto_rawDescGZIP(), []int{3}
}

func (x *Act3Message) GetChallenge() []byte {
	if x != nil {
		return x.Challenge
	}
	return nil
}

var File_pkg_net_gen_pb_handshake_proto protoreflect.FileDescriptor

var file_pkg_net_gen_pb_handshake_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x6b, 0x67, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62,
	0x2f, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x03, 0x6e, 0x65, 0x74, 0x22, 0x63, 0x0a, 0x11, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61,
	0x6b, 0x65, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x44, 0x22, 0x3f, 0x0a, 0x0b, 0x41, 0x63,
	0x74, 0x31, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e,
	0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x5d, 0x0a, 0x0b, 0x41,
	0x63, 0x74, 0x32, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f,
	0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x2b, 0x0a, 0x0b, 0x41, 0x63,
	0x74, 0x33, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x68, 0x61,
	0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x68,
	0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_net_gen_pb_handshake_proto_rawDescOnce sync.Once
	file_pkg_net_gen_pb_handshake_proto_rawDescData = file_pkg_net_gen_pb_handshake_proto_rawDesc
)

func file_pkg_net_gen_pb_handshake_proto_rawDescGZIP() []byte {
	file_pkg_net_gen_pb_handshake_proto_rawDescOnce.Do(func() {
		file_pkg_net_gen_pb_handshake_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_net_gen_pb_handshake_proto_rawDescData)
	})
	return file_pkg_net_gen_pb_handshake_proto_rawDescData
}

var file_pkg_net_gen_pb_handshake_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_net_gen_pb_handshake_proto_goTypes = []interface{}{
	(*HandshakeEnvelope)(nil), // 0: net.HandshakeEnvelope
	(*Act1Message)(nil),       // 1: net.Act1Message
	(*Act2Message)(nil),       // 2: net.Act2Message
	(*Act3Message)(nil),       // 3: net.Act3Message
}
var file_pkg_net_gen_pb_handshake_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_net_gen_pb_handshake_proto_init() }
func file_pkg_net_gen_pb_handshake_proto_init() {
	if File_pkg_net_gen_pb_handshake_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_net_gen_pb_handshake_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandshakeEnvelope); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_net_gen_pb_handshake_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Act1Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_net_gen_pb_handshake_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Act2Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_net_gen_pb_handshake_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Act3Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_net_gen_pb_handshake_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_net_gen_pb_handshake_proto_goTypes,
		DependencyIndexes: file_pkg_net_gen_pb_handshake_proto_depIdxs,
		MessageInfos:      file_pkg_net_gen_pb_handshake_proto_msgTypes,
	}.Build()
	File_pkg_net_gen_pb_handshake_proto = out.File
	file_pkg_net_gen_pb_handshake_proto_rawDesc = nil
	file_pkg_net_gen_pb_handshake_proto_goTypes = nil
	file_pkg_net_gen_pb_handshake_proto_depIdxs = nil
}

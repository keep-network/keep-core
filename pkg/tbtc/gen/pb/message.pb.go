// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.7.1
// source: pkg/tbtc/gen/pb/message.proto

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

type StopPill struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AttemptNumber uint64 `protobuf:"varint,1,opt,name=attemptNumber,proto3" json:"attemptNumber,omitempty"`
}

func (x *StopPill) Reset() {
	*x = StopPill{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_tbtc_gen_pb_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopPill) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopPill) ProtoMessage() {}

func (x *StopPill) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_tbtc_gen_pb_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopPill.ProtoReflect.Descriptor instead.
func (*StopPill) Descriptor() ([]byte, []int) {
	return file_pkg_tbtc_gen_pb_message_proto_rawDescGZIP(), []int{0}
}

func (x *StopPill) GetAttemptNumber() uint64 {
	if x != nil {
		return x.AttemptNumber
	}
	return 0
}

type SigningAnnouncementMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderID      uint32 `protobuf:"varint,1,opt,name=senderID,proto3" json:"senderID,omitempty"`
	Message       []byte `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	AttemptNumber uint64 `protobuf:"varint,3,opt,name=attemptNumber,proto3" json:"attemptNumber,omitempty"`
}

func (x *SigningAnnouncementMessage) Reset() {
	*x = SigningAnnouncementMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_tbtc_gen_pb_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SigningAnnouncementMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SigningAnnouncementMessage) ProtoMessage() {}

func (x *SigningAnnouncementMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_tbtc_gen_pb_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SigningAnnouncementMessage.ProtoReflect.Descriptor instead.
func (*SigningAnnouncementMessage) Descriptor() ([]byte, []int) {
	return file_pkg_tbtc_gen_pb_message_proto_rawDescGZIP(), []int{1}
}

func (x *SigningAnnouncementMessage) GetSenderID() uint32 {
	if x != nil {
		return x.SenderID
	}
	return 0
}

func (x *SigningAnnouncementMessage) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *SigningAnnouncementMessage) GetAttemptNumber() uint64 {
	if x != nil {
		return x.AttemptNumber
	}
	return 0
}

var File_pkg_tbtc_gen_pb_message_proto protoreflect.FileDescriptor

var file_pkg_tbtc_gen_pb_message_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x62, 0x74, 0x63, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70,
	0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x74, 0x62, 0x74, 0x63, 0x22, 0x30, 0x0a, 0x08, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x69, 0x6c,
	0x6c, 0x12, 0x24, 0x0a, 0x0d, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70,
	0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x78, 0x0a, 0x1a, 0x53, 0x69, 0x67, 0x6e, 0x69,
	0x6e, 0x67, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49,
	0x44, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x61,
	0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0d, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_pkg_tbtc_gen_pb_message_proto_rawDescOnce sync.Once
	file_pkg_tbtc_gen_pb_message_proto_rawDescData = file_pkg_tbtc_gen_pb_message_proto_rawDesc
)

func file_pkg_tbtc_gen_pb_message_proto_rawDescGZIP() []byte {
	file_pkg_tbtc_gen_pb_message_proto_rawDescOnce.Do(func() {
		file_pkg_tbtc_gen_pb_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_tbtc_gen_pb_message_proto_rawDescData)
	})
	return file_pkg_tbtc_gen_pb_message_proto_rawDescData
}

var file_pkg_tbtc_gen_pb_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_tbtc_gen_pb_message_proto_goTypes = []interface{}{
	(*StopPill)(nil),                   // 0: tbtc.StopPill
	(*SigningAnnouncementMessage)(nil), // 1: tbtc.SigningAnnouncementMessage
}
var file_pkg_tbtc_gen_pb_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_tbtc_gen_pb_message_proto_init() }
func file_pkg_tbtc_gen_pb_message_proto_init() {
	if File_pkg_tbtc_gen_pb_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_tbtc_gen_pb_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopPill); i {
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
		file_pkg_tbtc_gen_pb_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SigningAnnouncementMessage); i {
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
			RawDescriptor: file_pkg_tbtc_gen_pb_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_tbtc_gen_pb_message_proto_goTypes,
		DependencyIndexes: file_pkg_tbtc_gen_pb_message_proto_depIdxs,
		MessageInfos:      file_pkg_tbtc_gen_pb_message_proto_msgTypes,
	}.Build()
	File_pkg_tbtc_gen_pb_message_proto = out.File
	file_pkg_tbtc_gen_pb_message_proto_rawDesc = nil
	file_pkg_tbtc_gen_pb_message_proto_goTypes = nil
	file_pkg_tbtc_gen_pb_message_proto_depIdxs = nil
}

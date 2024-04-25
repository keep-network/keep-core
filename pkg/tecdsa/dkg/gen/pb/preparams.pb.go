// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.7.1
// source: pkg/tecdsa/dkg/gen/pb/preparams.proto

package pb

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type PreParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data              *PreParams_LocalPreParams `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	CreationTimestamp *timestamp.Timestamp      `protobuf:"bytes,2,opt,name=creationTimestamp,proto3" json:"creationTimestamp,omitempty"`
}

func (x *PreParams) Reset() {
	*x = PreParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreParams) ProtoMessage() {}

func (x *PreParams) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreParams.ProtoReflect.Descriptor instead.
func (*PreParams) Descriptor() ([]byte, []int) {
	return file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescGZIP(), []int{0}
}

func (x *PreParams) GetData() *PreParams_LocalPreParams {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PreParams) GetCreationTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.CreationTimestamp
	}
	return nil
}

type PreParams_PublicKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	N []byte `protobuf:"bytes,1,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *PreParams_PublicKey) Reset() {
	*x = PreParams_PublicKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreParams_PublicKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreParams_PublicKey) ProtoMessage() {}

func (x *PreParams_PublicKey) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreParams_PublicKey.ProtoReflect.Descriptor instead.
func (*PreParams_PublicKey) Descriptor() ([]byte, []int) {
	return file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescGZIP(), []int{0, 0}
}

func (x *PreParams_PublicKey) GetN() []byte {
	if x != nil {
		return x.N
	}
	return nil
}

type PreParams_PrivateKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublicKey *PreParams_PublicKey `protobuf:"bytes,1,opt,name=publicKey,proto3" json:"publicKey,omitempty"`
	LambdaN   []byte               `protobuf:"bytes,2,opt,name=lambdaN,proto3" json:"lambdaN,omitempty"`
	PhiN      []byte               `protobuf:"bytes,3,opt,name=phiN,proto3" json:"phiN,omitempty"`
}

func (x *PreParams_PrivateKey) Reset() {
	*x = PreParams_PrivateKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreParams_PrivateKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreParams_PrivateKey) ProtoMessage() {}

func (x *PreParams_PrivateKey) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreParams_PrivateKey.ProtoReflect.Descriptor instead.
func (*PreParams_PrivateKey) Descriptor() ([]byte, []int) {
	return file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescGZIP(), []int{0, 1}
}

func (x *PreParams_PrivateKey) GetPublicKey() *PreParams_PublicKey {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *PreParams_PrivateKey) GetLambdaN() []byte {
	if x != nil {
		return x.LambdaN
	}
	return nil
}

func (x *PreParams_PrivateKey) GetPhiN() []byte {
	if x != nil {
		return x.PhiN
	}
	return nil
}

type PreParams_LocalPreParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PaillierSK *PreParams_PrivateKey `protobuf:"bytes,1,opt,name=paillierSK,proto3" json:"paillierSK,omitempty"`
	NTilde     []byte                `protobuf:"bytes,2,opt,name=nTilde,proto3" json:"nTilde,omitempty"`
	H1I        []byte                `protobuf:"bytes,3,opt,name=h1i,proto3" json:"h1i,omitempty"`
	H2I        []byte                `protobuf:"bytes,4,opt,name=h2i,proto3" json:"h2i,omitempty"`
	Alpha      []byte                `protobuf:"bytes,5,opt,name=alpha,proto3" json:"alpha,omitempty"`
	Beta       []byte                `protobuf:"bytes,6,opt,name=beta,proto3" json:"beta,omitempty"`
	P          []byte                `protobuf:"bytes,7,opt,name=p,proto3" json:"p,omitempty"`
	Q          []byte                `protobuf:"bytes,8,opt,name=q,proto3" json:"q,omitempty"`
}

func (x *PreParams_LocalPreParams) Reset() {
	*x = PreParams_LocalPreParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreParams_LocalPreParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreParams_LocalPreParams) ProtoMessage() {}

func (x *PreParams_LocalPreParams) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreParams_LocalPreParams.ProtoReflect.Descriptor instead.
func (*PreParams_LocalPreParams) Descriptor() ([]byte, []int) {
	return file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescGZIP(), []int{0, 2}
}

func (x *PreParams_LocalPreParams) GetPaillierSK() *PreParams_PrivateKey {
	if x != nil {
		return x.PaillierSK
	}
	return nil
}

func (x *PreParams_LocalPreParams) GetNTilde() []byte {
	if x != nil {
		return x.NTilde
	}
	return nil
}

func (x *PreParams_LocalPreParams) GetH1I() []byte {
	if x != nil {
		return x.H1I
	}
	return nil
}

func (x *PreParams_LocalPreParams) GetH2I() []byte {
	if x != nil {
		return x.H2I
	}
	return nil
}

func (x *PreParams_LocalPreParams) GetAlpha() []byte {
	if x != nil {
		return x.Alpha
	}
	return nil
}

func (x *PreParams_LocalPreParams) GetBeta() []byte {
	if x != nil {
		return x.Beta
	}
	return nil
}

func (x *PreParams_LocalPreParams) GetP() []byte {
	if x != nil {
		return x.P
	}
	return nil
}

func (x *PreParams_LocalPreParams) GetQ() []byte {
	if x != nil {
		return x.Q
	}
	return nil
}

var File_pkg_tecdsa_dkg_gen_pb_preparams_proto protoreflect.FileDescriptor

var file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDesc = []byte{
	0x0a, 0x25, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x65, 0x63, 0x64, 0x73, 0x61, 0x2f, 0x64, 0x6b, 0x67,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x64, 0x6b, 0x67, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe7, 0x03,
	0x0a, 0x09, 0x50, 0x72, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x31, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x64, 0x6b, 0x67, 0x2e,
	0x50, 0x72, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x50,
	0x72, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x48,
	0x0a, 0x11, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x11, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1a, 0x19, 0x0a, 0x09, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x0c, 0x0a, 0x01, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x01, 0x6e, 0x1a, 0x72, 0x0a, 0x0a, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65,
	0x79, 0x12, 0x36, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x6b, 0x67, 0x2e, 0x50, 0x72, 0x65, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x09,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x61, 0x6d,
	0x62, 0x64, 0x61, 0x4e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6c, 0x61, 0x6d, 0x62,
	0x64, 0x61, 0x4e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x68, 0x69, 0x4e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x70, 0x68, 0x69, 0x4e, 0x1a, 0xcd, 0x01, 0x0a, 0x0e, 0x4c, 0x6f, 0x63, 0x61,
	0x6c, 0x50, 0x72, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x70, 0x61,
	0x69, 0x6c, 0x6c, 0x69, 0x65, 0x72, 0x53, 0x4b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x64, 0x6b, 0x67, 0x2e, 0x50, 0x72, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x2e, 0x50,
	0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x52, 0x0a, 0x70, 0x61, 0x69, 0x6c, 0x6c,
	0x69, 0x65, 0x72, 0x53, 0x4b, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x54, 0x69, 0x6c, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x6e, 0x54, 0x69, 0x6c, 0x64, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x68, 0x31, 0x69, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x68, 0x31, 0x69, 0x12,
	0x10, 0x0a, 0x03, 0x68, 0x32, 0x69, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x68, 0x32,
	0x69, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x65, 0x74, 0x61, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x65, 0x74, 0x61, 0x12, 0x0c, 0x0a, 0x01, 0x70,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x70, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x71, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescOnce sync.Once
	file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescData = file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDesc
)

func file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescGZIP() []byte {
	file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescOnce.Do(func() {
		file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescData)
	})
	return file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDescData
}

var file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_tecdsa_dkg_gen_pb_preparams_proto_goTypes = []interface{}{
	(*PreParams)(nil),                // 0: dkg.PreParams
	(*PreParams_PublicKey)(nil),      // 1: dkg.PreParams.PublicKey
	(*PreParams_PrivateKey)(nil),     // 2: dkg.PreParams.PrivateKey
	(*PreParams_LocalPreParams)(nil), // 3: dkg.PreParams.LocalPreParams
	(*timestamp.Timestamp)(nil),      // 4: google.protobuf.Timestamp
}
var file_pkg_tecdsa_dkg_gen_pb_preparams_proto_depIdxs = []int32{
	3, // 0: dkg.PreParams.data:type_name -> dkg.PreParams.LocalPreParams
	4, // 1: dkg.PreParams.creationTimestamp:type_name -> google.protobuf.Timestamp
	1, // 2: dkg.PreParams.PrivateKey.publicKey:type_name -> dkg.PreParams.PublicKey
	2, // 3: dkg.PreParams.LocalPreParams.paillierSK:type_name -> dkg.PreParams.PrivateKey
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pkg_tecdsa_dkg_gen_pb_preparams_proto_init() }
func file_pkg_tecdsa_dkg_gen_pb_preparams_proto_init() {
	if File_pkg_tecdsa_dkg_gen_pb_preparams_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreParams); i {
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
		file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreParams_PublicKey); i {
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
		file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreParams_PrivateKey); i {
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
		file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreParams_LocalPreParams); i {
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
			RawDescriptor: file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_tecdsa_dkg_gen_pb_preparams_proto_goTypes,
		DependencyIndexes: file_pkg_tecdsa_dkg_gen_pb_preparams_proto_depIdxs,
		MessageInfos:      file_pkg_tecdsa_dkg_gen_pb_preparams_proto_msgTypes,
	}.Build()
	File_pkg_tecdsa_dkg_gen_pb_preparams_proto = out.File
	file_pkg_tecdsa_dkg_gen_pb_preparams_proto_rawDesc = nil
	file_pkg_tecdsa_dkg_gen_pb_preparams_proto_goTypes = nil
	file_pkg_tecdsa_dkg_gen_pb_preparams_proto_depIdxs = nil
}

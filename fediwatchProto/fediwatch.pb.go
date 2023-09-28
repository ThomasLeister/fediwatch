// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.6.1
// source: fediwatch.proto

package fediwatchProto

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

type Connection_Direction int32

const (
	Connection_UNKNOWN Connection_Direction = 0
	Connection_IN      Connection_Direction = 1
	Connection_OUT     Connection_Direction = 2
)

// Enum value maps for Connection_Direction.
var (
	Connection_Direction_name = map[int32]string{
		0: "UNKNOWN",
		1: "IN",
		2: "OUT",
	}
	Connection_Direction_value = map[string]int32{
		"UNKNOWN": 0,
		"IN":      1,
		"OUT":     2,
	}
)

func (x Connection_Direction) Enum() *Connection_Direction {
	p := new(Connection_Direction)
	*p = x
	return p
}

func (x Connection_Direction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Connection_Direction) Descriptor() protoreflect.EnumDescriptor {
	return file_fediwatch_proto_enumTypes[0].Descriptor()
}

func (Connection_Direction) Type() protoreflect.EnumType {
	return &file_fediwatch_proto_enumTypes[0]
}

func (x Connection_Direction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Connection_Direction.Descriptor instead.
func (Connection_Direction) EnumDescriptor() ([]byte, []int) {
	return file_fediwatch_proto_rawDescGZIP(), []int{0, 0}
}

type Connection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dir Connection_Direction `protobuf:"varint,1,opt,name=dir,proto3,enum=fediwatch.Connection_Direction" json:"dir,omitempty"`
	Lat float32              `protobuf:"fixed32,2,opt,name=lat,proto3" json:"lat,omitempty"`
	Lng float32              `protobuf:"fixed32,3,opt,name=lng,proto3" json:"lng,omitempty"`
}

func (x *Connection) Reset() {
	*x = Connection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fediwatch_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Connection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Connection) ProtoMessage() {}

func (x *Connection) ProtoReflect() protoreflect.Message {
	mi := &file_fediwatch_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Connection.ProtoReflect.Descriptor instead.
func (*Connection) Descriptor() ([]byte, []int) {
	return file_fediwatch_proto_rawDescGZIP(), []int{0}
}

func (x *Connection) GetDir() Connection_Direction {
	if x != nil {
		return x.Dir
	}
	return Connection_UNKNOWN
}

func (x *Connection) GetLat() float32 {
	if x != nil {
		return x.Lat
	}
	return 0
}

func (x *Connection) GetLng() float32 {
	if x != nil {
		return x.Lng
	}
	return 0
}

var File_fediwatch_proto protoreflect.FileDescriptor

var file_fediwatch_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x66, 0x65, 0x64, 0x69, 0x77, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x66, 0x65, 0x64, 0x69, 0x77, 0x61, 0x74, 0x63, 0x68, 0x22, 0x8e, 0x01, 0x0a,
	0x0a, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x31, 0x0a, 0x03, 0x64,
	0x69, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x66, 0x65, 0x64, 0x69, 0x77,
	0x61, 0x74, 0x63, 0x68, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x64, 0x69, 0x72, 0x12, 0x10,
	0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x6c, 0x61, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x6c,
	0x6e, 0x67, 0x22, 0x29, 0x0a, 0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02,
	0x49, 0x4e, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x4f, 0x55, 0x54, 0x10, 0x02, 0x42, 0x12, 0x5a,
	0x10, 0x2e, 0x2f, 0x66, 0x65, 0x64, 0x69, 0x77, 0x61, 0x74, 0x63, 0x68, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fediwatch_proto_rawDescOnce sync.Once
	file_fediwatch_proto_rawDescData = file_fediwatch_proto_rawDesc
)

func file_fediwatch_proto_rawDescGZIP() []byte {
	file_fediwatch_proto_rawDescOnce.Do(func() {
		file_fediwatch_proto_rawDescData = protoimpl.X.CompressGZIP(file_fediwatch_proto_rawDescData)
	})
	return file_fediwatch_proto_rawDescData
}

var file_fediwatch_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_fediwatch_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_fediwatch_proto_goTypes = []interface{}{
	(Connection_Direction)(0), // 0: fediwatch.Connection.Direction
	(*Connection)(nil),        // 1: fediwatch.Connection
}
var file_fediwatch_proto_depIdxs = []int32{
	0, // 0: fediwatch.Connection.dir:type_name -> fediwatch.Connection.Direction
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_fediwatch_proto_init() }
func file_fediwatch_proto_init() {
	if File_fediwatch_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fediwatch_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Connection); i {
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
			RawDescriptor: file_fediwatch_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fediwatch_proto_goTypes,
		DependencyIndexes: file_fediwatch_proto_depIdxs,
		EnumInfos:         file_fediwatch_proto_enumTypes,
		MessageInfos:      file_fediwatch_proto_msgTypes,
	}.Build()
	File_fediwatch_proto = out.File
	file_fediwatch_proto_rawDesc = nil
	file_fediwatch_proto_goTypes = nil
	file_fediwatch_proto_depIdxs = nil
}
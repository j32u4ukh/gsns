// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: Edge.proto

package pbgo

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

type Edge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// {"default": "AI", "primary_key": "default"}
	Index  int32 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	UserId int32 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Target int32 `protobuf:"varint,3,opt,name=target,proto3" json:"target,omitempty"`
	// {"ignore": "true"}
	Targets []int32 `protobuf:"varint,4,rep,packed,name=targets,proto3" json:"targets,omitempty"`
	// {"default": "current_timestamp()"}
	CreateTime *TimeStamp `protobuf:"bytes,5,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// {"default": "current_timestamp()", "update": "current_timestamp()"}
	UpdateTime *TimeStamp `protobuf:"bytes,6,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	// {"ignore": "true"}
	CreateUtc int64 `protobuf:"varint,7,opt,name=create_utc,json=createUtc,proto3" json:"create_utc,omitempty"`
	// {"ignore": "true"}
	UpdateUtc int64 `protobuf:"varint,8,opt,name=update_utc,json=updateUtc,proto3" json:"update_utc,omitempty"`
}

func (x *Edge) Reset() {
	*x = Edge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Edge_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Edge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Edge) ProtoMessage() {}

func (x *Edge) ProtoReflect() protoreflect.Message {
	mi := &file_Edge_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Edge.ProtoReflect.Descriptor instead.
func (*Edge) Descriptor() ([]byte, []int) {
	return file_Edge_proto_rawDescGZIP(), []int{0}
}

func (x *Edge) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Edge) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Edge) GetTarget() int32 {
	if x != nil {
		return x.Target
	}
	return 0
}

func (x *Edge) GetTargets() []int32 {
	if x != nil {
		return x.Targets
	}
	return nil
}

func (x *Edge) GetCreateTime() *TimeStamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *Edge) GetUpdateTime() *TimeStamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

func (x *Edge) GetCreateUtc() int64 {
	if x != nil {
		return x.CreateUtc
	}
	return 0
}

func (x *Edge) GetUpdateUtc() int64 {
	if x != nil {
		return x.UpdateUtc
	}
	return 0
}

var File_Edge_proto protoreflect.FileDescriptor

var file_Edge_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x45, 0x64, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x54, 0x69,
	0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xff, 0x01,
	0x0a, 0x04, 0x45, 0x64, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x05, 0x52, 0x07,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x53, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x74, 0x63, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x74, 0x63,
	0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x74, 0x63, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x74, 0x63, 0x42,
	0x08, 0x5a, 0x06, 0x2e, 0x3b, 0x70, 0x62, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_Edge_proto_rawDescOnce sync.Once
	file_Edge_proto_rawDescData = file_Edge_proto_rawDesc
)

func file_Edge_proto_rawDescGZIP() []byte {
	file_Edge_proto_rawDescOnce.Do(func() {
		file_Edge_proto_rawDescData = protoimpl.X.CompressGZIP(file_Edge_proto_rawDescData)
	})
	return file_Edge_proto_rawDescData
}

var file_Edge_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_Edge_proto_goTypes = []interface{}{
	(*Edge)(nil),      // 0: Edge
	(*TimeStamp)(nil), // 1: TimeStamp
}
var file_Edge_proto_depIdxs = []int32{
	1, // 0: Edge.create_time:type_name -> TimeStamp
	1, // 1: Edge.update_time:type_name -> TimeStamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_Edge_proto_init() }
func file_Edge_proto_init() {
	if File_Edge_proto != nil {
		return
	}
	file_TimeStamp_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_Edge_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Edge); i {
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
			RawDescriptor: file_Edge_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_Edge_proto_goTypes,
		DependencyIndexes: file_Edge_proto_depIdxs,
		MessageInfos:      file_Edge_proto_msgTypes,
	}.Build()
	File_Edge_proto = out.File
	file_Edge_proto_rawDesc = nil
	file_Edge_proto_goTypes = nil
	file_Edge_proto_depIdxs = nil
}

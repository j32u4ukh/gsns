// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.3
// source: AccountArray.proto

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

type AccountArray struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Accounts []*Account `protobuf:"bytes,1,rep,name=accounts,proto3" json:"accounts,omitempty"`
}

func (x *AccountArray) Reset() {
	*x = AccountArray{}
	if protoimpl.UnsafeEnabled {
		mi := &file_AccountArray_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountArray) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountArray) ProtoMessage() {}

func (x *AccountArray) ProtoReflect() protoreflect.Message {
	mi := &file_AccountArray_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountArray.ProtoReflect.Descriptor instead.
func (*AccountArray) Descriptor() ([]byte, []int) {
	return file_AccountArray_proto_rawDescGZIP(), []int{0}
}

func (x *AccountArray) GetAccounts() []*Account {
	if x != nil {
		return x.Accounts
	}
	return nil
}

var File_AccountArray_proto protoreflect.FileDescriptor

var file_AccountArray_proto_rawDesc = []byte{
	0x0a, 0x12, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x72, 0x72, 0x61, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x34, 0x0a, 0x0c, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x72,
	0x72, 0x61, 0x79, 0x12, 0x24, 0x0a, 0x08, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52,
	0x08, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x3b, 0x70,
	0x62, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_AccountArray_proto_rawDescOnce sync.Once
	file_AccountArray_proto_rawDescData = file_AccountArray_proto_rawDesc
)

func file_AccountArray_proto_rawDescGZIP() []byte {
	file_AccountArray_proto_rawDescOnce.Do(func() {
		file_AccountArray_proto_rawDescData = protoimpl.X.CompressGZIP(file_AccountArray_proto_rawDescData)
	})
	return file_AccountArray_proto_rawDescData
}

var file_AccountArray_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_AccountArray_proto_goTypes = []interface{}{
	(*AccountArray)(nil), // 0: AccountArray
	(*Account)(nil),      // 1: Account
}
var file_AccountArray_proto_depIdxs = []int32{
	1, // 0: AccountArray.accounts:type_name -> Account
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_AccountArray_proto_init() }
func file_AccountArray_proto_init() {
	if File_AccountArray_proto != nil {
		return
	}
	file_Account_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_AccountArray_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountArray); i {
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
			RawDescriptor: file_AccountArray_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_AccountArray_proto_goTypes,
		DependencyIndexes: file_AccountArray_proto_depIdxs,
		MessageInfos:      file_AccountArray_proto_msgTypes,
	}.Build()
	File_AccountArray_proto = out.File
	file_AccountArray_proto_rawDesc = nil
	file_AccountArray_proto_goTypes = nil
	file_AccountArray_proto_depIdxs = nil
}

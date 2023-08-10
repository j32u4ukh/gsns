// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: Agreement.proto

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

type Agreement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd          int32          `protobuf:"varint,1,opt,name=cmd,proto3" json:"cmd,omitempty"`
	Service      int32          `protobuf:"varint,2,opt,name=service,proto3" json:"service,omitempty"`
	ReturnCode   int32          `protobuf:"varint,3,opt,name=return_code,json=returnCode,proto3" json:"return_code,omitempty"`
	Msg          string         `protobuf:"bytes,4,opt,name=msg,proto3" json:"msg,omitempty"`
	Cid          int32          `protobuf:"varint,5,opt,name=cid,proto3" json:"cid,omitempty"`
	Accounts     []*Account     `protobuf:"bytes,6,rep,name=accounts,proto3" json:"accounts,omitempty"`
	Users        []*User        `protobuf:"bytes,7,rep,name=users,proto3" json:"users,omitempty"`
	PostMessages []*PostMessage `protobuf:"bytes,8,rep,name=post_messages,json=postMessages,proto3" json:"post_messages,omitempty"`
	Cipher       string         `protobuf:"bytes,9,opt,name=cipher,proto3" json:"cipher,omitempty"`
	Identity     int32          `protobuf:"varint,10,opt,name=identity,proto3" json:"identity,omitempty"`
	Edges        []*Edge        `protobuf:"bytes,11,rep,name=edges,proto3" json:"edges,omitempty"`
	StartUtc     int64          `protobuf:"varint,12,opt,name=start_utc,json=startUtc,proto3" json:"start_utc,omitempty"`
	StopUtc      int64          `protobuf:"varint,13,opt,name=stop_utc,json=stopUtc,proto3" json:"stop_utc,omitempty"`
}

func (x *Agreement) Reset() {
	*x = Agreement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Agreement_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Agreement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Agreement) ProtoMessage() {}

func (x *Agreement) ProtoReflect() protoreflect.Message {
	mi := &file_Agreement_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Agreement.ProtoReflect.Descriptor instead.
func (*Agreement) Descriptor() ([]byte, []int) {
	return file_Agreement_proto_rawDescGZIP(), []int{0}
}

func (x *Agreement) GetCmd() int32 {
	if x != nil {
		return x.Cmd
	}
	return 0
}

func (x *Agreement) GetService() int32 {
	if x != nil {
		return x.Service
	}
	return 0
}

func (x *Agreement) GetReturnCode() int32 {
	if x != nil {
		return x.ReturnCode
	}
	return 0
}

func (x *Agreement) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Agreement) GetCid() int32 {
	if x != nil {
		return x.Cid
	}
	return 0
}

func (x *Agreement) GetAccounts() []*Account {
	if x != nil {
		return x.Accounts
	}
	return nil
}

func (x *Agreement) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *Agreement) GetPostMessages() []*PostMessage {
	if x != nil {
		return x.PostMessages
	}
	return nil
}

func (x *Agreement) GetCipher() string {
	if x != nil {
		return x.Cipher
	}
	return ""
}

func (x *Agreement) GetIdentity() int32 {
	if x != nil {
		return x.Identity
	}
	return 0
}

func (x *Agreement) GetEdges() []*Edge {
	if x != nil {
		return x.Edges
	}
	return nil
}

func (x *Agreement) GetStartUtc() int64 {
	if x != nil {
		return x.StartUtc
	}
	return 0
}

func (x *Agreement) GetStopUtc() int64 {
	if x != nil {
		return x.StopUtc
	}
	return 0
}

var File_Agreement_proto protoreflect.FileDescriptor

var file_Agreement_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x41, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x50, 0x6f,
	0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0a, 0x45, 0x64, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfb, 0x02, 0x0a, 0x09,
	0x41, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x6d, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x5f,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x72, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x69, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x63, 0x69, 0x64, 0x12, 0x24, 0x0a, 0x08, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x08, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73,
	0x12, 0x1b, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x31, 0x0a,
	0x0d, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x08,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x0c, 0x70, 0x6f, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x12, 0x1b, 0x0a, 0x05, 0x65, 0x64, 0x67, 0x65, 0x73, 0x18, 0x0b, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x45, 0x64, 0x67, 0x65, 0x52, 0x05, 0x65, 0x64, 0x67, 0x65,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x75, 0x74, 0x63, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x74, 0x61, 0x72, 0x74, 0x55, 0x74, 0x63, 0x12, 0x19,
	0x0a, 0x08, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x75, 0x74, 0x63, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x73, 0x74, 0x6f, 0x70, 0x55, 0x74, 0x63, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x3b, 0x70,
	0x62, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_Agreement_proto_rawDescOnce sync.Once
	file_Agreement_proto_rawDescData = file_Agreement_proto_rawDesc
)

func file_Agreement_proto_rawDescGZIP() []byte {
	file_Agreement_proto_rawDescOnce.Do(func() {
		file_Agreement_proto_rawDescData = protoimpl.X.CompressGZIP(file_Agreement_proto_rawDescData)
	})
	return file_Agreement_proto_rawDescData
}

var file_Agreement_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_Agreement_proto_goTypes = []interface{}{
	(*Agreement)(nil),   // 0: Agreement
	(*Account)(nil),     // 1: Account
	(*User)(nil),        // 2: User
	(*PostMessage)(nil), // 3: PostMessage
	(*Edge)(nil),        // 4: Edge
}
var file_Agreement_proto_depIdxs = []int32{
	1, // 0: Agreement.accounts:type_name -> Account
	2, // 1: Agreement.users:type_name -> User
	3, // 2: Agreement.post_messages:type_name -> PostMessage
	4, // 3: Agreement.edges:type_name -> Edge
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_Agreement_proto_init() }
func file_Agreement_proto_init() {
	if File_Agreement_proto != nil {
		return
	}
	file_Account_proto_init()
	file_User_proto_init()
	file_PostMessage_proto_init()
	file_Edge_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_Agreement_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Agreement); i {
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
			RawDescriptor: file_Agreement_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_Agreement_proto_goTypes,
		DependencyIndexes: file_Agreement_proto_depIdxs,
		MessageInfos:      file_Agreement_proto_msgTypes,
	}.Build()
	File_Agreement_proto = out.File
	file_Agreement_proto_rawDesc = nil
	file_Agreement_proto_goTypes = nil
	file_Agreement_proto_depIdxs = nil
}

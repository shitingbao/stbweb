// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.11.4
// source: stbserver.proto

package stbserver

import (
	context "context"
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Character struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Xaxis    int64       `protobuf:"varint,1,opt,name=xaxis,proto3" json:"xaxis,omitempty"`
	Yaxis    int64       `protobuf:"varint,2,opt,name=yaxis,proto3" json:"yaxis,omitempty"`
	Zaxis    int64       `protobuf:"varint,3,opt,name=zaxis,proto3" json:"zaxis,omitempty"`
	Area     string      `protobuf:"bytes,4,opt,name=area,proto3" json:"area,omitempty"`
	Name     string      `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Skill    []*Skill    `protobuf:"bytes,6,rep,name=skill,proto3" json:"skill,omitempty"`
	Summoner []*Summoner `protobuf:"bytes,7,rep,name=summoner,proto3" json:"summoner,omitempty"`
}

func (x *Character) Reset() {
	*x = Character{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stbserver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Character) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Character) ProtoMessage() {}

func (x *Character) ProtoReflect() protoreflect.Message {
	mi := &file_stbserver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Character.ProtoReflect.Descriptor instead.
func (*Character) Descriptor() ([]byte, []int) {
	return file_stbserver_proto_rawDescGZIP(), []int{0}
}

func (x *Character) GetXaxis() int64 {
	if x != nil {
		return x.Xaxis
	}
	return 0
}

func (x *Character) GetYaxis() int64 {
	if x != nil {
		return x.Yaxis
	}
	return 0
}

func (x *Character) GetZaxis() int64 {
	if x != nil {
		return x.Zaxis
	}
	return 0
}

func (x *Character) GetArea() string {
	if x != nil {
		return x.Area
	}
	return ""
}

func (x *Character) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Character) GetSkill() []*Skill {
	if x != nil {
		return x.Skill
	}
	return nil
}

func (x *Character) GetSummoner() []*Summoner {
	if x != nil {
		return x.Summoner
	}
	return nil
}

type Skill struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ordinary float64 `protobuf:"fixed64,1,opt,name=ordinary,proto3" json:"ordinary,omitempty"`
	Qkill    string  `protobuf:"bytes,2,opt,name=qkill,proto3" json:"qkill,omitempty"`
	Wkill    string  `protobuf:"bytes,3,opt,name=wkill,proto3" json:"wkill,omitempty"`
	Ekill    string  `protobuf:"bytes,4,opt,name=ekill,proto3" json:"ekill,omitempty"`
	Rkill    string  `protobuf:"bytes,5,opt,name=rkill,proto3" json:"rkill,omitempty"`
}

func (x *Skill) Reset() {
	*x = Skill{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stbserver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Skill) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Skill) ProtoMessage() {}

func (x *Skill) ProtoReflect() protoreflect.Message {
	mi := &file_stbserver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Skill.ProtoReflect.Descriptor instead.
func (*Skill) Descriptor() ([]byte, []int) {
	return file_stbserver_proto_rawDescGZIP(), []int{1}
}

func (x *Skill) GetOrdinary() float64 {
	if x != nil {
		return x.Ordinary
	}
	return 0
}

func (x *Skill) GetQkill() string {
	if x != nil {
		return x.Qkill
	}
	return ""
}

func (x *Skill) GetWkill() string {
	if x != nil {
		return x.Wkill
	}
	return ""
}

func (x *Skill) GetEkill() string {
	if x != nil {
		return x.Ekill
	}
	return ""
}

func (x *Skill) GetRkill() string {
	if x != nil {
		return x.Rkill
	}
	return ""
}

type Summoner struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dkill string `protobuf:"bytes,1,opt,name=dkill,proto3" json:"dkill,omitempty"`
	Fkill string `protobuf:"bytes,2,opt,name=fkill,proto3" json:"fkill,omitempty"`
}

func (x *Summoner) Reset() {
	*x = Summoner{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stbserver_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Summoner) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Summoner) ProtoMessage() {}

func (x *Summoner) ProtoReflect() protoreflect.Message {
	mi := &file_stbserver_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Summoner.ProtoReflect.Descriptor instead.
func (*Summoner) Descriptor() ([]byte, []int) {
	return file_stbserver_proto_rawDescGZIP(), []int{2}
}

func (x *Summoner) GetDkill() string {
	if x != nil {
		return x.Dkill
	}
	return ""
}

func (x *Summoner) GetFkill() string {
	if x != nil {
		return x.Fkill
	}
	return ""
}

type Identity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Idcard string `protobuf:"bytes,1,opt,name=idcard,proto3" json:"idcard,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Identity) Reset() {
	*x = Identity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stbserver_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Identity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Identity) ProtoMessage() {}

func (x *Identity) ProtoReflect() protoreflect.Message {
	mi := &file_stbserver_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Identity.ProtoReflect.Descriptor instead.
func (*Identity) Descriptor() ([]byte, []int) {
	return file_stbserver_proto_rawDescGZIP(), []int{3}
}

func (x *Identity) GetIdcard() string {
	if x != nil {
		return x.Idcard
	}
	return ""
}

func (x *Identity) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_stbserver_proto protoreflect.FileDescriptor

var file_stbserver_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x73, 0x74, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x73, 0x74, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0xce, 0x01, 0x0a,
	0x09, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x78, 0x61,
	0x78, 0x69, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x78, 0x61, 0x78, 0x69, 0x73,
	0x12, 0x14, 0x0a, 0x05, 0x79, 0x61, 0x78, 0x69, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x79, 0x61, 0x78, 0x69, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x7a, 0x61, 0x78, 0x69, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x7a, 0x61, 0x78, 0x69, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x61, 0x72, 0x65, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72, 0x65, 0x61,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x73, 0x6b, 0x69, 0x6c, 0x6c, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x74, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x53, 0x6b, 0x69, 0x6c, 0x6c, 0x52, 0x05, 0x73, 0x6b, 0x69, 0x6c, 0x6c, 0x12, 0x2f, 0x0a, 0x08,
	0x73, 0x75, 0x6d, 0x6d, 0x6f, 0x6e, 0x65, 0x72, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x73, 0x74, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x6f,
	0x6e, 0x65, 0x72, 0x52, 0x08, 0x73, 0x75, 0x6d, 0x6d, 0x6f, 0x6e, 0x65, 0x72, 0x22, 0x7b, 0x0a,
	0x05, 0x53, 0x6b, 0x69, 0x6c, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61,
	0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61,
	0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x6b, 0x69, 0x6c, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x71, 0x6b, 0x69, 0x6c, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x6b, 0x69, 0x6c,
	0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x77, 0x6b, 0x69, 0x6c, 0x6c, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x6b, 0x69, 0x6c, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x6b, 0x69, 0x6c, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6b, 0x69, 0x6c, 0x6c, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x6b, 0x69, 0x6c, 0x6c, 0x22, 0x36, 0x0a, 0x08, 0x53, 0x75,
	0x6d, 0x6d, 0x6f, 0x6e, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x6b, 0x69, 0x6c, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x6b, 0x69, 0x6c, 0x6c, 0x12, 0x14, 0x0a, 0x05,
	0x66, 0x6b, 0x69, 0x6c, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x6b, 0x69,
	0x6c, 0x6c, 0x22, 0x36, 0x0a, 0x08, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x69, 0x64, 0x63, 0x61, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x69, 0x64, 0x63, 0x61, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0x98, 0x02, 0x0a, 0x09, 0x53,
	0x74, 0x62, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x3e, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53,
	0x75, 0x6d, 0x6d, 0x6f, 0x6e, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x13, 0x2e, 0x73, 0x74,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x1a, 0x14, 0x2e, 0x73, 0x74, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x61,
	0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x0f, 0x50, 0x75, 0x74, 0x53,
	0x75, 0x6d, 0x6d, 0x6f, 0x6e, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x13, 0x2e, 0x73, 0x74,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x1a, 0x14, 0x2e, 0x73, 0x74, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x61,
	0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x22, 0x00, 0x28, 0x01, 0x12, 0x43, 0x0a, 0x12, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x53, 0x75, 0x6d, 0x6d, 0x6f, 0x6e, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x13, 0x2e, 0x73, 0x74, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x49, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x73, 0x74, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x22, 0x00, 0x30, 0x01, 0x12,
	0x44, 0x0a, 0x11, 0x53, 0x68, 0x61, 0x72, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x6f, 0x6e, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x13, 0x2e, 0x73, 0x74, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x73, 0x74, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x22,
	0x00, 0x28, 0x01, 0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stbserver_proto_rawDescOnce sync.Once
	file_stbserver_proto_rawDescData = file_stbserver_proto_rawDesc
)

func file_stbserver_proto_rawDescGZIP() []byte {
	file_stbserver_proto_rawDescOnce.Do(func() {
		file_stbserver_proto_rawDescData = protoimpl.X.CompressGZIP(file_stbserver_proto_rawDescData)
	})
	return file_stbserver_proto_rawDescData
}

var file_stbserver_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_stbserver_proto_goTypes = []interface{}{
	(*Character)(nil), // 0: stbserver.Character
	(*Skill)(nil),     // 1: stbserver.Skill
	(*Summoner)(nil),  // 2: stbserver.Summoner
	(*Identity)(nil),  // 3: stbserver.Identity
}
var file_stbserver_proto_depIdxs = []int32{
	1, // 0: stbserver.Character.skill:type_name -> stbserver.Skill
	2, // 1: stbserver.Character.summoner:type_name -> stbserver.Summoner
	3, // 2: stbserver.StbServer.GetSummonerInfo:input_type -> stbserver.Identity
	3, // 3: stbserver.StbServer.PutSummonerInfo:input_type -> stbserver.Identity
	3, // 4: stbserver.StbServer.GetAllSummonerInfo:input_type -> stbserver.Identity
	3, // 5: stbserver.StbServer.ShareSummonerInfo:input_type -> stbserver.Identity
	0, // 6: stbserver.StbServer.GetSummonerInfo:output_type -> stbserver.Character
	0, // 7: stbserver.StbServer.PutSummonerInfo:output_type -> stbserver.Character
	0, // 8: stbserver.StbServer.GetAllSummonerInfo:output_type -> stbserver.Character
	0, // 9: stbserver.StbServer.ShareSummonerInfo:output_type -> stbserver.Character
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_stbserver_proto_init() }
func file_stbserver_proto_init() {
	if File_stbserver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_stbserver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Character); i {
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
		file_stbserver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Skill); i {
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
		file_stbserver_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Summoner); i {
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
		file_stbserver_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Identity); i {
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
			RawDescriptor: file_stbserver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stbserver_proto_goTypes,
		DependencyIndexes: file_stbserver_proto_depIdxs,
		MessageInfos:      file_stbserver_proto_msgTypes,
	}.Build()
	File_stbserver_proto = out.File
	file_stbserver_proto_rawDesc = nil
	file_stbserver_proto_goTypes = nil
	file_stbserver_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StbServerClient is the client API for StbServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StbServerClient interface {
	//rpc ServerTest()returns(){}不能使用参数或者返回值为空的服务
	GetSummonerInfo(ctx context.Context, in *Identity, opts ...grpc.CallOption) (*Character, error)
	PutSummonerInfo(ctx context.Context, opts ...grpc.CallOption) (StbServer_PutSummonerInfoClient, error)
	GetAllSummonerInfo(ctx context.Context, in *Identity, opts ...grpc.CallOption) (StbServer_GetAllSummonerInfoClient, error)
	ShareSummonerInfo(ctx context.Context, opts ...grpc.CallOption) (StbServer_ShareSummonerInfoClient, error)
}

type stbServerClient struct {
	cc grpc.ClientConnInterface
}

func NewStbServerClient(cc grpc.ClientConnInterface) StbServerClient {
	return &stbServerClient{cc}
}

func (c *stbServerClient) GetSummonerInfo(ctx context.Context, in *Identity, opts ...grpc.CallOption) (*Character, error) {
	out := new(Character)
	err := c.cc.Invoke(ctx, "/stbserver.StbServer/GetSummonerInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stbServerClient) PutSummonerInfo(ctx context.Context, opts ...grpc.CallOption) (StbServer_PutSummonerInfoClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StbServer_serviceDesc.Streams[0], "/stbserver.StbServer/PutSummonerInfo", opts...)
	if err != nil {
		return nil, err
	}
	x := &stbServerPutSummonerInfoClient{stream}
	return x, nil
}

type StbServer_PutSummonerInfoClient interface {
	Send(*Identity) error
	CloseAndRecv() (*Character, error)
	grpc.ClientStream
}

type stbServerPutSummonerInfoClient struct {
	grpc.ClientStream
}

func (x *stbServerPutSummonerInfoClient) Send(m *Identity) error {
	return x.ClientStream.SendMsg(m)
}

func (x *stbServerPutSummonerInfoClient) CloseAndRecv() (*Character, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Character)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *stbServerClient) GetAllSummonerInfo(ctx context.Context, in *Identity, opts ...grpc.CallOption) (StbServer_GetAllSummonerInfoClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StbServer_serviceDesc.Streams[1], "/stbserver.StbServer/GetAllSummonerInfo", opts...)
	if err != nil {
		return nil, err
	}
	x := &stbServerGetAllSummonerInfoClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StbServer_GetAllSummonerInfoClient interface {
	Recv() (*Character, error)
	grpc.ClientStream
}

type stbServerGetAllSummonerInfoClient struct {
	grpc.ClientStream
}

func (x *stbServerGetAllSummonerInfoClient) Recv() (*Character, error) {
	m := new(Character)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *stbServerClient) ShareSummonerInfo(ctx context.Context, opts ...grpc.CallOption) (StbServer_ShareSummonerInfoClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StbServer_serviceDesc.Streams[2], "/stbserver.StbServer/ShareSummonerInfo", opts...)
	if err != nil {
		return nil, err
	}
	x := &stbServerShareSummonerInfoClient{stream}
	return x, nil
}

type StbServer_ShareSummonerInfoClient interface {
	Send(*Identity) error
	Recv() (*Character, error)
	grpc.ClientStream
}

type stbServerShareSummonerInfoClient struct {
	grpc.ClientStream
}

func (x *stbServerShareSummonerInfoClient) Send(m *Identity) error {
	return x.ClientStream.SendMsg(m)
}

func (x *stbServerShareSummonerInfoClient) Recv() (*Character, error) {
	m := new(Character)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StbServerServer is the server API for StbServer service.
type StbServerServer interface {
	//rpc ServerTest()returns(){}不能使用参数或者返回值为空的服务
	GetSummonerInfo(context.Context, *Identity) (*Character, error)
	PutSummonerInfo(StbServer_PutSummonerInfoServer) error
	GetAllSummonerInfo(*Identity, StbServer_GetAllSummonerInfoServer) error
	ShareSummonerInfo(StbServer_ShareSummonerInfoServer) error
}

// UnimplementedStbServerServer can be embedded to have forward compatible implementations.
type UnimplementedStbServerServer struct {
}

func (*UnimplementedStbServerServer) GetSummonerInfo(context.Context, *Identity) (*Character, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSummonerInfo not implemented")
}
func (*UnimplementedStbServerServer) PutSummonerInfo(StbServer_PutSummonerInfoServer) error {
	return status.Errorf(codes.Unimplemented, "method PutSummonerInfo not implemented")
}
func (*UnimplementedStbServerServer) GetAllSummonerInfo(*Identity, StbServer_GetAllSummonerInfoServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAllSummonerInfo not implemented")
}
func (*UnimplementedStbServerServer) ShareSummonerInfo(StbServer_ShareSummonerInfoServer) error {
	return status.Errorf(codes.Unimplemented, "method ShareSummonerInfo not implemented")
}

func RegisterStbServerServer(s *grpc.Server, srv StbServerServer) {
	s.RegisterService(&_StbServer_serviceDesc, srv)
}

func _StbServer_GetSummonerInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Identity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StbServerServer).GetSummonerInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stbserver.StbServer/GetSummonerInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StbServerServer).GetSummonerInfo(ctx, req.(*Identity))
	}
	return interceptor(ctx, in, info, handler)
}

func _StbServer_PutSummonerInfo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StbServerServer).PutSummonerInfo(&stbServerPutSummonerInfoServer{stream})
}

type StbServer_PutSummonerInfoServer interface {
	SendAndClose(*Character) error
	Recv() (*Identity, error)
	grpc.ServerStream
}

type stbServerPutSummonerInfoServer struct {
	grpc.ServerStream
}

func (x *stbServerPutSummonerInfoServer) SendAndClose(m *Character) error {
	return x.ServerStream.SendMsg(m)
}

func (x *stbServerPutSummonerInfoServer) Recv() (*Identity, error) {
	m := new(Identity)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _StbServer_GetAllSummonerInfo_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Identity)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StbServerServer).GetAllSummonerInfo(m, &stbServerGetAllSummonerInfoServer{stream})
}

type StbServer_GetAllSummonerInfoServer interface {
	Send(*Character) error
	grpc.ServerStream
}

type stbServerGetAllSummonerInfoServer struct {
	grpc.ServerStream
}

func (x *stbServerGetAllSummonerInfoServer) Send(m *Character) error {
	return x.ServerStream.SendMsg(m)
}

func _StbServer_ShareSummonerInfo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StbServerServer).ShareSummonerInfo(&stbServerShareSummonerInfoServer{stream})
}

type StbServer_ShareSummonerInfoServer interface {
	Send(*Character) error
	Recv() (*Identity, error)
	grpc.ServerStream
}

type stbServerShareSummonerInfoServer struct {
	grpc.ServerStream
}

func (x *stbServerShareSummonerInfoServer) Send(m *Character) error {
	return x.ServerStream.SendMsg(m)
}

func (x *stbServerShareSummonerInfoServer) Recv() (*Identity, error) {
	m := new(Identity)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _StbServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "stbserver.StbServer",
	HandlerType: (*StbServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSummonerInfo",
			Handler:    _StbServer_GetSummonerInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PutSummonerInfo",
			Handler:       _StbServer_PutSummonerInfo_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetAllSummonerInfo",
			Handler:       _StbServer_GetAllSummonerInfo_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ShareSummonerInfo",
			Handler:       _StbServer_ShareSummonerInfo_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "stbserver.proto",
}

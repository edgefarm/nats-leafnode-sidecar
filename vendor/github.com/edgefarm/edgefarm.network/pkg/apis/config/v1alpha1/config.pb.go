// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

package v1alpha1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type DesiredStateRequest struct {
	AccountName          string   `protobuf:"bytes,1,opt,name=accountName,proto3" json:"accountName,omitempty"`
	Username             []string `protobuf:"bytes,2,rep,name=username,proto3" json:"username,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DesiredStateRequest) Reset()         { *m = DesiredStateRequest{} }
func (m *DesiredStateRequest) String() string { return proto.CompactTextString(m) }
func (*DesiredStateRequest) ProtoMessage()    {}
func (*DesiredStateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{0}
}

func (m *DesiredStateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DesiredStateRequest.Unmarshal(m, b)
}
func (m *DesiredStateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DesiredStateRequest.Marshal(b, m, deterministic)
}
func (m *DesiredStateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DesiredStateRequest.Merge(m, src)
}
func (m *DesiredStateRequest) XXX_Size() int {
	return xxx_messageInfo_DesiredStateRequest.Size(m)
}
func (m *DesiredStateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DesiredStateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DesiredStateRequest proto.InternalMessageInfo

func (m *DesiredStateRequest) GetAccountName() string {
	if m != nil {
		return m.AccountName
	}
	return ""
}

func (m *DesiredStateRequest) GetUsername() []string {
	if m != nil {
		return m.Username
	}
	return nil
}

type Credentials struct {
	UserAccountName      string   `protobuf:"bytes,1,opt,name=userAccountName,proto3" json:"userAccountName,omitempty"`
	Username             string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	Creds                string   `protobuf:"bytes,4,opt,name=creds,proto3" json:"creds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Credentials) Reset()         { *m = Credentials{} }
func (m *Credentials) String() string { return proto.CompactTextString(m) }
func (*Credentials) ProtoMessage()    {}
func (*Credentials) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{1}
}

func (m *Credentials) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Credentials.Unmarshal(m, b)
}
func (m *Credentials) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Credentials.Marshal(b, m, deterministic)
}
func (m *Credentials) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Credentials.Merge(m, src)
}
func (m *Credentials) XXX_Size() int {
	return xxx_messageInfo_Credentials.Size(m)
}
func (m *Credentials) XXX_DiscardUnknown() {
	xxx_messageInfo_Credentials.DiscardUnknown(m)
}

var xxx_messageInfo_Credentials proto.InternalMessageInfo

func (m *Credentials) GetUserAccountName() string {
	if m != nil {
		return m.UserAccountName
	}
	return ""
}

func (m *Credentials) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Credentials) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *Credentials) GetCreds() string {
	if m != nil {
		return m.Creds
	}
	return ""
}

type DesiredStateResponse struct {
	Creds                []*Credentials `protobuf:"bytes,1,rep,name=creds,proto3" json:"creds,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *DesiredStateResponse) Reset()         { *m = DesiredStateResponse{} }
func (m *DesiredStateResponse) String() string { return proto.CompactTextString(m) }
func (*DesiredStateResponse) ProtoMessage()    {}
func (*DesiredStateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{2}
}

func (m *DesiredStateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DesiredStateResponse.Unmarshal(m, b)
}
func (m *DesiredStateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DesiredStateResponse.Marshal(b, m, deterministic)
}
func (m *DesiredStateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DesiredStateResponse.Merge(m, src)
}
func (m *DesiredStateResponse) XXX_Size() int {
	return xxx_messageInfo_DesiredStateResponse.Size(m)
}
func (m *DesiredStateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DesiredStateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DesiredStateResponse proto.InternalMessageInfo

func (m *DesiredStateResponse) GetCreds() []*Credentials {
	if m != nil {
		return m.Creds
	}
	return nil
}

type DeleteAccountRequest struct {
	AccountName          string   `protobuf:"bytes,1,opt,name=accountName,proto3" json:"accountName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteAccountRequest) Reset()         { *m = DeleteAccountRequest{} }
func (m *DeleteAccountRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteAccountRequest) ProtoMessage()    {}
func (*DeleteAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{3}
}

func (m *DeleteAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteAccountRequest.Unmarshal(m, b)
}
func (m *DeleteAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteAccountRequest.Marshal(b, m, deterministic)
}
func (m *DeleteAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteAccountRequest.Merge(m, src)
}
func (m *DeleteAccountRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteAccountRequest.Size(m)
}
func (m *DeleteAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteAccountRequest proto.InternalMessageInfo

func (m *DeleteAccountRequest) GetAccountName() string {
	if m != nil {
		return m.AccountName
	}
	return ""
}

type DeleteAccountResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteAccountResponse) Reset()         { *m = DeleteAccountResponse{} }
func (m *DeleteAccountResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteAccountResponse) ProtoMessage()    {}
func (*DeleteAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{4}
}

func (m *DeleteAccountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteAccountResponse.Unmarshal(m, b)
}
func (m *DeleteAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteAccountResponse.Marshal(b, m, deterministic)
}
func (m *DeleteAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteAccountResponse.Merge(m, src)
}
func (m *DeleteAccountResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteAccountResponse.Size(m)
}
func (m *DeleteAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteAccountResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*DesiredStateRequest)(nil), "v1alpha1.DesiredStateRequest")
	proto.RegisterType((*Credentials)(nil), "v1alpha1.Credentials")
	proto.RegisterType((*DesiredStateResponse)(nil), "v1alpha1.DesiredStateResponse")
	proto.RegisterType((*DeleteAccountRequest)(nil), "v1alpha1.DeleteAccountRequest")
	proto.RegisterType((*DeleteAccountResponse)(nil), "v1alpha1.DeleteAccountResponse")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor_3eaf2c85e69e9ea4) }

var fileDescriptor_3eaf2c85e69e9ea4 = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x4a, 0xf3, 0x40,
	0x10, 0xc7, 0xd9, 0x2f, 0x9f, 0xd2, 0x4e, 0x52, 0x84, 0xb5, 0xc5, 0x10, 0x50, 0x43, 0x4e, 0x01,
	0x21, 0xd0, 0x7a, 0x11, 0x3c, 0x69, 0xbc, 0x2a, 0x92, 0xdc, 0xbc, 0xad, 0xc9, 0xa8, 0x81, 0xb8,
	0x1b, 0x77, 0x37, 0xf5, 0x09, 0x7c, 0x1f, 0x1f, 0x51, 0x92, 0x74, 0x35, 0x49, 0x2d, 0x78, 0x9c,
	0xfd, 0xcd, 0xfc, 0xe7, 0x3f, 0x33, 0x0b, 0x4e, 0x26, 0xf8, 0x53, 0xf1, 0x1c, 0x55, 0x52, 0x68,
	0x41, 0x27, 0xeb, 0x25, 0x2b, 0xab, 0x17, 0xb6, 0x0c, 0x52, 0x38, 0xbc, 0x41, 0x55, 0x48, 0xcc,
	0x53, 0xcd, 0x34, 0x26, 0xf8, 0x56, 0xa3, 0xd2, 0xd4, 0x07, 0x9b, 0x65, 0x99, 0xa8, 0xb9, 0xbe,
	0x63, 0xaf, 0xe8, 0x12, 0x9f, 0x84, 0xd3, 0xa4, 0xff, 0x44, 0x3d, 0x98, 0xd4, 0x0a, 0x25, 0x6f,
	0xf0, 0x3f, 0xdf, 0x0a, 0xa7, 0xc9, 0x77, 0x1c, 0x7c, 0x10, 0xb0, 0x63, 0x89, 0x39, 0x72, 0x5d,
	0xb0, 0x52, 0xd1, 0x10, 0x0e, 0x1a, 0x76, 0xb5, 0xa5, 0x38, 0x7e, 0x1e, 0xa9, 0x92, 0xbe, 0x6a,
	0xc3, 0x2a, 0xa6, 0xd4, 0xbb, 0x90, 0xb9, 0x6b, 0x75, 0xcc, 0xc4, 0x74, 0x0e, 0x7b, 0x99, 0xc4,
	0x5c, 0xb9, 0xff, 0x5b, 0xd0, 0x05, 0x41, 0x0c, 0xf3, 0xe1, 0x70, 0xaa, 0x12, 0x5c, 0x21, 0x3d,
	0x33, 0xd9, 0xc4, 0xb7, 0x42, 0x7b, 0xb5, 0x88, 0xcc, 0x3a, 0xa2, 0x9e, 0x6b, 0x23, 0x72, 0xd1,
	0x88, 0x94, 0xa8, 0x71, 0xe3, 0xf3, 0xcf, 0x2b, 0x0a, 0x8e, 0x60, 0x31, 0xaa, 0xec, 0xfa, 0xaf,
	0x3e, 0x09, 0xcc, 0xe2, 0xf6, 0x1e, 0x29, 0xca, 0x75, 0x91, 0x21, 0xbd, 0x05, 0xa7, 0xef, 0x94,
	0x1e, 0xff, 0x58, 0xfa, 0xe5, 0x3c, 0xde, 0xc9, 0x2e, 0xbc, 0x19, 0xf0, 0x1e, 0x66, 0x83, 0xce,
	0x74, 0x50, 0xb0, 0x3d, 0x8c, 0x77, 0xba, 0x93, 0x77, 0x8a, 0xd7, 0xce, 0x03, 0x44, 0x97, 0x26,
	0xe7, 0x71, 0xbf, 0xfd, 0x46, 0xe7, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x20, 0xfd, 0xb5, 0x0b,
	0x56, 0x02, 0x00, 0x00,
}

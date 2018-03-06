// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

package vcs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type User struct {
	Id    int64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Login string `protobuf:"bytes,2,opt,name=login" json:"login,omitempty"`
	Url   string `protobuf:"bytes,5,opt,name=url" json:"url,omitempty"`
	Type  string `protobuf:"bytes,21,opt,name=type" json:"type,omitempty"`
	Name  string `protobuf:"bytes,23,opt,name=name" json:"name,omitempty"`
	Email string `protobuf:"bytes,27,opt,name=email" json:"email,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *User) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *User) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *User) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func init() {
	proto.RegisterType((*User)(nil), "vcs.User")
}

func init() { proto.RegisterFile("user.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x2d, 0x4e, 0x2d,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x4b, 0x2e, 0x56, 0xaa, 0xe1, 0x62, 0x09,
	0x2d, 0x4e, 0x2d, 0x12, 0xe2, 0xe3, 0x62, 0xca, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x0e,
	0x62, 0xca, 0x4c, 0x11, 0x12, 0xe1, 0x62, 0xcd, 0xc9, 0x4f, 0xcf, 0xcc, 0x93, 0x60, 0x52, 0x60,
	0xd4, 0xe0, 0x0c, 0x82, 0x70, 0x84, 0x04, 0xb8, 0x98, 0x4b, 0x8b, 0x72, 0x24, 0x58, 0xc1, 0x62,
	0x20, 0xa6, 0x90, 0x10, 0x17, 0x4b, 0x49, 0x65, 0x41, 0xaa, 0x84, 0x28, 0x58, 0x08, 0xcc, 0x06,
	0x89, 0xe5, 0x25, 0xe6, 0xa6, 0x4a, 0x88, 0x43, 0xc4, 0x40, 0x6c, 0x90, 0x79, 0xa9, 0xb9, 0x89,
	0x99, 0x39, 0x12, 0xd2, 0x10, 0xf3, 0xc0, 0x9c, 0x24, 0x36, 0xb0, 0x4b, 0x8c, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x29, 0x29, 0x53, 0x96, 0x97, 0x00, 0x00, 0x00,
}
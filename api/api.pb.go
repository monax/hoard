// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grant "github.com/monax/hoard/v5/grant"
	reference "github.com/monax/hoard/v5/reference"
	stores "github.com/monax/hoard/v5/stores"
	grpc "google.golang.org/grpc"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type GrantAndGrantSpec struct {
	Grant *grant.Grant `protobuf:"bytes,1,opt,name=Grant,proto3" json:"Grant,omitempty"`
	// The type of grant to output
	GrantSpec            *grant.Spec `protobuf:"bytes,2,opt,name=GrantSpec,proto3" json:"GrantSpec,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GrantAndGrantSpec) Reset()         { *m = GrantAndGrantSpec{} }
func (m *GrantAndGrantSpec) String() string { return proto.CompactTextString(m) }
func (*GrantAndGrantSpec) ProtoMessage()    {}
func (*GrantAndGrantSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}
func (m *GrantAndGrantSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GrantAndGrantSpec.Unmarshal(m, b)
}
func (m *GrantAndGrantSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GrantAndGrantSpec.Marshal(b, m, deterministic)
}
func (m *GrantAndGrantSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrantAndGrantSpec.Merge(m, src)
}
func (m *GrantAndGrantSpec) XXX_Size() int {
	return xxx_messageInfo_GrantAndGrantSpec.Size(m)
}
func (m *GrantAndGrantSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_GrantAndGrantSpec.DiscardUnknown(m)
}

var xxx_messageInfo_GrantAndGrantSpec proto.InternalMessageInfo

func (m *GrantAndGrantSpec) GetGrant() *grant.Grant {
	if m != nil {
		return m.Grant
	}
	return nil
}

func (m *GrantAndGrantSpec) GetGrantSpec() *grant.Spec {
	if m != nil {
		return m.GrantSpec
	}
	return nil
}

type PlaintextAndGrantSpec struct {
	Plaintext *Plaintext `protobuf:"bytes,1,opt,name=Plaintext,proto3" json:"Plaintext,omitempty"`
	// The type of grant to output
	GrantSpec            *grant.Spec `protobuf:"bytes,2,opt,name=GrantSpec,proto3" json:"GrantSpec,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *PlaintextAndGrantSpec) Reset()         { *m = PlaintextAndGrantSpec{} }
func (m *PlaintextAndGrantSpec) String() string { return proto.CompactTextString(m) }
func (*PlaintextAndGrantSpec) ProtoMessage()    {}
func (*PlaintextAndGrantSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}
func (m *PlaintextAndGrantSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlaintextAndGrantSpec.Unmarshal(m, b)
}
func (m *PlaintextAndGrantSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlaintextAndGrantSpec.Marshal(b, m, deterministic)
}
func (m *PlaintextAndGrantSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlaintextAndGrantSpec.Merge(m, src)
}
func (m *PlaintextAndGrantSpec) XXX_Size() int {
	return xxx_messageInfo_PlaintextAndGrantSpec.Size(m)
}
func (m *PlaintextAndGrantSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_PlaintextAndGrantSpec.DiscardUnknown(m)
}

var xxx_messageInfo_PlaintextAndGrantSpec proto.InternalMessageInfo

func (m *PlaintextAndGrantSpec) GetPlaintext() *Plaintext {
	if m != nil {
		return m.Plaintext
	}
	return nil
}

func (m *PlaintextAndGrantSpec) GetGrantSpec() *grant.Spec {
	if m != nil {
		return m.GrantSpec
	}
	return nil
}

type ReferenceAndGrantSpec struct {
	Reference *reference.Ref `protobuf:"bytes,1,opt,name=Reference,proto3" json:"Reference,omitempty"`
	// The type of grant to output
	GrantSpec            *grant.Spec `protobuf:"bytes,2,opt,name=GrantSpec,proto3" json:"GrantSpec,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ReferenceAndGrantSpec) Reset()         { *m = ReferenceAndGrantSpec{} }
func (m *ReferenceAndGrantSpec) String() string { return proto.CompactTextString(m) }
func (*ReferenceAndGrantSpec) ProtoMessage()    {}
func (*ReferenceAndGrantSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}
func (m *ReferenceAndGrantSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReferenceAndGrantSpec.Unmarshal(m, b)
}
func (m *ReferenceAndGrantSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReferenceAndGrantSpec.Marshal(b, m, deterministic)
}
func (m *ReferenceAndGrantSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReferenceAndGrantSpec.Merge(m, src)
}
func (m *ReferenceAndGrantSpec) XXX_Size() int {
	return xxx_messageInfo_ReferenceAndGrantSpec.Size(m)
}
func (m *ReferenceAndGrantSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_ReferenceAndGrantSpec.DiscardUnknown(m)
}

var xxx_messageInfo_ReferenceAndGrantSpec proto.InternalMessageInfo

func (m *ReferenceAndGrantSpec) GetReference() *reference.Ref {
	if m != nil {
		return m.Reference
	}
	return nil
}

func (m *ReferenceAndGrantSpec) GetGrantSpec() *grant.Spec {
	if m != nil {
		return m.GrantSpec
	}
	return nil
}

type Plaintext struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Salt                 []byte   `protobuf:"bytes,2,opt,name=Salt,proto3" json:"Salt,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Plaintext) Reset()         { *m = Plaintext{} }
func (m *Plaintext) String() string { return proto.CompactTextString(m) }
func (*Plaintext) ProtoMessage()    {}
func (*Plaintext) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}
func (m *Plaintext) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Plaintext.Unmarshal(m, b)
}
func (m *Plaintext) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Plaintext.Marshal(b, m, deterministic)
}
func (m *Plaintext) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Plaintext.Merge(m, src)
}
func (m *Plaintext) XXX_Size() int {
	return xxx_messageInfo_Plaintext.Size(m)
}
func (m *Plaintext) XXX_DiscardUnknown() {
	xxx_messageInfo_Plaintext.DiscardUnknown(m)
}

var xxx_messageInfo_Plaintext proto.InternalMessageInfo

func (m *Plaintext) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Plaintext) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

type Ciphertext struct {
	EncryptedData        []byte   `protobuf:"bytes,1,opt,name=EncryptedData,proto3" json:"EncryptedData,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ciphertext) Reset()         { *m = Ciphertext{} }
func (m *Ciphertext) String() string { return proto.CompactTextString(m) }
func (*Ciphertext) ProtoMessage()    {}
func (*Ciphertext) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}
func (m *Ciphertext) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ciphertext.Unmarshal(m, b)
}
func (m *Ciphertext) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ciphertext.Marshal(b, m, deterministic)
}
func (m *Ciphertext) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ciphertext.Merge(m, src)
}
func (m *Ciphertext) XXX_Size() int {
	return xxx_messageInfo_Ciphertext.Size(m)
}
func (m *Ciphertext) XXX_DiscardUnknown() {
	xxx_messageInfo_Ciphertext.DiscardUnknown(m)
}

var xxx_messageInfo_Ciphertext proto.InternalMessageInfo

func (m *Ciphertext) GetEncryptedData() []byte {
	if m != nil {
		return m.EncryptedData
	}
	return nil
}

type ReferenceAndCiphertext struct {
	Reference            *reference.Ref `protobuf:"bytes,1,opt,name=Reference,proto3" json:"Reference,omitempty"`
	Ciphertext           *Ciphertext    `protobuf:"bytes,2,opt,name=Ciphertext,proto3" json:"Ciphertext,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ReferenceAndCiphertext) Reset()         { *m = ReferenceAndCiphertext{} }
func (m *ReferenceAndCiphertext) String() string { return proto.CompactTextString(m) }
func (*ReferenceAndCiphertext) ProtoMessage()    {}
func (*ReferenceAndCiphertext) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}
func (m *ReferenceAndCiphertext) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReferenceAndCiphertext.Unmarshal(m, b)
}
func (m *ReferenceAndCiphertext) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReferenceAndCiphertext.Marshal(b, m, deterministic)
}
func (m *ReferenceAndCiphertext) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReferenceAndCiphertext.Merge(m, src)
}
func (m *ReferenceAndCiphertext) XXX_Size() int {
	return xxx_messageInfo_ReferenceAndCiphertext.Size(m)
}
func (m *ReferenceAndCiphertext) XXX_DiscardUnknown() {
	xxx_messageInfo_ReferenceAndCiphertext.DiscardUnknown(m)
}

var xxx_messageInfo_ReferenceAndCiphertext proto.InternalMessageInfo

func (m *ReferenceAndCiphertext) GetReference() *reference.Ref {
	if m != nil {
		return m.Reference
	}
	return nil
}

func (m *ReferenceAndCiphertext) GetCiphertext() *Ciphertext {
	if m != nil {
		return m.Ciphertext
	}
	return nil
}

type Address struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()    {}
func (*Address) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{6}
}
func (m *Address) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Address.Unmarshal(m, b)
}
func (m *Address) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Address.Marshal(b, m, deterministic)
}
func (m *Address) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Address.Merge(m, src)
}
func (m *Address) XXX_Size() int {
	return xxx_messageInfo_Address.Size(m)
}
func (m *Address) XXX_DiscardUnknown() {
	xxx_messageInfo_Address.DiscardUnknown(m)
}

var xxx_messageInfo_Address proto.InternalMessageInfo

func (m *Address) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func init() {
	proto.RegisterType((*GrantAndGrantSpec)(nil), "api.GrantAndGrantSpec")
	proto.RegisterType((*PlaintextAndGrantSpec)(nil), "api.PlaintextAndGrantSpec")
	proto.RegisterType((*ReferenceAndGrantSpec)(nil), "api.ReferenceAndGrantSpec")
	proto.RegisterType((*Plaintext)(nil), "api.Plaintext")
	proto.RegisterType((*Ciphertext)(nil), "api.Ciphertext")
	proto.RegisterType((*ReferenceAndCiphertext)(nil), "api.ReferenceAndCiphertext")
	proto.RegisterType((*Address)(nil), "api.Address")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 486 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x51, 0x8b, 0xd3, 0x40,
	0x10, 0xa6, 0xb6, 0x36, 0x64, 0x1a, 0x3d, 0x5d, 0xf0, 0x28, 0x11, 0x51, 0xa2, 0x3d, 0x3c, 0x90,
	0x44, 0x5a, 0xb9, 0xf7, 0xf3, 0x4e, 0x0e, 0xdf, 0xca, 0x06, 0x1f, 0xf4, 0x6d, 0xdb, 0x4c, 0xdb,
	0x40, 0x2e, 0x1b, 0x36, 0x5b, 0x3d, 0x41, 0xf0, 0x37, 0xfb, 0x0f, 0x24, 0x9b, 0x6d, 0xb2, 0xd9,
	0x1c, 0x07, 0x7d, 0xca, 0xce, 0x37, 0xdf, 0x37, 0xdf, 0xce, 0xce, 0x10, 0x70, 0x59, 0x91, 0x86,
	0x85, 0xe0, 0x92, 0x93, 0x21, 0x2b, 0x52, 0xff, 0x44, 0xe0, 0x06, 0x05, 0xe6, 0x6b, 0xac, 0x51,
	0x7f, 0xb2, 0x15, 0x2c, 0x97, 0x3a, 0xf0, 0x4a, 0xc9, 0x05, 0x96, 0x75, 0x14, 0xac, 0xe0, 0xf9,
	0x4d, 0x95, 0xbc, 0xcc, 0x13, 0xf5, 0x8d, 0x0b, 0x5c, 0x93, 0x00, 0x1e, 0xab, 0x60, 0x3a, 0x78,
	0x33, 0x78, 0x3f, 0x99, 0x7b, 0x61, 0xad, 0x57, 0x18, 0xad, 0x53, 0xe4, 0x1c, 0xdc, 0x46, 0x30,
	0x7d, 0xa4, 0x78, 0x13, 0xcd, 0xab, 0x20, 0xda, 0x66, 0x83, 0x02, 0x5e, 0x2c, 0x33, 0x96, 0xe6,
	0x12, 0xef, 0xba, 0x3e, 0x1f, 0xc0, 0x6d, 0x12, 0xda, 0xeb, 0x69, 0x58, 0x35, 0xd3, 0xa0, 0xb4,
	0x25, 0x1c, 0xe9, 0x48, 0x0f, 0x6f, 0x60, 0x3b, 0x36, 0x89, 0xc6, 0xb1, 0x7d, 0x2e, 0x8a, 0x1b,
	0xda, 0x12, 0x8e, 0x71, 0x5c, 0x18, 0xad, 0x10, 0x02, 0xa3, 0x6b, 0x26, 0x99, 0x32, 0xf0, 0xa8,
	0x3a, 0x57, 0x58, 0xcc, 0x32, 0xa9, 0xca, 0x78, 0x54, 0x9d, 0x83, 0x39, 0xc0, 0x55, 0x5a, 0xec,
	0x50, 0x28, 0xd5, 0x3b, 0x78, 0xf2, 0x25, 0x5f, 0x8b, 0xdf, 0x85, 0xc4, 0xc4, 0x90, 0x77, 0xc1,
	0xe0, 0x17, 0x9c, 0x9a, 0xad, 0x19, 0xfa, 0xe3, 0x7a, 0x8b, 0x4c, 0x6f, 0xdd, 0xdc, 0x89, 0x7a,
	0xfc, 0x16, 0xa6, 0x06, 0x25, 0x78, 0x0b, 0xce, 0x65, 0x92, 0x08, 0x2c, 0x4b, 0x32, 0x6d, 0x8e,
	0xfa, 0x8e, 0x87, 0x70, 0xfe, 0x6f, 0xa0, 0x57, 0x87, 0x7c, 0x84, 0x51, 0x8c, 0x2c, 0x23, 0xbe,
	0xaa, 0x79, 0xef, 0x34, 0xfc, 0xce, 0x62, 0x91, 0x33, 0x18, 0x7f, 0xcb, 0xcb, 0x4a, 0xd3, 0xc1,
	0x7d, 0xab, 0x09, 0x12, 0xc2, 0x98, 0xa2, 0xe2, 0x9d, 0xaa, 0xda, 0xbd, 0xfd, 0xb5, 0xea, 0x2e,
	0xc0, 0x59, 0xee, 0xa5, 0x71, 0x99, 0x7b, 0x97, 0xd1, 0x12, 0x9d, 0x83, 0x5b, 0x5f, 0xe6, 0x06,
	0x65, 0xef, 0x3e, 0x9d, 0x22, 0xf3, 0xef, 0xe0, 0x5e, 0x65, 0xc8, 0xea, 0x21, 0xcc, 0x60, 0xb8,
	0xdc, 0x4b, 0x62, 0x71, 0x7a, 0x3d, 0xcc, 0x60, 0x58, 0x15, 0xb6, 0xe0, 0x5e, 0xe9, 0x3f, 0x00,
	0x7a, 0xfa, 0x29, 0xcf, 0xc9, 0x05, 0x38, 0x3a, 0xea, 0xd5, 0x7f, 0xd9, 0x7b, 0x65, 0x63, 0x31,
	0x2e, 0xc0, 0xb9, 0xc6, 0x5a, 0xf7, 0x10, 0xaf, 0xe7, 0xfe, 0x17, 0x9c, 0x58, 0x72, 0xc1, 0xb6,
	0x48, 0x66, 0x30, 0x5a, 0xee, 0xcb, 0x1d, 0xb1, 0x37, 0xc4, 0xf7, 0x14, 0x70, 0x58, 0x0c, 0x45,
	0xcb, 0xaa, 0x01, 0x1a, 0xa8, 0x6f, 0x8b, 0xc8, 0x19, 0x8c, 0x62, 0xc9, 0xa4, 0x45, 0x7b, 0x16,
	0xea, 0x3f, 0x53, 0x95, 0xfb, 0x9a, 0x6f, 0xf8, 0xe7, 0xd7, 0x3f, 0x5e, 0x6d, 0x53, 0xb9, 0xdb,
	0xaf, 0xc2, 0x35, 0xbf, 0x8d, 0x6e, 0x79, 0xce, 0xee, 0xa2, 0x1d, 0x67, 0x22, 0x89, 0x7e, 0x7e,
	0x8a, 0x58, 0x91, 0xae, 0xc6, 0xea, 0x27, 0xb6, 0xf8, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x0f,
	0x13, 0x6f, 0x02, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GrantClient is the client API for Grant service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GrantClient interface {
	// Seal a Reference to create a Grant
	Seal(ctx context.Context, in *ReferenceAndGrantSpec, opts ...grpc.CallOption) (*grant.Grant, error)
	// Unseal a Grant to recover the Reference
	Unseal(ctx context.Context, in *grant.Grant, opts ...grpc.CallOption) (*reference.Ref, error)
	// Convert one grant to another grant to re-share with another party or just
	// to change grant type
	Reseal(ctx context.Context, in *GrantAndGrantSpec, opts ...grpc.CallOption) (*grant.Grant, error)
	// Put a Plaintext and returned the sealed Reference as a Grant
	PutSeal(ctx context.Context, in *PlaintextAndGrantSpec, opts ...grpc.CallOption) (*grant.Grant, error)
	// Unseal a Grant and follow the Reference to return a Plaintext
	UnsealGet(ctx context.Context, in *grant.Grant, opts ...grpc.CallOption) (*Plaintext, error)
}

type grantClient struct {
	cc *grpc.ClientConn
}

func NewGrantClient(cc *grpc.ClientConn) GrantClient {
	return &grantClient{cc}
}

func (c *grantClient) Seal(ctx context.Context, in *ReferenceAndGrantSpec, opts ...grpc.CallOption) (*grant.Grant, error) {
	out := new(grant.Grant)
	err := c.cc.Invoke(ctx, "/api.Grant/Seal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grantClient) Unseal(ctx context.Context, in *grant.Grant, opts ...grpc.CallOption) (*reference.Ref, error) {
	out := new(reference.Ref)
	err := c.cc.Invoke(ctx, "/api.Grant/Unseal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grantClient) Reseal(ctx context.Context, in *GrantAndGrantSpec, opts ...grpc.CallOption) (*grant.Grant, error) {
	out := new(grant.Grant)
	err := c.cc.Invoke(ctx, "/api.Grant/Reseal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grantClient) PutSeal(ctx context.Context, in *PlaintextAndGrantSpec, opts ...grpc.CallOption) (*grant.Grant, error) {
	out := new(grant.Grant)
	err := c.cc.Invoke(ctx, "/api.Grant/PutSeal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grantClient) UnsealGet(ctx context.Context, in *grant.Grant, opts ...grpc.CallOption) (*Plaintext, error) {
	out := new(Plaintext)
	err := c.cc.Invoke(ctx, "/api.Grant/UnsealGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrantServer is the server API for Grant service.
type GrantServer interface {
	// Seal a Reference to create a Grant
	Seal(context.Context, *ReferenceAndGrantSpec) (*grant.Grant, error)
	// Unseal a Grant to recover the Reference
	Unseal(context.Context, *grant.Grant) (*reference.Ref, error)
	// Convert one grant to another grant to re-share with another party or just
	// to change grant type
	Reseal(context.Context, *GrantAndGrantSpec) (*grant.Grant, error)
	// Put a Plaintext and returned the sealed Reference as a Grant
	PutSeal(context.Context, *PlaintextAndGrantSpec) (*grant.Grant, error)
	// Unseal a Grant and follow the Reference to return a Plaintext
	UnsealGet(context.Context, *grant.Grant) (*Plaintext, error)
}

func RegisterGrantServer(s *grpc.Server, srv GrantServer) {
	s.RegisterService(&_Grant_serviceDesc, srv)
}

func _Grant_Seal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReferenceAndGrantSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrantServer).Seal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Grant/Seal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrantServer).Seal(ctx, req.(*ReferenceAndGrantSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grant_Unseal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(grant.Grant)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrantServer).Unseal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Grant/Unseal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrantServer).Unseal(ctx, req.(*grant.Grant))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grant_Reseal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrantAndGrantSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrantServer).Reseal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Grant/Reseal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrantServer).Reseal(ctx, req.(*GrantAndGrantSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grant_PutSeal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaintextAndGrantSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrantServer).PutSeal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Grant/PutSeal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrantServer).PutSeal(ctx, req.(*PlaintextAndGrantSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grant_UnsealGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(grant.Grant)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrantServer).UnsealGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Grant/UnsealGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrantServer).UnsealGet(ctx, req.(*grant.Grant))
	}
	return interceptor(ctx, in, info, handler)
}

var _Grant_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Grant",
	HandlerType: (*GrantServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Seal",
			Handler:    _Grant_Seal_Handler,
		},
		{
			MethodName: "Unseal",
			Handler:    _Grant_Unseal_Handler,
		},
		{
			MethodName: "Reseal",
			Handler:    _Grant_Reseal_Handler,
		},
		{
			MethodName: "PutSeal",
			Handler:    _Grant_PutSeal_Handler,
		},
		{
			MethodName: "UnsealGet",
			Handler:    _Grant_UnsealGet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

// CleartextClient is the client API for Cleartext service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CleartextClient interface {
	// Push some plaintext data into storage and get its deterministically
	// generated secret reference.
	Put(ctx context.Context, in *Plaintext, opts ...grpc.CallOption) (*reference.Ref, error)
	// Provide a secret reference to an encrypted blob and get the plaintext
	// data back.
	Get(ctx context.Context, in *reference.Ref, opts ...grpc.CallOption) (*Plaintext, error)
}

type cleartextClient struct {
	cc *grpc.ClientConn
}

func NewCleartextClient(cc *grpc.ClientConn) CleartextClient {
	return &cleartextClient{cc}
}

func (c *cleartextClient) Put(ctx context.Context, in *Plaintext, opts ...grpc.CallOption) (*reference.Ref, error) {
	out := new(reference.Ref)
	err := c.cc.Invoke(ctx, "/api.Cleartext/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cleartextClient) Get(ctx context.Context, in *reference.Ref, opts ...grpc.CallOption) (*Plaintext, error) {
	out := new(Plaintext)
	err := c.cc.Invoke(ctx, "/api.Cleartext/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CleartextServer is the server API for Cleartext service.
type CleartextServer interface {
	// Push some plaintext data into storage and get its deterministically
	// generated secret reference.
	Put(context.Context, *Plaintext) (*reference.Ref, error)
	// Provide a secret reference to an encrypted blob and get the plaintext
	// data back.
	Get(context.Context, *reference.Ref) (*Plaintext, error)
}

func RegisterCleartextServer(s *grpc.Server, srv CleartextServer) {
	s.RegisterService(&_Cleartext_serviceDesc, srv)
}

func _Cleartext_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Plaintext)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CleartextServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Cleartext/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CleartextServer).Put(ctx, req.(*Plaintext))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cleartext_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(reference.Ref)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CleartextServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Cleartext/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CleartextServer).Get(ctx, req.(*reference.Ref))
	}
	return interceptor(ctx, in, info, handler)
}

var _Cleartext_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Cleartext",
	HandlerType: (*CleartextServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _Cleartext_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Cleartext_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

// EncryptionClient is the client API for Encryption service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EncryptionClient interface {
	// Encrypt some data and get its deterministically generated
	// secret reference including its address without storing the data.
	Encrypt(ctx context.Context, in *Plaintext, opts ...grpc.CallOption) (*ReferenceAndCiphertext, error)
	// Decrypt the provided data by supplying it alongside its secret
	// reference. The address is not used for decryption and may be omitted.
	Decrypt(ctx context.Context, in *ReferenceAndCiphertext, opts ...grpc.CallOption) (*Plaintext, error)
}

type encryptionClient struct {
	cc *grpc.ClientConn
}

func NewEncryptionClient(cc *grpc.ClientConn) EncryptionClient {
	return &encryptionClient{cc}
}

func (c *encryptionClient) Encrypt(ctx context.Context, in *Plaintext, opts ...grpc.CallOption) (*ReferenceAndCiphertext, error) {
	out := new(ReferenceAndCiphertext)
	err := c.cc.Invoke(ctx, "/api.Encryption/Encrypt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *encryptionClient) Decrypt(ctx context.Context, in *ReferenceAndCiphertext, opts ...grpc.CallOption) (*Plaintext, error) {
	out := new(Plaintext)
	err := c.cc.Invoke(ctx, "/api.Encryption/Decrypt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EncryptionServer is the server API for Encryption service.
type EncryptionServer interface {
	// Encrypt some data and get its deterministically generated
	// secret reference including its address without storing the data.
	Encrypt(context.Context, *Plaintext) (*ReferenceAndCiphertext, error)
	// Decrypt the provided data by supplying it alongside its secret
	// reference. The address is not used for decryption and may be omitted.
	Decrypt(context.Context, *ReferenceAndCiphertext) (*Plaintext, error)
}

func RegisterEncryptionServer(s *grpc.Server, srv EncryptionServer) {
	s.RegisterService(&_Encryption_serviceDesc, srv)
}

func _Encryption_Encrypt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Plaintext)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EncryptionServer).Encrypt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Encryption/Encrypt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EncryptionServer).Encrypt(ctx, req.(*Plaintext))
	}
	return interceptor(ctx, in, info, handler)
}

func _Encryption_Decrypt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReferenceAndCiphertext)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EncryptionServer).Decrypt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Encryption/Decrypt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EncryptionServer).Decrypt(ctx, req.(*ReferenceAndCiphertext))
	}
	return interceptor(ctx, in, info, handler)
}

var _Encryption_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Encryption",
	HandlerType: (*EncryptionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Encrypt",
			Handler:    _Encryption_Encrypt_Handler,
		},
		{
			MethodName: "Decrypt",
			Handler:    _Encryption_Decrypt_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StorageClient interface {
	// Insert the (presumably) encrypted data provided and get the its address.
	Push(ctx context.Context, in *Ciphertext, opts ...grpc.CallOption) (*Address, error)
	// Retrieve the (presumably) encrypted data stored at address.
	Pull(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Ciphertext, error)
	// Get some information about the encrypted blob stored at an address,
	// including whether it exists.
	Stat(ctx context.Context, in *Address, opts ...grpc.CallOption) (*stores.StatInfo, error)
}

type storageClient struct {
	cc *grpc.ClientConn
}

func NewStorageClient(cc *grpc.ClientConn) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) Push(ctx context.Context, in *Ciphertext, opts ...grpc.CallOption) (*Address, error) {
	out := new(Address)
	err := c.cc.Invoke(ctx, "/api.Storage/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) Pull(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Ciphertext, error) {
	out := new(Ciphertext)
	err := c.cc.Invoke(ctx, "/api.Storage/Pull", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) Stat(ctx context.Context, in *Address, opts ...grpc.CallOption) (*stores.StatInfo, error) {
	out := new(stores.StatInfo)
	err := c.cc.Invoke(ctx, "/api.Storage/Stat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageServer is the server API for Storage service.
type StorageServer interface {
	// Insert the (presumably) encrypted data provided and get the its address.
	Push(context.Context, *Ciphertext) (*Address, error)
	// Retrieve the (presumably) encrypted data stored at address.
	Pull(context.Context, *Address) (*Ciphertext, error)
	// Get some information about the encrypted blob stored at an address,
	// including whether it exists.
	Stat(context.Context, *Address) (*stores.StatInfo, error)
}

func RegisterStorageServer(s *grpc.Server, srv StorageServer) {
	s.RegisterService(&_Storage_serviceDesc, srv)
}

func _Storage_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ciphertext)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Push(ctx, req.(*Ciphertext))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_Pull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Pull(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/Pull",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Pull(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_Stat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).Stat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/Stat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).Stat(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

var _Storage_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _Storage_Push_Handler,
		},
		{
			MethodName: "Pull",
			Handler:    _Storage_Pull_Handler,
		},
		{
			MethodName: "Stat",
			Handler:    _Storage_Stat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

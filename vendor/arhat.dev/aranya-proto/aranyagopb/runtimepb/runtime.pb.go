// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: runtime/runtime.proto

package runtimepb

import (
	bytes "bytes"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strconv "strconv"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type PacketType int32

const (
	_INVALID_RUNTIME_DATA PacketType = 0
	// 1-9: runtime basic operation
	CMD_GET_INFO     PacketType = 1
	CMD_EXEC         PacketType = 2
	CMD_ATTACH       PacketType = 3
	CMD_LOGS         PacketType = 4
	CMD_TTY_RESIZE   PacketType = 5
	CMD_PORT_FORWARD PacketType = 6
	MSG_RUNTIME_INFO PacketType = 9
	MSG_ERROR        PacketType = 10
	// 11-19: container image / application bundle operations
	CMD_IMAGE_LIST   PacketType = 11
	CMD_IMAGE_ENSURE PacketType = 12
	CMD_IMAGE_DELETE PacketType = 13
	// 21-29: image msgs
	MSG_IMAGE_STATUS      PacketType = 21
	MSG_IMAGE_STATUS_LIST PacketType = 22
	// 31-39: pod provisioning
	CMD_POD_LIST   PacketType = 31
	CMD_POD_ENSURE PacketType = 32
	CMD_POD_DELETE PacketType = 33
	// 41-49: pod msgs
	MSG_POD_STATUS      PacketType = 41
	MSG_POD_STATUS_LIST PacketType = 42
)

var PacketType_name = map[int32]string{
	0:  "_INVALID_RUNTIME_DATA",
	1:  "CMD_GET_INFO",
	2:  "CMD_EXEC",
	3:  "CMD_ATTACH",
	4:  "CMD_LOGS",
	5:  "CMD_TTY_RESIZE",
	6:  "CMD_PORT_FORWARD",
	9:  "MSG_RUNTIME_INFO",
	10: "MSG_ERROR",
	11: "CMD_IMAGE_LIST",
	12: "CMD_IMAGE_ENSURE",
	13: "CMD_IMAGE_DELETE",
	21: "MSG_IMAGE_STATUS",
	22: "MSG_IMAGE_STATUS_LIST",
	31: "CMD_POD_LIST",
	32: "CMD_POD_ENSURE",
	33: "CMD_POD_DELETE",
	41: "MSG_POD_STATUS",
	42: "MSG_POD_STATUS_LIST",
}

var PacketType_value = map[string]int32{
	"_INVALID_RUNTIME_DATA": 0,
	"CMD_GET_INFO":          1,
	"CMD_EXEC":              2,
	"CMD_ATTACH":            3,
	"CMD_LOGS":              4,
	"CMD_TTY_RESIZE":        5,
	"CMD_PORT_FORWARD":      6,
	"MSG_RUNTIME_INFO":      9,
	"MSG_ERROR":             10,
	"CMD_IMAGE_LIST":        11,
	"CMD_IMAGE_ENSURE":      12,
	"CMD_IMAGE_DELETE":      13,
	"MSG_IMAGE_STATUS":      21,
	"MSG_IMAGE_STATUS_LIST": 22,
	"CMD_POD_LIST":          31,
	"CMD_POD_ENSURE":        32,
	"CMD_POD_DELETE":        33,
	"MSG_POD_STATUS":        41,
	"MSG_POD_STATUS_LIST":   42,
}

func (PacketType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ab2d388759451feb, []int{0}
}

type Packet struct {
	Kind    PacketType `protobuf:"varint,1,opt,name=kind,proto3,enum=runtime.PacketType" json:"kind,omitempty"`
	Payload []byte     `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *Packet) Reset()      { *m = Packet{} }
func (*Packet) ProtoMessage() {}
func (*Packet) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab2d388759451feb, []int{0}
}
func (m *Packet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Packet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Packet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Packet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet.Merge(m, src)
}
func (m *Packet) XXX_Size() int {
	return m.Size()
}
func (m *Packet) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet.DiscardUnknown(m)
}

var xxx_messageInfo_Packet proto.InternalMessageInfo

func (m *Packet) GetKind() PacketType {
	if m != nil {
		return m.Kind
	}
	return _INVALID_RUNTIME_DATA
}

func (m *Packet) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

// runtime info to override node info
type RuntimeInfo struct {
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version       string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Os            string `protobuf:"bytes,3,opt,name=os,proto3" json:"os,omitempty"`
	OsImage       string `protobuf:"bytes,4,opt,name=os_image,json=osImage,proto3" json:"os_image,omitempty"`
	Arch          string `protobuf:"bytes,5,opt,name=arch,proto3" json:"arch,omitempty"`
	KernelVersion string `protobuf:"bytes,6,opt,name=kernel_version,json=kernelVersion,proto3" json:"kernel_version,omitempty"`
}

func (m *RuntimeInfo) Reset()      { *m = RuntimeInfo{} }
func (*RuntimeInfo) ProtoMessage() {}
func (*RuntimeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab2d388759451feb, []int{1}
}
func (m *RuntimeInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RuntimeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RuntimeInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RuntimeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuntimeInfo.Merge(m, src)
}
func (m *RuntimeInfo) XXX_Size() int {
	return m.Size()
}
func (m *RuntimeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_RuntimeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_RuntimeInfo proto.InternalMessageInfo

func (m *RuntimeInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RuntimeInfo) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *RuntimeInfo) GetOs() string {
	if m != nil {
		return m.Os
	}
	return ""
}

func (m *RuntimeInfo) GetOsImage() string {
	if m != nil {
		return m.OsImage
	}
	return ""
}

func (m *RuntimeInfo) GetArch() string {
	if m != nil {
		return m.Arch
	}
	return ""
}

func (m *RuntimeInfo) GetKernelVersion() string {
	if m != nil {
		return m.KernelVersion
	}
	return ""
}

func init() {
	proto.RegisterEnum("runtime.PacketType", PacketType_name, PacketType_value)
	proto.RegisterType((*Packet)(nil), "runtime.Packet")
	proto.RegisterType((*RuntimeInfo)(nil), "runtime.RuntimeInfo")
}

func init() { proto.RegisterFile("runtime/runtime.proto", fileDescriptor_ab2d388759451feb) }

var fileDescriptor_ab2d388759451feb = []byte{
	// 512 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0xc1, 0x6e, 0xd3, 0x30,
	0x1c, 0xc6, 0xe3, 0xb6, 0xeb, 0x56, 0xaf, 0xad, 0x2c, 0x8f, 0x42, 0x76, 0x31, 0x65, 0x12, 0x62,
	0x0c, 0xd1, 0x49, 0xf0, 0x04, 0xa1, 0xf1, 0x4a, 0x44, 0xdb, 0x54, 0x8e, 0x3b, 0x60, 0x97, 0xc8,
	0x5d, 0xc3, 0x56, 0x75, 0x4d, 0xaa, 0xa4, 0x4c, 0xea, 0x8d, 0x07, 0xe0, 0xc0, 0x2b, 0x70, 0xe3,
	0x51, 0x38, 0xf6, 0xb8, 0x23, 0x4d, 0x2f, 0x1c, 0xf7, 0x08, 0xc8, 0x71, 0x0c, 0x6c, 0xa7, 0xf8,
	0xfb, 0xfd, 0xfd, 0x7d, 0x5f, 0xfe, 0x92, 0x61, 0x23, 0xfe, 0x1c, 0x2e, 0x26, 0xb3, 0xe0, 0x38,
	0xff, 0xb6, 0xe6, 0x71, 0xb4, 0x88, 0xf0, 0x76, 0x2e, 0x0f, 0xde, 0xc1, 0xf2, 0x40, 0x9c, 0x4f,
	0x83, 0x05, 0x7e, 0x06, 0x4b, 0xd3, 0x49, 0x38, 0x36, 0x41, 0x13, 0x1c, 0xd6, 0x5f, 0xed, 0xb5,
	0xb4, 0x41, 0x8d, 0xf9, 0x72, 0x1e, 0xb0, 0xec, 0x02, 0x36, 0xe1, 0xf6, 0x5c, 0x2c, 0xaf, 0x22,
	0x31, 0x36, 0x0b, 0x4d, 0x70, 0x58, 0x65, 0x5a, 0x1e, 0x7c, 0x07, 0x70, 0x97, 0x29, 0x9b, 0x13,
	0x7e, 0x8a, 0x30, 0x86, 0xa5, 0x50, 0xcc, 0x82, 0x2c, 0xb2, 0xc2, 0xb2, 0xb3, 0x74, 0x5f, 0x07,
	0x71, 0x32, 0x89, 0xc2, 0xcc, 0x5d, 0x61, 0x5a, 0xe2, 0x3a, 0x2c, 0x44, 0x89, 0x59, 0xcc, 0x60,
	0x21, 0x4a, 0xf0, 0x3e, 0xdc, 0x89, 0x12, 0x7f, 0x32, 0x13, 0x17, 0x81, 0x59, 0x52, 0x57, 0xa3,
	0xc4, 0x91, 0x52, 0x06, 0x8b, 0xf8, 0xfc, 0xd2, 0xdc, 0x52, 0xc1, 0xf2, 0x8c, 0x9f, 0xc2, 0xfa,
	0x34, 0x88, 0xc3, 0xe0, 0xca, 0xd7, 0xf9, 0xe5, 0x6c, 0x5a, 0x53, 0xf4, 0x54, 0xc1, 0xa3, 0xaf,
	0x45, 0x08, 0xff, 0xad, 0x84, 0xf7, 0x61, 0xc3, 0x77, 0xfa, 0xa7, 0x56, 0xd7, 0xb1, 0x7d, 0x36,
	0xec, 0x73, 0xa7, 0x47, 0x7d, 0xdb, 0xe2, 0x16, 0x32, 0x30, 0x82, 0xd5, 0x76, 0xcf, 0xf6, 0x3b,
	0x94, 0xfb, 0x4e, 0xff, 0xc4, 0x45, 0x00, 0x57, 0xe1, 0x8e, 0x24, 0xf4, 0x03, 0x6d, 0xa3, 0x02,
	0xae, 0x43, 0x28, 0x95, 0xc5, 0xb9, 0xd5, 0x7e, 0x8b, 0x8a, 0x7a, 0xda, 0x75, 0x3b, 0x1e, 0x2a,
	0x61, 0x0c, 0xeb, 0x52, 0x71, 0xfe, 0xd1, 0x67, 0xd4, 0x73, 0xce, 0x28, 0xda, 0xc2, 0x0f, 0x20,
	0x92, 0x6c, 0xe0, 0x32, 0xee, 0x9f, 0xb8, 0xec, 0xbd, 0xc5, 0x6c, 0x54, 0x96, 0xb4, 0xe7, 0x75,
	0xfe, 0xb6, 0x67, 0x5d, 0x15, 0x5c, 0x83, 0x15, 0x49, 0x29, 0x63, 0x2e, 0x43, 0x50, 0xc7, 0x39,
	0x3d, 0xab, 0x43, 0xfd, 0xae, 0xe3, 0x71, 0xb4, 0xab, 0xe3, 0x14, 0xa3, 0x7d, 0x6f, 0xc8, 0x28,
	0xaa, 0xde, 0xa5, 0x36, 0xed, 0x52, 0x4e, 0x51, 0x4d, 0x97, 0x28, 0xea, 0x71, 0x8b, 0x0f, 0x3d,
	0xd4, 0x90, 0xdb, 0xdf, 0xa7, 0x2a, 0xfc, 0xa1, 0xde, 0x7e, 0xe0, 0xda, 0x8a, 0x3c, 0xd6, 0xbf,
	0x20, 0x49, 0x5e, 0xd6, 0xfc, 0x9f, 0xe5, 0x55, 0x4f, 0x24, 0x93, 0xa1, 0x92, 0xe5, 0x45, 0xcf,
	0xf1, 0x23, 0xb8, 0x77, 0x97, 0xa9, 0xd0, 0xa3, 0x37, 0x62, 0xb5, 0x26, 0xc6, 0xcd, 0x9a, 0x18,
	0xb7, 0x6b, 0x02, 0xbe, 0xa4, 0x04, 0xfc, 0x48, 0x09, 0xf8, 0x99, 0x12, 0xb0, 0x4a, 0x09, 0xf8,
	0x95, 0x12, 0xf0, 0x3b, 0x25, 0xc6, 0x6d, 0x4a, 0xc0, 0xb7, 0x0d, 0x31, 0x56, 0x1b, 0x62, 0xdc,
	0x6c, 0x88, 0x71, 0xf6, 0x42, 0xc4, 0x97, 0x62, 0xd1, 0x1a, 0x07, 0xd7, 0xc7, 0x22, 0x16, 0xe1,
	0x52, 0xbc, 0xcc, 0x9e, 0x75, 0x2e, 0x2e, 0xa2, 0xf9, 0x48, 0xbf, 0xf7, 0xf9, 0x68, 0x54, 0xce,
	0x66, 0xaf, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x76, 0xe7, 0x54, 0x52, 0x0b, 0x03, 0x00, 0x00,
}

func (x PacketType) String() string {
	s, ok := PacketType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *Packet) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Packet)
	if !ok {
		that2, ok := that.(Packet)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Kind != that1.Kind {
		return false
	}
	if !bytes.Equal(this.Payload, that1.Payload) {
		return false
	}
	return true
}
func (this *RuntimeInfo) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RuntimeInfo)
	if !ok {
		that2, ok := that.(RuntimeInfo)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Version != that1.Version {
		return false
	}
	if this.Os != that1.Os {
		return false
	}
	if this.OsImage != that1.OsImage {
		return false
	}
	if this.Arch != that1.Arch {
		return false
	}
	if this.KernelVersion != that1.KernelVersion {
		return false
	}
	return true
}
func (this *Packet) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&runtimepb.Packet{")
	s = append(s, "Kind: "+fmt.Sprintf("%#v", this.Kind)+",\n")
	s = append(s, "Payload: "+fmt.Sprintf("%#v", this.Payload)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *RuntimeInfo) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&runtimepb.RuntimeInfo{")
	s = append(s, "Name: "+fmt.Sprintf("%#v", this.Name)+",\n")
	s = append(s, "Version: "+fmt.Sprintf("%#v", this.Version)+",\n")
	s = append(s, "Os: "+fmt.Sprintf("%#v", this.Os)+",\n")
	s = append(s, "OsImage: "+fmt.Sprintf("%#v", this.OsImage)+",\n")
	s = append(s, "Arch: "+fmt.Sprintf("%#v", this.Arch)+",\n")
	s = append(s, "KernelVersion: "+fmt.Sprintf("%#v", this.KernelVersion)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringRuntime(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Packet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Packet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Packet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintRuntime(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x12
	}
	if m.Kind != 0 {
		i = encodeVarintRuntime(dAtA, i, uint64(m.Kind))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *RuntimeInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RuntimeInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RuntimeInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.KernelVersion) > 0 {
		i -= len(m.KernelVersion)
		copy(dAtA[i:], m.KernelVersion)
		i = encodeVarintRuntime(dAtA, i, uint64(len(m.KernelVersion)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Arch) > 0 {
		i -= len(m.Arch)
		copy(dAtA[i:], m.Arch)
		i = encodeVarintRuntime(dAtA, i, uint64(len(m.Arch)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.OsImage) > 0 {
		i -= len(m.OsImage)
		copy(dAtA[i:], m.OsImage)
		i = encodeVarintRuntime(dAtA, i, uint64(len(m.OsImage)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Os) > 0 {
		i -= len(m.Os)
		copy(dAtA[i:], m.Os)
		i = encodeVarintRuntime(dAtA, i, uint64(len(m.Os)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Version) > 0 {
		i -= len(m.Version)
		copy(dAtA[i:], m.Version)
		i = encodeVarintRuntime(dAtA, i, uint64(len(m.Version)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintRuntime(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRuntime(dAtA []byte, offset int, v uint64) int {
	offset -= sovRuntime(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Packet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Kind != 0 {
		n += 1 + sovRuntime(uint64(m.Kind))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovRuntime(uint64(l))
	}
	return n
}

func (m *RuntimeInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovRuntime(uint64(l))
	}
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovRuntime(uint64(l))
	}
	l = len(m.Os)
	if l > 0 {
		n += 1 + l + sovRuntime(uint64(l))
	}
	l = len(m.OsImage)
	if l > 0 {
		n += 1 + l + sovRuntime(uint64(l))
	}
	l = len(m.Arch)
	if l > 0 {
		n += 1 + l + sovRuntime(uint64(l))
	}
	l = len(m.KernelVersion)
	if l > 0 {
		n += 1 + l + sovRuntime(uint64(l))
	}
	return n
}

func sovRuntime(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRuntime(x uint64) (n int) {
	return sovRuntime(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Packet) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Packet{`,
		`Kind:` + fmt.Sprintf("%v", this.Kind) + `,`,
		`Payload:` + fmt.Sprintf("%v", this.Payload) + `,`,
		`}`,
	}, "")
	return s
}
func (this *RuntimeInfo) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RuntimeInfo{`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`Version:` + fmt.Sprintf("%v", this.Version) + `,`,
		`Os:` + fmt.Sprintf("%v", this.Os) + `,`,
		`OsImage:` + fmt.Sprintf("%v", this.OsImage) + `,`,
		`Arch:` + fmt.Sprintf("%v", this.Arch) + `,`,
		`KernelVersion:` + fmt.Sprintf("%v", this.KernelVersion) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringRuntime(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Packet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRuntime
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Packet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Packet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Kind", wireType)
			}
			m.Kind = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Kind |= PacketType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthRuntime
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthRuntime
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRuntime(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRuntime
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRuntime
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RuntimeInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRuntime
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RuntimeInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RuntimeInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRuntime
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRuntime
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRuntime
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRuntime
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Os", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRuntime
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRuntime
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Os = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OsImage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRuntime
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRuntime
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OsImage = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Arch", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRuntime
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRuntime
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Arch = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KernelVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRuntime
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRuntime
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KernelVersion = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRuntime(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRuntime
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRuntime
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRuntime(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRuntime
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRuntime
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthRuntime
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRuntime
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRuntime
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRuntime        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRuntime          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRuntime = fmt.Errorf("proto: unexpected end of group")
)

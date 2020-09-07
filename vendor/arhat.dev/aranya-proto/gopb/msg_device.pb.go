// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: msg_device.proto

// +build !nodev

package gopb

import (
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

type DeviceStatus_State int32

const (
	DEVICE_STATE_UNKNOWN   DeviceStatus_State = 0
	DEVICE_STATE_CREATED   DeviceStatus_State = 1
	DEVICE_STATE_CONNECTED DeviceStatus_State = 2
	DEVICE_STATE_ERRORED   DeviceStatus_State = 3
)

var DeviceStatus_State_name = map[int32]string{
	0: "DEVICE_STATE_UNKNOWN",
	1: "DEVICE_STATE_CREATED",
	2: "DEVICE_STATE_CONNECTED",
	3: "DEVICE_STATE_ERRORED",
}

var DeviceStatus_State_value = map[string]int32{
	"DEVICE_STATE_UNKNOWN":   0,
	"DEVICE_STATE_CREATED":   1,
	"DEVICE_STATE_CONNECTED": 2,
	"DEVICE_STATE_ERRORED":   3,
}

func (DeviceStatus_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0fd5ffbd6aabb0be, []int{0, 0}
}

type DeviceStatus struct {
	Id    string             `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	State DeviceStatus_State `protobuf:"varint,2,opt,name=state,proto3,enum=aranya.DeviceStatus_State" json:"state,omitempty"`
	// message description for this state
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *DeviceStatus) Reset()      { *m = DeviceStatus{} }
func (*DeviceStatus) ProtoMessage() {}
func (*DeviceStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fd5ffbd6aabb0be, []int{0}
}
func (m *DeviceStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DeviceStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DeviceStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DeviceStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeviceStatus.Merge(m, src)
}
func (m *DeviceStatus) XXX_Size() int {
	return m.Size()
}
func (m *DeviceStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_DeviceStatus.DiscardUnknown(m)
}

var xxx_messageInfo_DeviceStatus proto.InternalMessageInfo

func (m *DeviceStatus) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DeviceStatus) GetState() DeviceStatus_State {
	if m != nil {
		return m.State
	}
	return DEVICE_STATE_UNKNOWN
}

func (m *DeviceStatus) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type DeviceStatusList struct {
	Devices []*DeviceStatus `protobuf:"bytes,1,rep,name=devices,proto3" json:"devices,omitempty"`
}

func (m *DeviceStatusList) Reset()      { *m = DeviceStatusList{} }
func (*DeviceStatusList) ProtoMessage() {}
func (*DeviceStatusList) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fd5ffbd6aabb0be, []int{1}
}
func (m *DeviceStatusList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DeviceStatusList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DeviceStatusList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DeviceStatusList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeviceStatusList.Merge(m, src)
}
func (m *DeviceStatusList) XXX_Size() int {
	return m.Size()
}
func (m *DeviceStatusList) XXX_DiscardUnknown() {
	xxx_messageInfo_DeviceStatusList.DiscardUnknown(m)
}

var xxx_messageInfo_DeviceStatusList proto.InternalMessageInfo

func (m *DeviceStatusList) GetDevices() []*DeviceStatus {
	if m != nil {
		return m.Devices
	}
	return nil
}

func init() {
	proto.RegisterEnum("aranya.DeviceStatus_State", DeviceStatus_State_name, DeviceStatus_State_value)
	proto.RegisterType((*DeviceStatus)(nil), "aranya.DeviceStatus")
	proto.RegisterType((*DeviceStatusList)(nil), "aranya.DeviceStatusList")
}

func init() { proto.RegisterFile("msg_device.proto", fileDescriptor_0fd5ffbd6aabb0be) }

var fileDescriptor_0fd5ffbd6aabb0be = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0x2d, 0x4e, 0x8f,
	0x4f, 0x49, 0x2d, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4b, 0x2c,
	0x4a, 0xcc, 0xab, 0x4c, 0x54, 0xba, 0xcb, 0xc8, 0xc5, 0xe3, 0x02, 0x96, 0x08, 0x2e, 0x49, 0x2c,
	0x29, 0x2d, 0x16, 0xe2, 0xe3, 0x62, 0xca, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x62,
	0xca, 0x4c, 0x11, 0x32, 0xe0, 0x62, 0x2d, 0x2e, 0x49, 0x2c, 0x49, 0x95, 0x60, 0x52, 0x60, 0xd4,
	0xe0, 0x33, 0x92, 0xd2, 0x83, 0x68, 0xd4, 0x43, 0xd6, 0xa4, 0x07, 0xa2, 0x52, 0x83, 0x20, 0x0a,
	0x85, 0x24, 0xb8, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x25, 0x98, 0xc1, 0xc6, 0xc0,
	0xb8, 0x4a, 0x85, 0x5c, 0xac, 0xc1, 0x50, 0x25, 0x22, 0x2e, 0xae, 0x61, 0x9e, 0xce, 0xae, 0xf1,
	0xc1, 0x21, 0x8e, 0x21, 0xae, 0xf1, 0xa1, 0x7e, 0xde, 0x7e, 0xfe, 0xe1, 0x7e, 0x02, 0x0c, 0x18,
	0x32, 0xce, 0x41, 0xae, 0x8e, 0x21, 0xae, 0x2e, 0x02, 0x8c, 0x42, 0x52, 0x5c, 0x62, 0xa8, 0x32,
	0xfe, 0x7e, 0x7e, 0xae, 0xce, 0x20, 0x39, 0x26, 0x0c, 0x5d, 0xae, 0x41, 0x41, 0xfe, 0x41, 0xae,
	0x2e, 0x02, 0xcc, 0x4a, 0x4e, 0x5c, 0x02, 0xc8, 0x2e, 0xf5, 0xc9, 0x2c, 0x2e, 0x11, 0xd2, 0xe3,
	0x62, 0x87, 0x84, 0x45, 0xb1, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0xb7, 0x91, 0x08, 0x36, 0x4f, 0x05,
	0xc1, 0x14, 0x39, 0x05, 0x5e, 0x78, 0x28, 0xc7, 0x70, 0xe3, 0xa1, 0x1c, 0xc3, 0x87, 0x87, 0x72,
	0x8c, 0x0d, 0x8f, 0xe4, 0x18, 0x57, 0x3c, 0x92, 0x63, 0x3c, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23,
	0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x5f, 0x3c, 0x92, 0x63, 0xf8, 0xf0, 0x48, 0x8e, 0x71, 0xc2,
	0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88, 0x92, 0x4e, 0x2c, 0xca,
	0x48, 0x2c, 0xd1, 0x4b, 0x49, 0x2d, 0xd3, 0x87, 0x98, 0xae, 0x0b, 0x0e, 0x79, 0xfd, 0xf4, 0xfc,
	0x82, 0xa4, 0x24, 0x36, 0x30, 0xdb, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x0e, 0x84, 0x8b, 0x4d,
	0x99, 0x01, 0x00, 0x00,
}

func (x DeviceStatus_State) String() string {
	s, ok := DeviceStatus_State_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *DeviceStatus) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*DeviceStatus)
	if !ok {
		that2, ok := that.(DeviceStatus)
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
	if this.Id != that1.Id {
		return false
	}
	if this.State != that1.State {
		return false
	}
	if this.Message != that1.Message {
		return false
	}
	return true
}
func (this *DeviceStatusList) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*DeviceStatusList)
	if !ok {
		that2, ok := that.(DeviceStatusList)
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
	if len(this.Devices) != len(that1.Devices) {
		return false
	}
	for i := range this.Devices {
		if !this.Devices[i].Equal(that1.Devices[i]) {
			return false
		}
	}
	return true
}
func (this *DeviceStatus) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&gopb.DeviceStatus{")
	s = append(s, "Id: "+fmt.Sprintf("%#v", this.Id)+",\n")
	s = append(s, "State: "+fmt.Sprintf("%#v", this.State)+",\n")
	s = append(s, "Message: "+fmt.Sprintf("%#v", this.Message)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *DeviceStatusList) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&gopb.DeviceStatusList{")
	if this.Devices != nil {
		s = append(s, "Devices: "+fmt.Sprintf("%#v", this.Devices)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringMsgDevice(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *DeviceStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeviceStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DeviceStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintMsgDevice(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x1a
	}
	if m.State != 0 {
		i = encodeVarintMsgDevice(dAtA, i, uint64(m.State))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintMsgDevice(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DeviceStatusList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeviceStatusList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DeviceStatusList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Devices) > 0 {
		for iNdEx := len(m.Devices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Devices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMsgDevice(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintMsgDevice(dAtA []byte, offset int, v uint64) int {
	offset -= sovMsgDevice(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DeviceStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovMsgDevice(uint64(l))
	}
	if m.State != 0 {
		n += 1 + sovMsgDevice(uint64(m.State))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovMsgDevice(uint64(l))
	}
	return n
}

func (m *DeviceStatusList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Devices) > 0 {
		for _, e := range m.Devices {
			l = e.Size()
			n += 1 + l + sovMsgDevice(uint64(l))
		}
	}
	return n
}

func sovMsgDevice(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMsgDevice(x uint64) (n int) {
	return sovMsgDevice(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *DeviceStatus) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&DeviceStatus{`,
		`Id:` + fmt.Sprintf("%v", this.Id) + `,`,
		`State:` + fmt.Sprintf("%v", this.State) + `,`,
		`Message:` + fmt.Sprintf("%v", this.Message) + `,`,
		`}`,
	}, "")
	return s
}
func (this *DeviceStatusList) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForDevices := "[]*DeviceStatus{"
	for _, f := range this.Devices {
		repeatedStringForDevices += strings.Replace(f.String(), "DeviceStatus", "DeviceStatus", 1) + ","
	}
	repeatedStringForDevices += "}"
	s := strings.Join([]string{`&DeviceStatusList{`,
		`Devices:` + repeatedStringForDevices + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringMsgDevice(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *DeviceStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgDevice
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
			return fmt.Errorf("proto: DeviceStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeviceStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgDevice
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
				return ErrInvalidLengthMsgDevice
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgDevice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgDevice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= DeviceStatus_State(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgDevice
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
				return ErrInvalidLengthMsgDevice
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgDevice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgDevice(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMsgDevice
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMsgDevice
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
func (m *DeviceStatusList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgDevice
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
			return fmt.Errorf("proto: DeviceStatusList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeviceStatusList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Devices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgDevice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMsgDevice
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgDevice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Devices = append(m.Devices, &DeviceStatus{})
			if err := m.Devices[len(m.Devices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgDevice(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMsgDevice
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMsgDevice
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
func skipMsgDevice(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMsgDevice
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
					return 0, ErrIntOverflowMsgDevice
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
					return 0, ErrIntOverflowMsgDevice
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
				return 0, ErrInvalidLengthMsgDevice
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMsgDevice
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMsgDevice
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMsgDevice        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMsgDevice          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMsgDevice = fmt.Errorf("proto: unexpected end of group")
)
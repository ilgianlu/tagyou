// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.0
// source: kura-payload.proto

package kura

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

type KuraPayload_KuraMetric_ValueType int32

const (
	KuraPayload_KuraMetric_DOUBLE KuraPayload_KuraMetric_ValueType = 0
	KuraPayload_KuraMetric_FLOAT  KuraPayload_KuraMetric_ValueType = 1
	KuraPayload_KuraMetric_INT64  KuraPayload_KuraMetric_ValueType = 2
	KuraPayload_KuraMetric_INT32  KuraPayload_KuraMetric_ValueType = 3
	KuraPayload_KuraMetric_BOOL   KuraPayload_KuraMetric_ValueType = 4
	KuraPayload_KuraMetric_STRING KuraPayload_KuraMetric_ValueType = 5
	KuraPayload_KuraMetric_BYTES  KuraPayload_KuraMetric_ValueType = 6
)

// Enum value maps for KuraPayload_KuraMetric_ValueType.
var (
	KuraPayload_KuraMetric_ValueType_name = map[int32]string{
		0: "DOUBLE",
		1: "FLOAT",
		2: "INT64",
		3: "INT32",
		4: "BOOL",
		5: "STRING",
		6: "BYTES",
	}
	KuraPayload_KuraMetric_ValueType_value = map[string]int32{
		"DOUBLE": 0,
		"FLOAT":  1,
		"INT64":  2,
		"INT32":  3,
		"BOOL":   4,
		"STRING": 5,
		"BYTES":  6,
	}
)

func (x KuraPayload_KuraMetric_ValueType) Enum() *KuraPayload_KuraMetric_ValueType {
	p := new(KuraPayload_KuraMetric_ValueType)
	*p = x
	return p
}

func (x KuraPayload_KuraMetric_ValueType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (KuraPayload_KuraMetric_ValueType) Descriptor() protoreflect.EnumDescriptor {
	return file_kura_payload_proto_enumTypes[0].Descriptor()
}

func (KuraPayload_KuraMetric_ValueType) Type() protoreflect.EnumType {
	return &file_kura_payload_proto_enumTypes[0]
}

func (x KuraPayload_KuraMetric_ValueType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use KuraPayload_KuraMetric_ValueType.Descriptor instead.
func (KuraPayload_KuraMetric_ValueType) EnumDescriptor() ([]byte, []int) {
	return file_kura_payload_proto_rawDescGZIP(), []int{0, 0, 0}
}

type KuraPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64                     `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Position  *KuraPayload_KuraPosition `protobuf:"bytes,2,opt,name=position,proto3" json:"position,omitempty"`
	Metric    []*KuraPayload_KuraMetric `protobuf:"bytes,5000,rep,name=metric,proto3" json:"metric,omitempty"` // can be zero, so optional
	Body      []byte                    `protobuf:"bytes,5001,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *KuraPayload) Reset() {
	*x = KuraPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kura_payload_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KuraPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KuraPayload) ProtoMessage() {}

func (x *KuraPayload) ProtoReflect() protoreflect.Message {
	mi := &file_kura_payload_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KuraPayload.ProtoReflect.Descriptor instead.
func (*KuraPayload) Descriptor() ([]byte, []int) {
	return file_kura_payload_proto_rawDescGZIP(), []int{0}
}

func (x *KuraPayload) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *KuraPayload) GetPosition() *KuraPayload_KuraPosition {
	if x != nil {
		return x.Position
	}
	return nil
}

func (x *KuraPayload) GetMetric() []*KuraPayload_KuraMetric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *KuraPayload) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type KuraPayload_KuraMetric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string                           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type        KuraPayload_KuraMetric_ValueType `protobuf:"varint,2,opt,name=type,proto3,enum=kura.KuraPayload_KuraMetric_ValueType" json:"type,omitempty"`
	DoubleValue float64                          `protobuf:"fixed64,3,opt,name=double_value,json=doubleValue,proto3" json:"double_value,omitempty"`
	FloatValue  float32                          `protobuf:"fixed32,4,opt,name=float_value,json=floatValue,proto3" json:"float_value,omitempty"`
	LongValue   int64                            `protobuf:"varint,5,opt,name=long_value,json=longValue,proto3" json:"long_value,omitempty"`
	IntValue    int32                            `protobuf:"varint,6,opt,name=int_value,json=intValue,proto3" json:"int_value,omitempty"`
	BoolValue   bool                             `protobuf:"varint,7,opt,name=bool_value,json=boolValue,proto3" json:"bool_value,omitempty"`
	StringValue string                           `protobuf:"bytes,8,opt,name=string_value,json=stringValue,proto3" json:"string_value,omitempty"`
	BytesValue  []byte                           `protobuf:"bytes,9,opt,name=bytes_value,json=bytesValue,proto3" json:"bytes_value,omitempty"`
}

func (x *KuraPayload_KuraMetric) Reset() {
	*x = KuraPayload_KuraMetric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kura_payload_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KuraPayload_KuraMetric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KuraPayload_KuraMetric) ProtoMessage() {}

func (x *KuraPayload_KuraMetric) ProtoReflect() protoreflect.Message {
	mi := &file_kura_payload_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KuraPayload_KuraMetric.ProtoReflect.Descriptor instead.
func (*KuraPayload_KuraMetric) Descriptor() ([]byte, []int) {
	return file_kura_payload_proto_rawDescGZIP(), []int{0, 0}
}

func (x *KuraPayload_KuraMetric) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *KuraPayload_KuraMetric) GetType() KuraPayload_KuraMetric_ValueType {
	if x != nil {
		return x.Type
	}
	return KuraPayload_KuraMetric_DOUBLE
}

func (x *KuraPayload_KuraMetric) GetDoubleValue() float64 {
	if x != nil {
		return x.DoubleValue
	}
	return 0
}

func (x *KuraPayload_KuraMetric) GetFloatValue() float32 {
	if x != nil {
		return x.FloatValue
	}
	return 0
}

func (x *KuraPayload_KuraMetric) GetLongValue() int64 {
	if x != nil {
		return x.LongValue
	}
	return 0
}

func (x *KuraPayload_KuraMetric) GetIntValue() int32 {
	if x != nil {
		return x.IntValue
	}
	return 0
}

func (x *KuraPayload_KuraMetric) GetBoolValue() bool {
	if x != nil {
		return x.BoolValue
	}
	return false
}

func (x *KuraPayload_KuraMetric) GetStringValue() string {
	if x != nil {
		return x.StringValue
	}
	return ""
}

func (x *KuraPayload_KuraMetric) GetBytesValue() []byte {
	if x != nil {
		return x.BytesValue
	}
	return nil
}

type KuraPayload_KuraPosition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude   float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude  float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	Altitude   float64 `protobuf:"fixed64,3,opt,name=altitude,proto3" json:"altitude,omitempty"`
	Precision  float64 `protobuf:"fixed64,4,opt,name=precision,proto3" json:"precision,omitempty"` // dilution of precision of the current satellite fix.
	Heading    float64 `protobuf:"fixed64,5,opt,name=heading,proto3" json:"heading,omitempty"`     // heading in degrees
	Speed      float64 `protobuf:"fixed64,6,opt,name=speed,proto3" json:"speed,omitempty"`         // meters per second
	Timestamp  int64   `protobuf:"varint,7,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Satellites int32   `protobuf:"varint,8,opt,name=satellites,proto3" json:"satellites,omitempty"` // number satellites locked by the GPS device
	Status     int32   `protobuf:"varint,9,opt,name=status,proto3" json:"status,omitempty"`         // status indicator for the GPS data: 1 = no GPS response; 2 = error in response; 4 = valid.
}

func (x *KuraPayload_KuraPosition) Reset() {
	*x = KuraPayload_KuraPosition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kura_payload_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KuraPayload_KuraPosition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KuraPayload_KuraPosition) ProtoMessage() {}

func (x *KuraPayload_KuraPosition) ProtoReflect() protoreflect.Message {
	mi := &file_kura_payload_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KuraPayload_KuraPosition.ProtoReflect.Descriptor instead.
func (*KuraPayload_KuraPosition) Descriptor() ([]byte, []int) {
	return file_kura_payload_proto_rawDescGZIP(), []int{0, 1}
}

func (x *KuraPayload_KuraPosition) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *KuraPayload_KuraPosition) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

func (x *KuraPayload_KuraPosition) GetAltitude() float64 {
	if x != nil {
		return x.Altitude
	}
	return 0
}

func (x *KuraPayload_KuraPosition) GetPrecision() float64 {
	if x != nil {
		return x.Precision
	}
	return 0
}

func (x *KuraPayload_KuraPosition) GetHeading() float64 {
	if x != nil {
		return x.Heading
	}
	return 0
}

func (x *KuraPayload_KuraPosition) GetSpeed() float64 {
	if x != nil {
		return x.Speed
	}
	return 0
}

func (x *KuraPayload_KuraPosition) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *KuraPayload_KuraPosition) GetSatellites() int32 {
	if x != nil {
		return x.Satellites
	}
	return 0
}

func (x *KuraPayload_KuraPosition) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

var File_kura_payload_proto protoreflect.FileDescriptor

var file_kura_payload_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6b, 0x75, 0x72, 0x61, 0x2d, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6b, 0x75, 0x72, 0x61, 0x22, 0xdb, 0x06, 0x0a, 0x0b, 0x4b,
	0x75, 0x72, 0x61, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x3a, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6b, 0x75, 0x72,
	0x61, 0x2e, 0x4b, 0x75, 0x72, 0x61, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x4b, 0x75,
	0x72, 0x61, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x35, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x88,
	0x27, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6b, 0x75, 0x72, 0x61, 0x2e, 0x4b, 0x75, 0x72,
	0x61, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x4b, 0x75, 0x72, 0x61, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x13, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x89, 0x27, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x1a, 0x9a, 0x03, 0x0a, 0x0a, 0x4b, 0x75, 0x72, 0x61, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x3a, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x26, 0x2e, 0x6b, 0x75, 0x72, 0x61, 0x2e, 0x4b, 0x75, 0x72, 0x61, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x4b, 0x75, 0x72, 0x61, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0a, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x6e, 0x67, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x6e, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x62, 0x6f, 0x6f, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21,
	0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x62, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x59, 0x0a, 0x09, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0a, 0x0a, 0x06, 0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x46,
	0x4c, 0x4f, 0x41, 0x54, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x49, 0x4e, 0x54, 0x36, 0x34, 0x10,
	0x02, 0x12, 0x09, 0x0a, 0x05, 0x49, 0x4e, 0x54, 0x33, 0x32, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04,
	0x42, 0x4f, 0x4f, 0x4c, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47,
	0x10, 0x05, 0x12, 0x09, 0x0a, 0x05, 0x42, 0x59, 0x54, 0x45, 0x53, 0x10, 0x06, 0x1a, 0x88, 0x02,
	0x0a, 0x0c, 0x4b, 0x75, 0x72, 0x61, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a,
	0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f,
	0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6c,
	0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x6c, 0x74, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x61, 0x6c, 0x74, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x70, 0x72, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x70, 0x65, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x73, 0x70, 0x65,
	0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x61, 0x74, 0x65, 0x6c, 0x6c, 0x69, 0x74, 0x65, 0x73, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x61, 0x74, 0x65, 0x6c, 0x6c, 0x69, 0x74, 0x65, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x1b, 0x5a, 0x19, 0x63, 0x6f, 0x64, 0x65,
	0x2e, 0x63, 0x61, 0x6d, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6d, 0x61, 0x67, 0x6e, 0x69,
	0x2f, 0x6b, 0x75, 0x72, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kura_payload_proto_rawDescOnce sync.Once
	file_kura_payload_proto_rawDescData = file_kura_payload_proto_rawDesc
)

func file_kura_payload_proto_rawDescGZIP() []byte {
	file_kura_payload_proto_rawDescOnce.Do(func() {
		file_kura_payload_proto_rawDescData = protoimpl.X.CompressGZIP(file_kura_payload_proto_rawDescData)
	})
	return file_kura_payload_proto_rawDescData
}

var file_kura_payload_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_kura_payload_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_kura_payload_proto_goTypes = []interface{}{
	(KuraPayload_KuraMetric_ValueType)(0), // 0: kura.KuraPayload.KuraMetric.ValueType
	(*KuraPayload)(nil),                   // 1: kura.KuraPayload
	(*KuraPayload_KuraMetric)(nil),        // 2: kura.KuraPayload.KuraMetric
	(*KuraPayload_KuraPosition)(nil),      // 3: kura.KuraPayload.KuraPosition
}
var file_kura_payload_proto_depIdxs = []int32{
	3, // 0: kura.KuraPayload.position:type_name -> kura.KuraPayload.KuraPosition
	2, // 1: kura.KuraPayload.metric:type_name -> kura.KuraPayload.KuraMetric
	0, // 2: kura.KuraPayload.KuraMetric.type:type_name -> kura.KuraPayload.KuraMetric.ValueType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_kura_payload_proto_init() }
func file_kura_payload_proto_init() {
	if File_kura_payload_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kura_payload_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KuraPayload); i {
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
		file_kura_payload_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KuraPayload_KuraMetric); i {
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
		file_kura_payload_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KuraPayload_KuraPosition); i {
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
			RawDescriptor: file_kura_payload_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kura_payload_proto_goTypes,
		DependencyIndexes: file_kura_payload_proto_depIdxs,
		EnumInfos:         file_kura_payload_proto_enumTypes,
		MessageInfos:      file_kura_payload_proto_msgTypes,
	}.Build()
	File_kura_payload_proto = out.File
	file_kura_payload_proto_rawDesc = nil
	file_kura_payload_proto_goTypes = nil
	file_kura_payload_proto_depIdxs = nil
}

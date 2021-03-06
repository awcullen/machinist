// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.20.1
// 	protoc        v3.11.4
// source: opcua/write.proto

package opcua

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// WriteValue message.
type WriteValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId      string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	AttributeId uint32   `protobuf:"varint,2,opt,name=attributeId,proto3" json:"attributeId,omitempty"`
	IndexRange  string   `protobuf:"bytes,3,opt,name=indexRange,proto3" json:"indexRange,omitempty"`
	Value       *Variant `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *WriteValue) Reset() {
	*x = WriteValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_opcua_write_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteValue) ProtoMessage() {}

func (x *WriteValue) ProtoReflect() protoreflect.Message {
	mi := &file_opcua_write_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteValue.ProtoReflect.Descriptor instead.
func (*WriteValue) Descriptor() ([]byte, []int) {
	return file_opcua_write_proto_rawDescGZIP(), []int{0}
}

func (x *WriteValue) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *WriteValue) GetAttributeId() uint32 {
	if x != nil {
		return x.AttributeId
	}
	return 0
}

func (x *WriteValue) GetIndexRange() string {
	if x != nil {
		return x.IndexRange
	}
	return ""
}

func (x *WriteValue) GetValue() *Variant {
	if x != nil {
		return x.Value
	}
	return nil
}

// WriteRequest message
type WriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestHandle uint32        `protobuf:"varint,1,opt,name=requestHandle,proto3" json:"requestHandle,omitempty"`
	NodesToWrite  []*WriteValue `protobuf:"bytes,2,rep,name=nodesToWrite,proto3" json:"nodesToWrite,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_opcua_write_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_opcua_write_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteRequest.ProtoReflect.Descriptor instead.
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return file_opcua_write_proto_rawDescGZIP(), []int{1}
}

func (x *WriteRequest) GetRequestHandle() uint32 {
	if x != nil {
		return x.RequestHandle
	}
	return 0
}

func (x *WriteRequest) GetNodesToWrite() []*WriteValue {
	if x != nil {
		return x.NodesToWrite
	}
	return nil
}

// WriteResponse message.
type WriteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestHandle uint32   `protobuf:"varint,1,opt,name=requestHandle,proto3" json:"requestHandle,omitempty"`
	StatusCode    uint32   `protobuf:"varint,2,opt,name=statusCode,proto3" json:"statusCode,omitempty"`
	Results       []uint32 `protobuf:"varint,3,rep,packed,name=results,proto3" json:"results,omitempty"`
}

func (x *WriteResponse) Reset() {
	*x = WriteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_opcua_write_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteResponse) ProtoMessage() {}

func (x *WriteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_opcua_write_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteResponse.ProtoReflect.Descriptor instead.
func (*WriteResponse) Descriptor() ([]byte, []int) {
	return file_opcua_write_proto_rawDescGZIP(), []int{2}
}

func (x *WriteResponse) GetRequestHandle() uint32 {
	if x != nil {
		return x.RequestHandle
	}
	return 0
}

func (x *WriteResponse) GetStatusCode() uint32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *WriteResponse) GetResults() []uint32 {
	if x != nil {
		return x.Results
	}
	return nil
}

var File_opcua_write_proto protoreflect.FileDescriptor

var file_opcua_write_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6f, 0x70, 0x63, 0x75, 0x61, 0x2f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6f, 0x70, 0x63, 0x75, 0x61, 0x1a, 0x13, 0x6f, 0x70, 0x63, 0x75,
	0x61, 0x2f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x8c, 0x01, 0x0a, 0x0a, 0x57, 0x72, 0x69, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x61, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6f, 0x70, 0x63, 0x75, 0x61, 0x2e,
	0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x6b,
	0x0a, 0x0c, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24,
	0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x12, 0x35, 0x0a, 0x0c, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x54, 0x6f, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6f, 0x70, 0x63,
	0x75, 0x61, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0c, 0x6e,
	0x6f, 0x64, 0x65, 0x73, 0x54, 0x6f, 0x57, 0x72, 0x69, 0x74, 0x65, 0x22, 0x6f, 0x0a, 0x0d, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x0d,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0d, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x32, 0x44, 0x0a, 0x0c,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x34, 0x0a, 0x05,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x13, 0x2e, 0x6f, 0x70, 0x63, 0x75, 0x61, 0x2e, 0x57, 0x72,
	0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x6f, 0x70, 0x63,
	0x75, 0x61, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x6f, 0x70, 0x63, 0x75, 0x61, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_opcua_write_proto_rawDescOnce sync.Once
	file_opcua_write_proto_rawDescData = file_opcua_write_proto_rawDesc
)

func file_opcua_write_proto_rawDescGZIP() []byte {
	file_opcua_write_proto_rawDescOnce.Do(func() {
		file_opcua_write_proto_rawDescData = protoimpl.X.CompressGZIP(file_opcua_write_proto_rawDescData)
	})
	return file_opcua_write_proto_rawDescData
}

var file_opcua_write_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_opcua_write_proto_goTypes = []interface{}{
	(*WriteValue)(nil),    // 0: opcua.WriteValue
	(*WriteRequest)(nil),  // 1: opcua.WriteRequest
	(*WriteResponse)(nil), // 2: opcua.WriteResponse
	(*Variant)(nil),       // 3: opcua.Variant
}
var file_opcua_write_proto_depIdxs = []int32{
	3, // 0: opcua.WriteValue.value:type_name -> opcua.Variant
	0, // 1: opcua.WriteRequest.nodesToWrite:type_name -> opcua.WriteValue
	1, // 2: opcua.WriteService.Write:input_type -> opcua.WriteRequest
	2, // 3: opcua.WriteService.Write:output_type -> opcua.WriteResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_opcua_write_proto_init() }
func file_opcua_write_proto_init() {
	if File_opcua_write_proto != nil {
		return
	}
	file_opcua_variant_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_opcua_write_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteValue); i {
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
		file_opcua_write_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteRequest); i {
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
		file_opcua_write_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteResponse); i {
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
			RawDescriptor: file_opcua_write_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_opcua_write_proto_goTypes,
		DependencyIndexes: file_opcua_write_proto_depIdxs,
		MessageInfos:      file_opcua_write_proto_msgTypes,
	}.Build()
	File_opcua_write_proto = out.File
	file_opcua_write_proto_rawDesc = nil
	file_opcua_write_proto_goTypes = nil
	file_opcua_write_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// WriteServiceClient is the client API for WriteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WriteServiceClient interface {
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
}

type writeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWriteServiceClient(cc grpc.ClientConnInterface) WriteServiceClient {
	return &writeServiceClient{cc}
}

func (c *writeServiceClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, "/opcua.WriteService/Write", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WriteServiceServer is the server API for WriteService service.
type WriteServiceServer interface {
	Write(context.Context, *WriteRequest) (*WriteResponse, error)
}

// UnimplementedWriteServiceServer can be embedded to have forward compatible implementations.
type UnimplementedWriteServiceServer struct {
}

func (*UnimplementedWriteServiceServer) Write(context.Context, *WriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}

func RegisterWriteServiceServer(s *grpc.Server, srv WriteServiceServer) {
	s.RegisterService(&_WriteService_serviceDesc, srv)
}

func _WriteService_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriteServiceServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/opcua.WriteService/Write",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriteServiceServer).Write(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _WriteService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "opcua.WriteService",
	HandlerType: (*WriteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Write",
			Handler:    _WriteService_Write_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "opcua/write.proto",
}

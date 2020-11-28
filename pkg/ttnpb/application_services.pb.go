// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/application_services.proto

package ttnpb

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() {
	proto.RegisterFile("lorawan-stack/api/application_services.proto", fileDescriptor_f6c42f4fe8e3c902)
}
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/application_services.proto", fileDescriptor_f6c42f4fe8e3c902)
}

var fileDescriptor_f6c42f4fe8e3c902 = []byte{
	// 855 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0x4f, 0x48, 0xdc, 0x58,
	0x1c, 0xce, 0x73, 0x97, 0x81, 0x7d, 0xba, 0x2b, 0xbe, 0x85, 0x5d, 0x18, 0xdd, 0xc7, 0x92, 0x45,
	0x5d, 0xdc, 0x9d, 0x64, 0xd7, 0x61, 0xb7, 0x58, 0xa4, 0xad, 0x7f, 0xca, 0x54, 0x6c, 0xa9, 0x28,
	0xbd, 0xcc, 0xc5, 0x66, 0xc6, 0x67, 0x0c, 0x33, 0x4d, 0xd2, 0xbc, 0xa7, 0x76, 0x1c, 0x04, 0xdb,
	0x93, 0x78, 0x6a, 0x29, 0x2d, 0xa5, 0xf4, 0xd0, 0x8b, 0xd4, 0x4b, 0xc1, 0xa3, 0xb7, 0x7a, 0xf4,
	0x28, 0xf4, 0x50, 0x8f, 0x4e, 0x52, 0xa8, 0x87, 0x16, 0x3c, 0x7a, 0x2c, 0x79, 0x49, 0x30, 0x99,
	0x19, 0x33, 0xce, 0xd8, 0x9b, 0xf3, 0xde, 0xf7, 0x7e, 0xbf, 0xef, 0xf7, 0xfd, 0xf2, 0x7d, 0x08,
	0xff, 0x2e, 0x1a, 0x96, 0xb2, 0xac, 0xe8, 0x29, 0xca, 0x94, 0x7c, 0x41, 0x56, 0x4c, 0x4d, 0x56,
	0x4c, 0xb3, 0xa8, 0xe5, 0x15, 0xa6, 0x19, 0xfa, 0x2c, 0x25, 0xd6, 0x92, 0x96, 0x27, 0x54, 0x32,
	0x2d, 0x83, 0x19, 0xe8, 0x27, 0xc6, 0x74, 0xc9, 0x7f, 0x21, 0x2d, 0xa5, 0x93, 0x3d, 0xaa, 0x61,
	0xa8, 0x45, 0xe2, 0x3d, 0xd3, 0x75, 0x83, 0xf1, 0x57, 0x3e, 0x3a, 0xd9, 0xed, 0xdf, 0xf2, 0x5f,
	0xb9, 0xc5, 0x79, 0x99, 0xdc, 0x33, 0x59, 0xc9, 0xbf, 0xfc, 0x23, 0xb6, 0xf1, 0xd9, 0x20, 0x6d,
	0x8e, 0xe8, 0x4c, 0x9b, 0xd7, 0x88, 0x15, 0xb4, 0xc1, 0xb5, 0x20, 0x4b, 0x53, 0x17, 0x98, 0x7f,
	0x3f, 0xf8, 0x25, 0x01, 0x7f, 0x1e, 0x39, 0x2d, 0x3d, 0x4d, 0x54, 0x8d, 0x32, 0xab, 0x84, 0x1c,
	0x00, 0x13, 0x63, 0x16, 0x51, 0x18, 0x41, 0x7f, 0x4a, 0xd1, 0xc1, 0x24, 0xef, 0x3c, 0xf2, 0xea,
	0xfe, 0x22, 0xa1, 0x2c, 0xd9, 0x5d, 0x8d, 0x0c, 0x61, 0xc4, 0x27, 0xe0, 0xd1, 0xfb, 0x8f, 0x4f,
	0xdb, 0x36, 0x80, 0x98, 0x96, 0x17, 0x29, 0xb1, 0xa8, 0x5c, 0xce, 0x1b, 0xc5, 0xa2, 0x92, 0x33,
	0x2c, 0x85, 0x19, 0x96, 0xe4, 0x9e, 0xcd, 0x6a, 0x73, 0x34, 0xf8, 0x63, 0x35, 0x3c, 0x32, 0xbd,
	0x0c, 0x06, 0xb2, 0x53, 0xe2, 0xa4, 0x6c, 0x58, 0xaa, 0xa2, 0x6b, 0x2b, 0xde, 0x61, 0x55, 0x85,
	0xf0, 0x1d, 0xaf, 0x54, 0x75, 0x50, 0x53, 0x11, 0x3d, 0x04, 0xf0, 0xbb, 0x0c, 0x61, 0xa8, 0xb7,
	0x9a, 0x78, 0x86, 0xb0, 0x66, 0xe7, 0xfb, 0x9f, 0x8f, 0xf7, 0x0f, 0x92, 0x22, 0x5d, 0xe4, 0x72,
	0xf8, 0x8b, 0x71, 0x49, 0x45, 0x7f, 0xaf, 0xa2, 0xcf, 0x00, 0x7e, 0x7f, 0x53, 0xa3, 0x0c, 0xf5,
	0x57, 0x57, 0x77, 0x4f, 0x43, 0x1d, 0x68, 0x40, 0xa3, 0x27, 0x86, 0x06, 0x15, 0x5f, 0x79, 0x3a,
	0x3f, 0x03, 0xe8, 0xc7, 0x08, 0x93, 0xec, 0x7f, 0xa8, 0x15, 0xe1, 0xb3, 0xb7, 0xd0, 0xb7, 0x54,
	0x1d, 0x6d, 0x00, 0x98, 0xb8, 0x63, 0xce, 0xd5, 0xfd, 0xb0, 0xbc, 0xf3, 0x66, 0x85, 0x1f, 0xe2,
	0xf3, 0xa6, 0x93, 0x31, 0xc2, 0x4b, 0x75, 0x84, 0x77, 0xf7, 0x6f, 0xc2, 0xc4, 0x38, 0x29, 0x12,
	0x46, 0x50, 0x5f, 0x4c, 0x87, 0x89, 0x53, 0x57, 0x25, 0x7f, 0x91, 0x3c, 0xdf, 0x4a, 0x81, 0x6f,
	0xa5, 0xeb, 0xae, 0x6f, 0xc5, 0x3e, 0x4e, 0xe2, 0xf7, 0x01, 0x1c, 0xbb, 0xfd, 0xd5, 0xc1, 0x77,
	0xed, 0xb0, 0x2b, 0x54, 0x7a, 0x24, 0x9f, 0x27, 0x94, 0xa2, 0x32, 0x84, 0xee, 0xb2, 0xa7, 0xb9,
	0x33, 0x9b, 0xe0, 0x52, 0x85, 0xf3, 0xde, 0x8b, 0x29, 0xce, 0xa5, 0x1f, 0xf5, 0xc6, 0x73, 0xf1,
	0x83, 0x00, 0xbd, 0x04, 0xb0, 0xc3, 0xb7, 0xf4, 0xd4, 0xc4, 0x24, 0x29, 0x21, 0xa9, 0xa1, 0xe1,
	0x3d, 0x60, 0xb0, 0x9d, 0x1a, 0x1e, 0xde, 0xb5, 0x38, 0xca, 0x79, 0x0c, 0x8b, 0x97, 0x9a, 0x73,
	0x84, 0x1b, 0x52, 0xa9, 0x02, 0x29, 0x71, 0x87, 0x3e, 0x07, 0xb0, 0x9d, 0xfb, 0x80, 0x97, 0xa4,
	0x28, 0xd5, 0xc0, 0x24, 0x3e, 0x2e, 0xa0, 0xf6, 0x6b, 0x7d, 0x6a, 0x54, 0xbc, 0xca, 0xb9, 0x0d,
	0xa1, 0x56, 0xb9, 0xb9, 0xaa, 0xfd, 0xe0, 0xa6, 0x84, 0x27, 0xd9, 0x5f, 0xf1, 0x01, 0x72, 0x3e,
	0xbd, 0x6e, 0x70, 0x4e, 0xa3, 0xe8, 0x5a, 0x8b, 0x9c, 0xe4, 0x72, 0x81, 0x94, 0x78, 0xa6, 0xbc,
	0x01, 0xb0, 0xc3, 0x37, 0xd3, 0x19, 0x2b, 0xad, 0xb1, 0xda, 0xf9, 0x28, 0xde, 0xe6, 0x14, 0x27,
	0x92, 0xe3, 0x2d, 0x53, 0x54, 0x4c, 0x6d, 0xb6, 0x40, 0x4a, 0x92, 0xef, 0xc0, 0x0f, 0x6d, 0xb0,
	0x33, 0x43, 0xd8, 0x58, 0x28, 0x51, 0xd0, 0xbf, 0xf1, 0x62, 0x86, 0xb1, 0x01, 0xdf, 0xfe, 0x3a,
	0x4f, 0xa2, 0x38, 0x6a, 0x1a, 0x3a, 0x25, 0xe2, 0x27, 0x2f, 0x1d, 0x0f, 0x41, 0x36, 0x87, 0xee,
	0x36, 0x39, 0x44, 0x38, 0xf6, 0x78, 0x92, 0x36, 0x0a, 0xd2, 0xec, 0x0a, 0x7a, 0x70, 0x91, 0x1e,
	0xe1, 0x24, 0x6d, 0x36, 0x75, 0xd1, 0x26, 0x80, 0x9d, 0x33, 0x8d, 0x94, 0x9d, 0x69, 0xa8, 0xec,
	0x59, 0x81, 0x97, 0xe1, 0x3a, 0x8e, 0x24, 0x87, 0x2f, 0x30, 0x20, 0x77, 0xf8, 0x5b, 0x00, 0xbb,
	0x5c, 0x13, 0x87, 0x9b, 0x53, 0x94, 0x6e, 0xe0, 0xf3, 0x08, 0x3a, 0xe0, 0xfa, 0x5b, 0x4d, 0x70,
	0x85, 0x51, 0xe2, 0x38, 0xa7, 0x7c, 0x05, 0x5d, 0x88, 0xf2, 0xe8, 0x26, 0xd8, 0xab, 0x60, 0xb0,
	0x5f, 0xc1, 0xe0, 0xa0, 0x82, 0x85, 0xc3, 0x0a, 0x16, 0x8e, 0x2a, 0x58, 0x38, 0xae, 0x60, 0xe1,
	0xa4, 0x82, 0xc1, 0x9a, 0x8d, 0xc1, 0xba, 0x8d, 0x85, 0x2d, 0x1b, 0x83, 0x6d, 0x1b, 0x0b, 0x3b,
	0x36, 0x16, 0x76, 0x6d, 0x2c, 0xec, 0xd9, 0x18, 0xec, 0xdb, 0x18, 0x1c, 0xd8, 0x58, 0x38, 0xb4,
	0x31, 0x38, 0xb2, 0xb1, 0x70, 0x6c, 0x63, 0x70, 0x62, 0x63, 0x61, 0xcd, 0xc1, 0xc2, 0xba, 0x83,
	0xc1, 0x63, 0x07, 0x0b, 0x2f, 0x1c, 0x0c, 0x5e, 0x3b, 0x58, 0xd8, 0x72, 0xb0, 0xb0, 0xed, 0x60,
	0xb0, 0xe3, 0x60, 0xb0, 0xeb, 0x60, 0x90, 0x95, 0x55, 0x43, 0x62, 0x0b, 0x84, 0x2d, 0x68, 0xba,
	0x4a, 0x25, 0x9d, 0xb0, 0x65, 0xc3, 0x2a, 0xc8, 0xd1, 0xff, 0xee, 0x96, 0xd2, 0xb2, 0x59, 0x50,
	0x65, 0xc6, 0x74, 0x33, 0x97, 0x4b, 0xf0, 0x85, 0xa5, 0xbf, 0x06, 0x00, 0x00, 0xff, 0xff, 0x4e,
	0x38, 0x8a, 0xd9, 0xc5, 0x0a, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ApplicationRegistryClient is the client API for ApplicationRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ApplicationRegistryClient interface {
	// Create a new application. This also sets the given organization or user as
	// first collaborator with all possible rights.
	Create(ctx context.Context, in *CreateApplicationRequest, opts ...grpc.CallOption) (*Application, error)
	// Get the application with the given identifiers, selecting the fields specified
	// in the field mask.
	// More or less fields may be returned, depending on the rights of the caller.
	Get(ctx context.Context, in *GetApplicationRequest, opts ...grpc.CallOption) (*Application, error)
	// List applications where the given user or organization is a direct collaborator.
	// If no user or organization is given, this returns the applications the caller
	// has access to.
	// Similar to Get, this selects the fields given by the field mask.
	// More or less fields may be returned, depending on the rights of the caller.
	List(ctx context.Context, in *ListApplicationsRequest, opts ...grpc.CallOption) (*Applications, error)
	// Update the application, changing the fields specified by the field mask to the provided values.
	Update(ctx context.Context, in *UpdateApplicationRequest, opts ...grpc.CallOption) (*Application, error)
	// Delete the application. This may not release the application ID for reuse.
	// All end devices must be deleted from the application before it can be deleted.
	Delete(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*types.Empty, error)
}

type applicationRegistryClient struct {
	cc *grpc.ClientConn
}

func NewApplicationRegistryClient(cc *grpc.ClientConn) ApplicationRegistryClient {
	return &applicationRegistryClient{cc}
}

func (c *applicationRegistryClient) Create(ctx context.Context, in *CreateApplicationRequest, opts ...grpc.CallOption) (*Application, error) {
	out := new(Application)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationRegistry/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationRegistryClient) Get(ctx context.Context, in *GetApplicationRequest, opts ...grpc.CallOption) (*Application, error) {
	out := new(Application)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationRegistry/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationRegistryClient) List(ctx context.Context, in *ListApplicationsRequest, opts ...grpc.CallOption) (*Applications, error) {
	out := new(Applications)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationRegistry/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationRegistryClient) Update(ctx context.Context, in *UpdateApplicationRequest, opts ...grpc.CallOption) (*Application, error) {
	out := new(Application)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationRegistry/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationRegistryClient) Delete(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationRegistry/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApplicationRegistryServer is the server API for ApplicationRegistry service.
type ApplicationRegistryServer interface {
	// Create a new application. This also sets the given organization or user as
	// first collaborator with all possible rights.
	Create(context.Context, *CreateApplicationRequest) (*Application, error)
	// Get the application with the given identifiers, selecting the fields specified
	// in the field mask.
	// More or less fields may be returned, depending on the rights of the caller.
	Get(context.Context, *GetApplicationRequest) (*Application, error)
	// List applications where the given user or organization is a direct collaborator.
	// If no user or organization is given, this returns the applications the caller
	// has access to.
	// Similar to Get, this selects the fields given by the field mask.
	// More or less fields may be returned, depending on the rights of the caller.
	List(context.Context, *ListApplicationsRequest) (*Applications, error)
	// Update the application, changing the fields specified by the field mask to the provided values.
	Update(context.Context, *UpdateApplicationRequest) (*Application, error)
	// Delete the application. This may not release the application ID for reuse.
	// All end devices must be deleted from the application before it can be deleted.
	Delete(context.Context, *ApplicationIdentifiers) (*types.Empty, error)
}

// UnimplementedApplicationRegistryServer can be embedded to have forward compatible implementations.
type UnimplementedApplicationRegistryServer struct {
}

func (*UnimplementedApplicationRegistryServer) Create(ctx context.Context, req *CreateApplicationRequest) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedApplicationRegistryServer) Get(ctx context.Context, req *GetApplicationRequest) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedApplicationRegistryServer) List(ctx context.Context, req *ListApplicationsRequest) (*Applications, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (*UnimplementedApplicationRegistryServer) Update(ctx context.Context, req *UpdateApplicationRequest) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedApplicationRegistryServer) Delete(ctx context.Context, req *ApplicationIdentifiers) (*types.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func RegisterApplicationRegistryServer(s *grpc.Server, srv ApplicationRegistryServer) {
	s.RegisterService(&_ApplicationRegistry_serviceDesc, srv)
}

func _ApplicationRegistry_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateApplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationRegistryServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationRegistry/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationRegistryServer).Create(ctx, req.(*CreateApplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationRegistry_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetApplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationRegistryServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationRegistry/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationRegistryServer).Get(ctx, req.(*GetApplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationRegistry_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListApplicationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationRegistryServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationRegistry/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationRegistryServer).List(ctx, req.(*ListApplicationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationRegistry_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateApplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationRegistryServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationRegistry/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationRegistryServer).Update(ctx, req.(*UpdateApplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationRegistry_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplicationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationRegistryServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationRegistry/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationRegistryServer).Delete(ctx, req.(*ApplicationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

var _ApplicationRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.ApplicationRegistry",
	HandlerType: (*ApplicationRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _ApplicationRegistry_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ApplicationRegistry_Get_Handler,
		},
		{
			MethodName: "List",
			Handler:    _ApplicationRegistry_List_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ApplicationRegistry_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ApplicationRegistry_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/application_services.proto",
}

// ApplicationAccessClient is the client API for ApplicationAccess service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ApplicationAccessClient interface {
	// List the rights the caller has on this application.
	ListRights(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*Rights, error)
	// Create an API key scoped to this application.
	CreateAPIKey(ctx context.Context, in *CreateApplicationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error)
	// List the API keys for this application.
	ListAPIKeys(ctx context.Context, in *ListApplicationAPIKeysRequest, opts ...grpc.CallOption) (*APIKeys, error)
	// Get a single API key of this application.
	GetAPIKey(ctx context.Context, in *GetApplicationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error)
	// Update the rights of an API key of the application.
	// This method can also be used to delete the API key, by giving it no rights.
	// The caller is required to have all assigned or/and removed rights.
	UpdateAPIKey(ctx context.Context, in *UpdateApplicationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error)
	// Get the rights of a collaborator (member) of the application.
	// Pseudo-rights in the response (such as the "_ALL" right) are not expanded.
	GetCollaborator(ctx context.Context, in *GetApplicationCollaboratorRequest, opts ...grpc.CallOption) (*GetCollaboratorResponse, error)
	// Set the rights of a collaborator (member) on the application.
	// This method can also be used to delete the collaborator, by giving them no rights.
	// The caller is required to have all assigned or/and removed rights.
	SetCollaborator(ctx context.Context, in *SetApplicationCollaboratorRequest, opts ...grpc.CallOption) (*types.Empty, error)
	// List the collaborators on this application.
	ListCollaborators(ctx context.Context, in *ListApplicationCollaboratorsRequest, opts ...grpc.CallOption) (*Collaborators, error)
}

type applicationAccessClient struct {
	cc *grpc.ClientConn
}

func NewApplicationAccessClient(cc *grpc.ClientConn) ApplicationAccessClient {
	return &applicationAccessClient{cc}
}

func (c *applicationAccessClient) ListRights(ctx context.Context, in *ApplicationIdentifiers, opts ...grpc.CallOption) (*Rights, error) {
	out := new(Rights)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationAccess/ListRights", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationAccessClient) CreateAPIKey(ctx context.Context, in *CreateApplicationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error) {
	out := new(APIKey)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationAccess/CreateAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationAccessClient) ListAPIKeys(ctx context.Context, in *ListApplicationAPIKeysRequest, opts ...grpc.CallOption) (*APIKeys, error) {
	out := new(APIKeys)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationAccess/ListAPIKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationAccessClient) GetAPIKey(ctx context.Context, in *GetApplicationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error) {
	out := new(APIKey)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationAccess/GetAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationAccessClient) UpdateAPIKey(ctx context.Context, in *UpdateApplicationAPIKeyRequest, opts ...grpc.CallOption) (*APIKey, error) {
	out := new(APIKey)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationAccess/UpdateAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationAccessClient) GetCollaborator(ctx context.Context, in *GetApplicationCollaboratorRequest, opts ...grpc.CallOption) (*GetCollaboratorResponse, error) {
	out := new(GetCollaboratorResponse)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationAccess/GetCollaborator", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationAccessClient) SetCollaborator(ctx context.Context, in *SetApplicationCollaboratorRequest, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationAccess/SetCollaborator", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationAccessClient) ListCollaborators(ctx context.Context, in *ListApplicationCollaboratorsRequest, opts ...grpc.CallOption) (*Collaborators, error) {
	out := new(Collaborators)
	err := c.cc.Invoke(ctx, "/ttn.lorawan.v3.ApplicationAccess/ListCollaborators", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApplicationAccessServer is the server API for ApplicationAccess service.
type ApplicationAccessServer interface {
	// List the rights the caller has on this application.
	ListRights(context.Context, *ApplicationIdentifiers) (*Rights, error)
	// Create an API key scoped to this application.
	CreateAPIKey(context.Context, *CreateApplicationAPIKeyRequest) (*APIKey, error)
	// List the API keys for this application.
	ListAPIKeys(context.Context, *ListApplicationAPIKeysRequest) (*APIKeys, error)
	// Get a single API key of this application.
	GetAPIKey(context.Context, *GetApplicationAPIKeyRequest) (*APIKey, error)
	// Update the rights of an API key of the application.
	// This method can also be used to delete the API key, by giving it no rights.
	// The caller is required to have all assigned or/and removed rights.
	UpdateAPIKey(context.Context, *UpdateApplicationAPIKeyRequest) (*APIKey, error)
	// Get the rights of a collaborator (member) of the application.
	// Pseudo-rights in the response (such as the "_ALL" right) are not expanded.
	GetCollaborator(context.Context, *GetApplicationCollaboratorRequest) (*GetCollaboratorResponse, error)
	// Set the rights of a collaborator (member) on the application.
	// This method can also be used to delete the collaborator, by giving them no rights.
	// The caller is required to have all assigned or/and removed rights.
	SetCollaborator(context.Context, *SetApplicationCollaboratorRequest) (*types.Empty, error)
	// List the collaborators on this application.
	ListCollaborators(context.Context, *ListApplicationCollaboratorsRequest) (*Collaborators, error)
}

// UnimplementedApplicationAccessServer can be embedded to have forward compatible implementations.
type UnimplementedApplicationAccessServer struct {
}

func (*UnimplementedApplicationAccessServer) ListRights(ctx context.Context, req *ApplicationIdentifiers) (*Rights, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRights not implemented")
}
func (*UnimplementedApplicationAccessServer) CreateAPIKey(ctx context.Context, req *CreateApplicationAPIKeyRequest) (*APIKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAPIKey not implemented")
}
func (*UnimplementedApplicationAccessServer) ListAPIKeys(ctx context.Context, req *ListApplicationAPIKeysRequest) (*APIKeys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAPIKeys not implemented")
}
func (*UnimplementedApplicationAccessServer) GetAPIKey(ctx context.Context, req *GetApplicationAPIKeyRequest) (*APIKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAPIKey not implemented")
}
func (*UnimplementedApplicationAccessServer) UpdateAPIKey(ctx context.Context, req *UpdateApplicationAPIKeyRequest) (*APIKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAPIKey not implemented")
}
func (*UnimplementedApplicationAccessServer) GetCollaborator(ctx context.Context, req *GetApplicationCollaboratorRequest) (*GetCollaboratorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCollaborator not implemented")
}
func (*UnimplementedApplicationAccessServer) SetCollaborator(ctx context.Context, req *SetApplicationCollaboratorRequest) (*types.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetCollaborator not implemented")
}
func (*UnimplementedApplicationAccessServer) ListCollaborators(ctx context.Context, req *ListApplicationCollaboratorsRequest) (*Collaborators, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCollaborators not implemented")
}

func RegisterApplicationAccessServer(s *grpc.Server, srv ApplicationAccessServer) {
	s.RegisterService(&_ApplicationAccess_serviceDesc, srv)
}

func _ApplicationAccess_ListRights_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplicationIdentifiers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationAccessServer).ListRights(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationAccess/ListRights",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationAccessServer).ListRights(ctx, req.(*ApplicationIdentifiers))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationAccess_CreateAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateApplicationAPIKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationAccessServer).CreateAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationAccess/CreateAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationAccessServer).CreateAPIKey(ctx, req.(*CreateApplicationAPIKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationAccess_ListAPIKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListApplicationAPIKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationAccessServer).ListAPIKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationAccess/ListAPIKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationAccessServer).ListAPIKeys(ctx, req.(*ListApplicationAPIKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationAccess_GetAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetApplicationAPIKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationAccessServer).GetAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationAccess/GetAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationAccessServer).GetAPIKey(ctx, req.(*GetApplicationAPIKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationAccess_UpdateAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateApplicationAPIKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationAccessServer).UpdateAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationAccess/UpdateAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationAccessServer).UpdateAPIKey(ctx, req.(*UpdateApplicationAPIKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationAccess_GetCollaborator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetApplicationCollaboratorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationAccessServer).GetCollaborator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationAccess/GetCollaborator",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationAccessServer).GetCollaborator(ctx, req.(*GetApplicationCollaboratorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationAccess_SetCollaborator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetApplicationCollaboratorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationAccessServer).SetCollaborator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationAccess/SetCollaborator",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationAccessServer).SetCollaborator(ctx, req.(*SetApplicationCollaboratorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationAccess_ListCollaborators_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListApplicationCollaboratorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationAccessServer).ListCollaborators(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ttn.lorawan.v3.ApplicationAccess/ListCollaborators",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationAccessServer).ListCollaborators(ctx, req.(*ListApplicationCollaboratorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ApplicationAccess_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ttn.lorawan.v3.ApplicationAccess",
	HandlerType: (*ApplicationAccessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRights",
			Handler:    _ApplicationAccess_ListRights_Handler,
		},
		{
			MethodName: "CreateAPIKey",
			Handler:    _ApplicationAccess_CreateAPIKey_Handler,
		},
		{
			MethodName: "ListAPIKeys",
			Handler:    _ApplicationAccess_ListAPIKeys_Handler,
		},
		{
			MethodName: "GetAPIKey",
			Handler:    _ApplicationAccess_GetAPIKey_Handler,
		},
		{
			MethodName: "UpdateAPIKey",
			Handler:    _ApplicationAccess_UpdateAPIKey_Handler,
		},
		{
			MethodName: "GetCollaborator",
			Handler:    _ApplicationAccess_GetCollaborator_Handler,
		},
		{
			MethodName: "SetCollaborator",
			Handler:    _ApplicationAccess_SetCollaborator_Handler,
		},
		{
			MethodName: "ListCollaborators",
			Handler:    _ApplicationAccess_ListCollaborators_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lorawan-stack/api/application_services.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.0
// source: category.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CategoryServiceClient is the client API for CategoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CategoryServiceClient interface {
	FindAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListCategoryResponse, error)
	FindAllByTodoID(ctx context.Context, in *GetCategoryByTodoIDRequest, opts ...grpc.CallOption) (*ListCategoryResponse, error)
	FindByIDs(ctx context.Context, in *GetCategoryByIDsRequest, opts ...grpc.CallOption) (*ListCategoryResponse, error)
	FindByID(ctx context.Context, in *GetCategoryByIDRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error)
	Create(ctx context.Context, in *CreateCategoryRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error)
	Delete(ctx context.Context, in *GetCategoryByIDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type categoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCategoryServiceClient(cc grpc.ClientConnInterface) CategoryServiceClient {
	return &categoryServiceClient{cc}
}

func (c *categoryServiceClient) FindAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListCategoryResponse, error) {
	out := new(ListCategoryResponse)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/FindAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) FindAllByTodoID(ctx context.Context, in *GetCategoryByTodoIDRequest, opts ...grpc.CallOption) (*ListCategoryResponse, error) {
	out := new(ListCategoryResponse)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/FindAllByTodoID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) FindByIDs(ctx context.Context, in *GetCategoryByIDsRequest, opts ...grpc.CallOption) (*ListCategoryResponse, error) {
	out := new(ListCategoryResponse)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/FindByIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) FindByID(ctx context.Context, in *GetCategoryByIDRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error) {
	out := new(GetCategoryResponse)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/FindByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) Create(ctx context.Context, in *CreateCategoryRequest, opts ...grpc.CallOption) (*GetCategoryResponse, error) {
	out := new(GetCategoryResponse)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) Delete(ctx context.Context, in *GetCategoryByIDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.CategoryService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CategoryServiceServer is the server API for CategoryService service.
// All implementations must embed UnimplementedCategoryServiceServer
// for forward compatibility
type CategoryServiceServer interface {
	FindAll(context.Context, *emptypb.Empty) (*ListCategoryResponse, error)
	FindAllByTodoID(context.Context, *GetCategoryByTodoIDRequest) (*ListCategoryResponse, error)
	FindByIDs(context.Context, *GetCategoryByIDsRequest) (*ListCategoryResponse, error)
	FindByID(context.Context, *GetCategoryByIDRequest) (*GetCategoryResponse, error)
	Create(context.Context, *CreateCategoryRequest) (*GetCategoryResponse, error)
	Delete(context.Context, *GetCategoryByIDRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedCategoryServiceServer()
}

// UnimplementedCategoryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCategoryServiceServer struct {
}

func (UnimplementedCategoryServiceServer) FindAll(context.Context, *emptypb.Empty) (*ListCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAll not implemented")
}
func (UnimplementedCategoryServiceServer) FindAllByTodoID(context.Context, *GetCategoryByTodoIDRequest) (*ListCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAllByTodoID not implemented")
}
func (UnimplementedCategoryServiceServer) FindByIDs(context.Context, *GetCategoryByIDsRequest) (*ListCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByIDs not implemented")
}
func (UnimplementedCategoryServiceServer) FindByID(context.Context, *GetCategoryByIDRequest) (*GetCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByID not implemented")
}
func (UnimplementedCategoryServiceServer) Create(context.Context, *CreateCategoryRequest) (*GetCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedCategoryServiceServer) Delete(context.Context, *GetCategoryByIDRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedCategoryServiceServer) mustEmbedUnimplementedCategoryServiceServer() {}

// UnsafeCategoryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CategoryServiceServer will
// result in compilation errors.
type UnsafeCategoryServiceServer interface {
	mustEmbedUnimplementedCategoryServiceServer()
}

func RegisterCategoryServiceServer(s grpc.ServiceRegistrar, srv CategoryServiceServer) {
	s.RegisterService(&CategoryService_ServiceDesc, srv)
}

func _CategoryService_FindAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).FindAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/FindAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).FindAll(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_FindAllByTodoID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryByTodoIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).FindAllByTodoID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/FindAllByTodoID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).FindAllByTodoID(ctx, req.(*GetCategoryByTodoIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_FindByIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryByIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).FindByIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/FindByIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).FindByIDs(ctx, req.(*GetCategoryByIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_FindByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).FindByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/FindByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).FindByID(ctx, req.(*GetCategoryByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).Create(ctx, req.(*CreateCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.CategoryService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).Delete(ctx, req.(*GetCategoryByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CategoryService_ServiceDesc is the grpc.ServiceDesc for CategoryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CategoryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.CategoryService",
	HandlerType: (*CategoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAll",
			Handler:    _CategoryService_FindAll_Handler,
		},
		{
			MethodName: "FindAllByTodoID",
			Handler:    _CategoryService_FindAllByTodoID_Handler,
		},
		{
			MethodName: "FindByIDs",
			Handler:    _CategoryService_FindByIDs_Handler,
		},
		{
			MethodName: "FindByID",
			Handler:    _CategoryService_FindByID_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _CategoryService_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _CategoryService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "category.proto",
}

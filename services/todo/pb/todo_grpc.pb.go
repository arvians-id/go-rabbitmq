// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.0
// source: todo.proto

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

// TodoServiceClient is the client API for TodoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TodoServiceClient interface {
	FindAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTodoResponse, error)
	FindByUserIDs(ctx context.Context, in *GetTodoByUserIDsRequest, opts ...grpc.CallOption) (*ListTodoResponse, error)
	FindByID(ctx context.Context, in *GetTodoByIDRequest, opts ...grpc.CallOption) (*GetTodoResponse, error)
	Create(ctx context.Context, in *CreateTodoRequest, opts ...grpc.CallOption) (*GetTodoResponse, error)
	Update(ctx context.Context, in *UpdateTodoRequest, opts ...grpc.CallOption) (*GetTodoResponse, error)
	Delete(ctx context.Context, in *GetTodoByIDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type todoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTodoServiceClient(cc grpc.ClientConnInterface) TodoServiceClient {
	return &todoServiceClient{cc}
}

func (c *todoServiceClient) FindAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTodoResponse, error) {
	out := new(ListTodoResponse)
	err := c.cc.Invoke(ctx, "/proto.TodoService/FindAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) FindByUserIDs(ctx context.Context, in *GetTodoByUserIDsRequest, opts ...grpc.CallOption) (*ListTodoResponse, error) {
	out := new(ListTodoResponse)
	err := c.cc.Invoke(ctx, "/proto.TodoService/FindByUserIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) FindByID(ctx context.Context, in *GetTodoByIDRequest, opts ...grpc.CallOption) (*GetTodoResponse, error) {
	out := new(GetTodoResponse)
	err := c.cc.Invoke(ctx, "/proto.TodoService/FindByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) Create(ctx context.Context, in *CreateTodoRequest, opts ...grpc.CallOption) (*GetTodoResponse, error) {
	out := new(GetTodoResponse)
	err := c.cc.Invoke(ctx, "/proto.TodoService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) Update(ctx context.Context, in *UpdateTodoRequest, opts ...grpc.CallOption) (*GetTodoResponse, error) {
	out := new(GetTodoResponse)
	err := c.cc.Invoke(ctx, "/proto.TodoService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) Delete(ctx context.Context, in *GetTodoByIDRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.TodoService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TodoServiceServer is the server API for TodoService service.
// All implementations must embed UnimplementedTodoServiceServer
// for forward compatibility
type TodoServiceServer interface {
	FindAll(context.Context, *emptypb.Empty) (*ListTodoResponse, error)
	FindByUserIDs(context.Context, *GetTodoByUserIDsRequest) (*ListTodoResponse, error)
	FindByID(context.Context, *GetTodoByIDRequest) (*GetTodoResponse, error)
	Create(context.Context, *CreateTodoRequest) (*GetTodoResponse, error)
	Update(context.Context, *UpdateTodoRequest) (*GetTodoResponse, error)
	Delete(context.Context, *GetTodoByIDRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedTodoServiceServer()
}

// UnimplementedTodoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTodoServiceServer struct {
}

func (UnimplementedTodoServiceServer) FindAll(context.Context, *emptypb.Empty) (*ListTodoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAll not implemented")
}
func (UnimplementedTodoServiceServer) FindByUserIDs(context.Context, *GetTodoByUserIDsRequest) (*ListTodoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByUserIDs not implemented")
}
func (UnimplementedTodoServiceServer) FindByID(context.Context, *GetTodoByIDRequest) (*GetTodoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByID not implemented")
}
func (UnimplementedTodoServiceServer) Create(context.Context, *CreateTodoRequest) (*GetTodoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedTodoServiceServer) Update(context.Context, *UpdateTodoRequest) (*GetTodoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedTodoServiceServer) Delete(context.Context, *GetTodoByIDRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedTodoServiceServer) mustEmbedUnimplementedTodoServiceServer() {}

// UnsafeTodoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TodoServiceServer will
// result in compilation errors.
type UnsafeTodoServiceServer interface {
	mustEmbedUnimplementedTodoServiceServer()
}

func RegisterTodoServiceServer(s grpc.ServiceRegistrar, srv TodoServiceServer) {
	s.RegisterService(&TodoService_ServiceDesc, srv)
}

func _TodoService_FindAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).FindAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TodoService/FindAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).FindAll(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TodoService_FindByUserIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTodoByUserIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).FindByUserIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TodoService/FindByUserIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).FindByUserIDs(ctx, req.(*GetTodoByUserIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TodoService_FindByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTodoByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).FindByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TodoService/FindByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).FindByID(ctx, req.(*GetTodoByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TodoService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTodoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TodoService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).Create(ctx, req.(*CreateTodoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TodoService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTodoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TodoService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).Update(ctx, req.(*UpdateTodoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TodoService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTodoByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TodoService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).Delete(ctx, req.(*GetTodoByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TodoService_ServiceDesc is the grpc.ServiceDesc for TodoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TodoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TodoService",
	HandlerType: (*TodoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAll",
			Handler:    _TodoService_FindAll_Handler,
		},
		{
			MethodName: "FindByUserIDs",
			Handler:    _TodoService_FindByUserIDs_Handler,
		},
		{
			MethodName: "FindByID",
			Handler:    _TodoService_FindByID_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _TodoService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _TodoService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _TodoService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "todo.proto",
}

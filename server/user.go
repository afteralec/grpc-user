package server

import (
	"context"

	"github.com/afteralec/grpc-user/proto"
	"github.com/afteralec/grpc-user/services/user"
	"github.com/afteralec/grpc-user/services/user/passphrase"
	"github.com/afteralec/grpc-user/services/user/username"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedUserServer
	user *user.Service
}

func (s *server) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterReply, error) {
	if err := username.IsValid(in.Username); err != nil {
		// TODO: Implement Error Details
		return nil, status.Error(codes.InvalidArgument, "this error message is unimplemented")
	}
	if !passphrase.IsValid(in.Password) {
		// TODO: Implement Error Details
		return nil, status.Error(codes.InvalidArgument, "this error message is unimplemented")
	}

	uid, err := s.user.Register(in.Username, in.Password)
	if err != nil {
		// TODO: Implement Error Details
		return nil, status.Error(codes.Internal, "this error message is unimplemented")
	}
	return &proto.RegisterReply{Id: uid}, nil
}

func (s *server) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginReply, error) {
	uid, err := s.user.Authenticate(in.Username, in.Password)
	if err != nil {
		// TODO: Implement Error Details
		return nil, status.Error(codes.Unauthenticated, "this error message is unimplemented")
	}

	return &proto.LoginReply{Verified: true, Id: uid}, nil
}

func (s *server) UserSettings(ctx context.Context, in *proto.UserSettingsRequest) (*proto.UserSettingsReply, error) {
	settings, err := s.user.UserSettings(ctx, in.Uid)
	if err != nil {
		// TODO: Implement Error Details
		return nil, status.Error(codes.Internal, "this error message is unimplemented")
	}

	return &proto.UserSettingsReply{Id: settings.ID, Uid: settings.UID, Theme: settings.Theme}, nil
}

func (s *server) SetUserSettingsTheme(ctx context.Context, in *proto.SetUserSettingsThemeRequest) (*proto.SetUserSettingsThemeReply, error) {
	settings, err := s.user.SetUserSettingsTheme(ctx, in.Uid, in.Theme)
	if err != nil {
		// TODO: Implement Error Details
		return nil, status.Error(codes.Internal, "this error message is unimplemented")
	}

	return &proto.SetUserSettingsThemeReply{Id: settings.ID, Uid: settings.UID, Theme: settings.Theme}, nil
}

func (s *server) Users(ctx context.Context, in *proto.UsersRequest) (*proto.UsersReply, error) {
	users, err := s.user.Users(ctx)
	if err != nil {
		// TODO: Implement Error Details
		return nil, status.Error(codes.Internal, "this error message is unimplemented")
	}

	replyUsers := []*proto.UsersReplyUser{}
	for _, user := range users {
		replyUsers = append(replyUsers, &proto.UsersReplyUser{
			Id:           user.ID,
			Username:     user.Username,
			PrimaryEmail: "test@web.site",
		})
	}
	return &proto.UsersReply{Users: replyUsers}, nil
}

func (s *server) UserPermissionDefinitions(ctx context.Context, in *proto.UserPermissionDefinitionsRequest) (*proto.UserPermissionDefinitionsReply, error) {
	permissions := []*proto.UserPermissionDefinitionsReplyPermission{}
	for _, permission := range user.AllPermissions {
		permissions = append(permissions, &proto.UserPermissionDefinitionsReplyPermission{
			Name:     permission.Name,
			Title:    permission.Title,
			About:    permission.About,
			Category: permission.Category,
		})
	}

	return &proto.UserPermissionDefinitionsReply{Permissions: permissions}, nil
}

func (s *server) UserPermissions(ctx context.Context, in *proto.UserPermissionsRequest) (*proto.UserPermissionsReply, error) {
	permissions, err := s.user.UserPermissions(ctx, in.Uid)
	if err != nil {
		// TODO: Implement Error Details
		return nil, status.Error(codes.Internal, "this error message is unimplemented")
	}

	names := []string{}
	for _, permission := range permissions {
		names = append(names, permission.Name)
	}

	return &proto.UserPermissionsReply{Uid: in.Uid, Names: names}, nil
}

func (s *server) GrantUserPermission(ctx context.Context, in *proto.GrantUserPermissionRequest) (*proto.GrantUserPermissionReply, error) {
	if user.IsRootPermission(in.Name) {
		// TODO: Implement Error Details
		return nil, status.Error(codes.PermissionDenied, "this error message is unimplemented")
	}
	id, err := s.user.GrantUserPermission(ctx, in.Uid, in.Iuid, in.Name)
	if err != nil {
		// TODO: Implement Error Details
		return nil, status.Error(codes.Internal, "this error message is unimplemented")
	}

	return &proto.GrantUserPermissionReply{Id: id}, nil
}

func (s *server) RevokeUserPermission(ctx context.Context, in *proto.RevokeUserPermissionRequest) (*proto.RevokeUserPermissionReply, error) {
	if user.IsRootPermission(in.Name) {
		// TODO: Implement Error Details
		return nil, status.Error(codes.PermissionDenied, "this error message is unimplemented")
	}
	id, err := s.user.RevokeUserPermission(ctx, in.Uid, in.Iuid, in.Name)
	if err != nil {
		// TODO: Implement Error Details
		return nil, status.Error(codes.Internal, "this error message is unimplemented")
	}

	return &proto.RevokeUserPermissionReply{Id: id}, nil
}

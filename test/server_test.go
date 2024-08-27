package test

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/afteralec/grpc-user/proto"
	"github.com/afteralec/grpc-user/services/user"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	TestUsername     = "testify"
	TestRootUsername = "tested"
	TestPassword     = "T3sted_tested"
)

func TestServer(t *testing.T) {
	newClient := func(addr string) (pb.MirrorClient, func()) {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		closeClient := func() {
			err := conn.Close()
			require.NoError(t, err)
		}

		client := pb.NewMirrorClient(conn)

		return client, closeClient
	}

	compose, err := tc.NewDockerComposeWith(tc.WithStackFiles("../compose.yml", "../test.compose.yml"))
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal))
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	err = compose.WaitForService("mirror", wait.ForExposedPort()).Up(ctx)
	require.NoError(t, err)

	server, err := compose.ServiceContainer(ctx, "mirror")
	require.NoError(t, err, "get server ServiceContainer")

	host, err := server.Host(ctx)
	require.NoError(t, err)
	port, err := server.MappedPort(ctx, "8009/tcp")
	require.NoError(t, err)
	addr := fmt.Sprintf("%s:%s", host, port.Port())

	client, closeClient := newClient(addr)
	t.Cleanup(closeClient)

	t.Run("Register Success", func(t *testing.T) {
		t.Parallel()
		ctx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer", "authorization", "some-secret-token")
		reply, err := client.Register(ctx, &pb.RegisterRequest{
			Username: TestUsername,
			Password: TestPassword,
		})
		require.NoError(t, err)
		require.NotNil(t, reply)
		require.Greater(t, reply.Id, int64(0))
	})

	type invalidInputTestcase struct {
		name        string
		username    string
		password    string
		shouldError bool
	}
	invalidInputTestcases := []invalidInputTestcase{{"Username too long", "testfieduntiltheendoftime", "T3sted_tested", true}}
	for _, tc := range invalidInputTestcases {
		tc := tc
		t.Run("Register Invalid Username/Password: "+tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer", "authorization", "some-secret-token")
			reply, err := client.Register(ctx, &pb.RegisterRequest{
				Username: tc.username,
				Password: tc.password,
			})
			if tc.shouldError {
				require.Error(t, err, "Register Invalid Username/Password")
				require.Nil(t, reply)
			} else {
				require.NoError(t, err, "Register Invalid Username/Password")
				require.NotNil(t, reply)
				require.Greater(t, reply.Id, int64(0))
			}
		})
	}

	t.Run("Register Conflict", func(t *testing.T) {
		t.Parallel()
		ctx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer", "authorization", "some-secret-token")

		reply, err := client.Register(ctx, &pb.RegisterRequest{
			Username: "testific",
			Password: "T3sted_tested",
		})
		require.NoError(t, err, "Register Conflict")
		require.NotNil(t, reply)

		reply, err = client.Register(ctx, &pb.RegisterRequest{
			Username: "testific",
			Password: "T3sted_tested",
		})
		require.Error(t, err, "Register Conflict")
		require.Nil(t, reply)
	})

	t.Run("Login Success", func(t *testing.T) {
		t.Parallel()
		registerReply, err := client.Register(ctx, &pb.RegisterRequest{
			Username: "testify",
			Password: "T3sted_tested",
		})
		require.NoError(t, err, "Login Success")
		require.NotNil(t, registerReply)

		loginReply, err := client.Login(ctx, &pb.LoginRequest{
			Username: "testify",
			Password: "T3sted_tested",
		})
		require.NoError(t, err, "Login Success")
		require.NotNil(t, loginReply)
		require.True(t, loginReply.Verified)
		require.Equal(t, registerReply.Id, loginReply.Id)
	})

	t.Run("Login Invalid Username", func(t *testing.T) {
		t.Parallel()
		registerReply, err := client.Register(ctx, &pb.RegisterRequest{
			Username: "testli",
			Password: "T3sted_tested",
		})
		require.NoError(t, err, "Login Success")
		require.NotNil(t, registerReply)

		loginReply, err := client.Login(ctx, &pb.LoginRequest{
			Username: "testli",
			Password: "T3sted_tested",
		})
		require.NoError(t, err)
		require.NotNil(t, loginReply)
		require.True(t, loginReply.Verified)
		require.Equal(t, registerReply.Id, loginReply.Id)

		loginReply, err = client.Login(ctx, &pb.LoginRequest{
			Username: "testly",
			Password: "T3sted_tested",
		})
		require.Error(t, err)
		require.Equal(t, codes.Unauthenticated, status.Code(err))
		require.Nil(t, loginReply)
	})

	t.Run("Login Invalid Password", func(t *testing.T) {
		t.Parallel()
		registerReply, err := client.Register(ctx, &pb.RegisterRequest{
			Username: "testla",
			Password: "T3sted_tested",
		})
		require.NoError(t, err)
		require.NotNil(t, registerReply)

		loginReply, err := client.Login(ctx, &pb.LoginRequest{
			Username: "testla",
			Password: "T3sted_tested",
		})
		require.NoError(t, err)
		require.NotNil(t, loginReply)
		require.True(t, loginReply.Verified)
		require.Equal(t, registerReply.Id, loginReply.Id)

		loginReply, err = client.Login(ctx, &pb.LoginRequest{
			Username: "testla",
			Password: "T3sted_tasted",
		})
		require.Error(t, err)
		require.Equal(t, codes.Unauthenticated, status.Code(err))
		require.Nil(t, loginReply)
	})

	t.Run("User Settings Success", func(t *testing.T) {
		t.Parallel()
		registerReply, err := client.Register(ctx, &pb.RegisterRequest{
			Username: "testls",
			Password: "T3sted_tested",
		})
		require.NoError(t, err)
		require.NotNil(t, registerReply)

		userSettingsReply, err := client.UserSettings(ctx, &pb.UserSettingsRequest{
			Uid: registerReply.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, userSettingsReply)
		require.Equal(t, registerReply.Id, userSettingsReply.Uid)
		require.Equal(t, user.ThemeDefault, userSettingsReply.Theme)
	})

	t.Run("Set User Settings Theme Success", func(t *testing.T) {
		t.Parallel()
		registerReply, err := client.Register(ctx, &pb.RegisterRequest{
			Username: "testll",
			Password: "T3sted_tested",
		})
		require.NoError(t, err)
		require.NotNil(t, registerReply)

		userSettingsReply, err := client.UserSettings(ctx, &pb.UserSettingsRequest{
			Uid: registerReply.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, userSettingsReply)
		require.Equal(t, registerReply.Id, userSettingsReply.Uid)
		require.Equal(t, user.ThemeDefault, userSettingsReply.Theme)

		theme, err := user.OtherTheme(user.ThemeDefault)
		require.NoError(t, err)
		setUserSettingsThemeReply, err := client.SetUserSettingsTheme(ctx, &pb.SetUserSettingsThemeRequest{
			Uid:   registerReply.Id,
			Theme: theme,
		})
		require.NoError(t, err)
		require.NotNil(t, setUserSettingsThemeReply)
		require.Equal(t, registerReply.Id, setUserSettingsThemeReply.Uid)
		require.Equal(t, theme, setUserSettingsThemeReply.Theme)

		userSettingsReply, err = client.UserSettings(ctx, &pb.UserSettingsRequest{
			Uid: registerReply.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, userSettingsReply)
		require.Equal(t, registerReply.Id, userSettingsReply.Uid)
		require.Equal(t, userSettingsReply.Theme, user.ThemeLight)
	})
}

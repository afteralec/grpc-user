package user

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/afteralec/grpc-user/db/query"
	"github.com/afteralec/grpc-user/services/user/passphrase"
	"github.com/afteralec/grpc-user/services/user/username"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
)

type Service struct {
	db     *sql.DB
	query  *query.Queries
	config *viper.Viper
}

func New(db *sql.DB, opts ...func(s *Service) error) (Service, error) {
	if db == nil {
		return Service{}, errors.New("cannot instantiate without a database connection")
	}
	// TODO: Get sensible defaults for this config
	service := Service{db: db, query: query.New(db), config: viper.New()}
	for _, opt := range opts {
		if err := opt(&service); err != nil {
			return Service{}, err
		}
	}
	if err := username.IsValid(service.config.GetString("root_username")); err != nil {
		return Service{}, err
	}
	return service, nil
}

func (s *Service) SyncRootPermissions(ctx context.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := s.query.WithTx(tx)

	if err := revokeAllRootUserPermissions(ctx, qtx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	tx, err = s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx = s.query.WithTx(tx)

	rootUsername := s.config.GetString("root_username")

	exists, err := s.userExistsWithUsername(ctx, rootUsername)
	if err != nil {
		return err
	}

	var uid int64
	if exists {
		u, err := qtx.GetUserByUsername(ctx, s.config.GetString("root_username"))
		if err != nil {
			return err
		}
		uid = u.ID
	} else {
		uid, err = s.Register(s.config.GetString("root_username"), s.config.GetString("root_passphrase"))
		if err != nil {
			return err
		}
	}

	if err := grantRootUserPermissions(ctx, qtx, uid); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Register(u, pass string) (int64, error) {
	hash, err := passphrase.Hash(pass, passphrase.NewParams())
	if err != nil {
		return 0, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	qtx := s.query.WithTx(tx)

	r, err := qtx.CreateUser(context.Background(), query.CreateUserParams{
		Username: u,
		PwHash:   hash,
	})
	if err != nil {
		return 0, err
		// TODO: Create a parsing layer for SQLite errors
		// if sqliteErr, ok := err.(sqlite3.Error); ok {
		// 	if sqliteErr.Code == sqlite3.Err {
		// 		fmt.Println("busy")
		// 	}
		// }
	}
	uid, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := qtx.CreateUserSettings(context.Background(), query.CreateUserSettingsParams{
		Theme: ThemeDefault,
		UID:   uid,
	}); err != nil {
		return 0, err
	}

	if u == s.config.GetString("root_username") {
		if err := revokeAllRootUserPermissions(context.Background(), qtx); err != nil {
			log.Printf("revoke all root permissions err: %v", err)
			return 0, err
		}
		if err := grantRootUserPermissions(context.Background(), qtx, uid); err != nil {
			log.Printf("grant all root permissions err: %v", err)
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return uid, nil
}

type UnauthenticatedError struct{}

func (e *UnauthenticatedError) Error() string {
	return "could not authenticate this username and password"
}

func (s *Service) Authenticate(u, pw string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	qtx := s.query.WithTx(tx)

	p, err := qtx.GetUserByUsername(context.Background(), u)
	if err != nil {
		rand.Seed(uint64(time.Now().UnixNano()))
		n := rand.Intn(2) + 2
		time.Sleep(time.Duration(n) * time.Second)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	ok, err := passphrase.Verify(pw, p.PwHash)
	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, &UnauthenticatedError{}
	}

	return p.ID, nil
}

func (s *Service) UserSettings(ctx context.Context, uid int64) (*query.UserSetting, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	qtx := s.query.WithTx(tx)

	settings, err := qtx.GetUserSettings(ctx, uid)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &settings, nil
}

func (s *Service) SetUserSettingsTheme(ctx context.Context, uid int64, theme string) (*query.UserSetting, error) {
	// TODO: Discrete error type
	if theme != ThemeLight && theme != ThemeDark {
		return nil, errors.New("invalid theme value")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	qtx := s.query.WithTx(tx)

	if err := qtx.UpdateUserSettingsTheme(ctx, query.UpdateUserSettingsThemeParams{
		UID:   uid,
		Theme: theme,
	}); err != nil {
		return nil, err
	}

	settings, err := qtx.GetUserSettings(ctx, uid)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &settings, nil
}

func (s *Service) Users(ctx context.Context) ([]query.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return []query.User{}, err
	}
	defer tx.Rollback()
	qtx := s.query.WithTx(tx)

	users, err := qtx.ListUsers(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return []query.User{}, nil
		}
		return []query.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return []query.User{}, err
	}

	return users, nil
}

func (s *Service) UserPermissions(ctx context.Context, uid int64) ([]query.UserPermission, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return []query.UserPermission{}, err
	}
	defer tx.Rollback()
	qtx := s.query.WithTx(tx)

	permissions, err := userPermissions(ctx, qtx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return []query.UserPermission{}, nil
		}
		return []query.UserPermission{}, err
	}

	if err := tx.Commit(); err != nil {
		return []query.UserPermission{}, err
	}

	return permissions, nil
}

func (s *Service) GrantUserPermission(ctx context.Context, uid, iuid int64, name string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	qtx := s.query.WithTx(tx)

	issuerPermissionRecords, err := userPermissions(ctx, qtx, iuid)
	if err != nil {
		return 0, err
	}
	issuerPermissions := NewPermissions(iuid, issuerPermissionRecords)
	if !issuerPermissions.CanGrant(name) {
		// TODO: Create a discrete error here
		return 0, errors.New("this issuer cannot grant this permission")
	}

	id, err := grantUserPermission(ctx, qtx, uid, iuid, name)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) RevokeUserPermission(ctx context.Context, uid, iuid int64, name string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	qtx := s.query.WithTx(tx)

	issuerPermissionRecords, err := userPermissions(ctx, qtx, iuid)
	if err != nil {
		return 0, err
	}
	issuerPermissions := NewPermissions(iuid, issuerPermissionRecords)
	if !issuerPermissions.CanRevoke(name) {
		// TODO: Create a discrete error here
		return 0, errors.New("this issuer cannot revoke this permission")
	}

	id, err := revokeUserPermission(ctx, qtx, uid, iuid, name)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) userExistsWithUsername(ctx context.Context, u string) (bool, error) {
	_, err := s.query.GetUserByUsername(ctx, u)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

package user

import "github.com/afteralec/grpc-user/db/query"

type Permission struct {
	Name     string
	Title    string
	About    string
	Category string
}

var PermissionGrantAll Permission = Permission{
	Name:     "grant-all",
	Title:    "Grant All Permissions",
	About:    "The root permission. Only one person should have this at a time.",
	Category: "Root",
}

var PermissionRevokeAll Permission = Permission{
	Name:     "revoke-all",
	Title:    "Revoke All Permissions",
	About:    "The root revocation permission. Only one person should have this at a time.",
	Category: "Root",
}

var PermissionReviewCharacterApplications Permission = Permission{
	Name:     "review-character-applications",
	Title:    "Review Character Applications",
	About:    "Enable this user to review Character Applications.",
	Category: "Reviewer",
}

var PermissionViewAllRooms Permission = Permission{
	Name:     "view-all-rooms",
	Title:    "View All Rooms",
	About:    "The permission to view (but not edit) all room data.",
	Category: "Room",
}

var PermissionCreateRoom Permission = Permission{
	Name:     "create-room",
	Title:    "Create Room",
	About:    "Create a new room, but not connect it to the grid.",
	Category: "Room",
}

var PermissionViewAllActorImages Permission = Permission{
	Name:     "view-all-actor-images",
	Title:    "View All Actor Images",
	About:    "View all Actor Images, i.e. in the main Actor Images list.",
	Category: "Actor",
}

var PermissionCreateActorImage Permission = Permission{
	Name:     "create-actor-image",
	Title:    "Create Actor Image",
	About:    "Create new actor via creating new Actor Images",
	Category: "Actor",
}

var PermissionCreateChangelog Permission = Permission{
	Name:     "create-changelog",
	Title:    "Create Changelogs",
	About:    "Draft and edit changelogs prior to release",
	Category: "Changelog",
}

var PermissionReleaseChangelog Permission = Permission{
	Name:     "release-changelog",
	Title:    "Release Changelogs",
	About:    "Release a changelog. Once released, it cannot be edited or revoked",
	Category: "Changelog",
}

var PermissionRevokeChangelog Permission = Permission{
	Name:     "revoke-changelog",
	Title:    "Revoke Changelogs",
	About:    "Revoke a changelog. This is an emergency measure.",
	Category: "Changelog",
}

var AllPermissions []Permission = []Permission{
	PermissionGrantAll,
	PermissionRevokeAll,
	PermissionReviewCharacterApplications,
	PermissionViewAllRooms,
	PermissionCreateRoom,
	PermissionViewAllActorImages,
	PermissionCreateActorImage,
	PermissionCreateChangelog,
	PermissionReleaseChangelog,
	PermissionRevokeChangelog,
}

var RootPermissions []Permission = []Permission{
	PermissionGrantAll,
	PermissionRevokeAll,
}

var NonRootPermissions []Permission = nonRootPermissions(AllPermissions, RootPermissions)

func nonRootPermissions(all []Permission, root []Permission) []Permission {
	rootByName := permissionsByName(root)

	nonRootPermissions := []Permission{}
	for _, permission := range all {
		_, ok := rootByName[permission.Name]
		if ok {
			continue
		}

		nonRootPermissions = append(nonRootPermissions, permission)
	}

	return nonRootPermissions
}

func permissionsByName(permissions []Permission) map[string]Permission {
	permissionsbyname := make(map[string]Permission)
	for _, permission := range permissions {
		permissionsbyname[permission.Name] = permission
	}
	return permissionsbyname
}

var (
	AllPermissionsByName     = permissionsByName(AllPermissions)
	RootPermissionsByName    = permissionsByName(RootPermissions)
	NonRootPermissionsByName = permissionsByName(NonRootPermissions)
)

type Permissions struct {
	inner map[string]bool
	list  []string
	uid   int64
}

func NewPermissions(uid int64, perms []query.UserPermission) Permissions {
	filtered := []query.UserPermission{}
	for _, perm := range perms {
		if IsValidPermissionName(perm.Name) {
			filtered = append(filtered, perm)
		}
	}
	list := []string{}
	for _, perm := range filtered {
		list = append(list, perm.Name)
	}
	permissionsMap := map[string]bool{}
	for _, perm := range filtered {
		permissionsMap[perm.Name] = true
	}
	return Permissions{
		uid:   uid,
		list:  list,
		inner: permissionsMap,
	}
}

func (p *Permissions) Has(name string) bool {
	_, ok := p.inner[name]
	return ok
}

func (p *Permissions) HasPermissionInSet(set []string) bool {
	for _, perm := range set {
		_, ok := p.inner[perm]
		if ok {
			return true
		}
	}
	return false
}

func (p *Permissions) HasAllPermissionsInSet(set []string) bool {
	for _, perm := range set {
		_, ok := p.inner[perm]
		if !ok {
			return false
		}
	}
	return true
}

func (p *Permissions) CanGrant(name string) bool {
	if !IsValidPermissionName(name) {
		return false
	}

	_, ok := RootPermissionsByName[name]
	if ok {
		return false
	}

	return p.Has(PermissionGrantAll.Name)
}

func (p *Permissions) CanRevoke(name string) bool {
	if !IsValidPermissionName(name) {
		return false
	}

	_, ok := RootPermissionsByName[name]
	if ok {
		return false
	}

	return p.Has(PermissionRevokeAll.Name)
}

func IsValidPermissionName(name string) bool {
	_, ok := AllPermissionsByName[name]
	return ok
}

func IsRootPermission(name string) bool {
	if !IsValidPermissionName(name) {
		return false
	}

	_, ok := RootPermissionsByName[name]
	return ok
}

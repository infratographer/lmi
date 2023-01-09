package storage

import (
	"context"

	apiv1 "github.com/infratographer/lmi/api/v1"
)

type Storage interface {
	GetAssignments(c context.Context, params *apiv1.GetAssignmentsParams) ([]*apiv1.Assignment, error)

	GetPermissions(c context.Context, params *apiv1.GetPermissionsParams) ([]*apiv1.Permission, error)

	GetRoles(c context.Context) ([]*apiv1.RoleInfo, error)

	CreateRole(c context.Context, role apiv1.NewRole) (*apiv1.Role, error)

	DeleteRole(c context.Context, id apiv1.EntityID) error

	GetRole(c context.Context, id apiv1.EntityID) (*apiv1.Role, error)

	UpdateRole(c context.Context, role *apiv1.Role) (*apiv1.Role, error)

	RemoveRoleAssignment(c context.Context, a apiv1.Assignment) error

	GetRoleAssignments(c context.Context, roleID apiv1.EntityID) ([]*apiv1.Assignment, error)

	AssignRole(c context.Context, roleID apiv1.EntityID, assignment apiv1.NewRoleAssignment) error

	RemoveRolePermission(c context.Context, id apiv1.EntityID, targetID apiv1.PermissionIdentifier) error

	GetRolePermissions(c context.Context, id apiv1.EntityID) ([]*apiv1.Permission, error)

	AddRolePermission(c context.Context, id apiv1.EntityID, targetID apiv1.PermissionIdentifier) error
}

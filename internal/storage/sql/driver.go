package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	apiv1 "github.com/infratographer/lmi/api/v1"
	"github.com/infratographer/lmi/internal/storage"
	"github.com/infratographer/lmi/internal/storage/sql/models"
)

type sqlDriver struct {
	db *sql.DB
}

// ensure we implement storage interface.
var _ storage.Storage = (*sqlDriver)(nil)

func NewSQLDriver(db *sql.DB) storage.Storage {
	return &sqlDriver{
		db: db,
	}
}

func (drv *sqlDriver) GetAssignments(
	c context.Context,
	params *apiv1.GetAssignmentsParams,
) ([]*apiv1.Assignment, error) {
	as, err := models.RoleAssignments(buildGetAssignmentsQuery(params)...).All(c, drv.db)
	if err != nil {
		return nil, fmt.Errorf("couldn't get role assignments: %w", err)
	}

	assignments := make([]*apiv1.Assignment, len(as))
	for i, a := range as {
		roleID, err := apiv1.ParseEntityID(a.RoleID)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse role ID for role assignment %s", a.ID)
		}

		assignments[i] = &apiv1.Assignment{
			Subject: a.SubjectID,
			Scope:   a.Scope,
			Role:    roleID,
		}
	}

	return assignments, nil
}

func buildGetAssignmentsQuery(params *apiv1.GetAssignmentsParams) []qm.QueryMod {
	mods := []qm.QueryMod{}
	if params != nil {
		if params.Subject != "" {
			mods = append(mods, qm.Where(models.RoleAssignmentColumns.SubjectID+"=?", params.Subject))
		}
		if params.Scope != "" {
			if len(mods) > 0 {
				mods = append(mods, qm.And(models.RoleAssignmentColumns.Scope+"=?", params.Scope))
			}

			mods = append(mods, qm.Where(models.RoleAssignmentColumns.Scope+"=?", params.Scope))
		}
		if params.Role != nil {
			if len(mods) > 0 {
				mods = append(mods, qm.And(models.RoleAssignmentColumns.RoleID+"=?", params.Role))
			}

			mods = append(mods, qm.Where(models.RoleAssignmentColumns.RoleID+"=?", params.Role))
		}
	}

	return mods
}

func (drv *sqlDriver) GetPermissions(
	c context.Context,
	params *apiv1.GetPermissionsParams,
) ([]*apiv1.Permission, error) {
	perms, err := models.Permissions(buildPermissionsQuery(params)...).All(c, drv.db)
	if err != nil {
		return nil, fmt.Errorf("couldn't get permissions: %w", err)
	}

	permissions := make([]*apiv1.Permission, len(perms))
	for i, p := range perms {
		permissions[i] = &apiv1.Permission{
			Target:      p.Target,
			Description: &p.Description,
		}
	}

	return permissions, nil
}

func buildPermissionsQuery(params *apiv1.GetPermissionsParams) []qm.QueryMod {
	mods := []qm.QueryMod{}
	if params != nil {
		if params.Target != nil {
			mods = append(mods, models.PermissionWhere.Target.EQ(*params.Target))
		}
	}

	return mods
}

func (drv *sqlDriver) GetRoles(c context.Context) ([]*apiv1.RoleInfo, error) {
	roles, err := models.Roles().All(c, drv.db)
	if err != nil {
		return nil, fmt.Errorf("couldn't get roles: %w", err)
	}

	rolesOut := make([]*apiv1.RoleInfo, len(roles))
	for i, r := range roles {
		roleID, err := apiv1.ParseEntityID(r.ID)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse role ID %s: %w", r.ID, err)
		}

		rolesOut[i] = &apiv1.RoleInfo{
			Id:          roleID,
			Name:        r.Name,
			Description: &r.Description,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	}

	return rolesOut, nil
}

func (drv *sqlDriver) CreateRole(c context.Context, newRole apiv1.NewRole) (*apiv1.Role, error) {
	r := &models.Role{
		Name: newRole.Name,
	}

	if newRole.Description != nil {
		r.Description = *newRole.Description
	}

	if err := r.Insert(c, drv.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("couldn't create role: %w", err)
	}

	roleID, err := apiv1.ParseEntityID(r.ID)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse role ID %s: %w", r.ID, err)
	}

	return &apiv1.Role{
		Id:          roleID,
		Name:        r.Name,
		Description: &r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}, nil
}

func (drv *sqlDriver) DeleteRole(c context.Context, id apiv1.EntityID) error {
	r, err := models.FindRole(c, drv.db, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return fmt.Errorf("couldn't find role: %w", err)
	}

	if _, err := r.Delete(c, drv.db); err != nil {
		return fmt.Errorf("couldn't delete role: %w", err)
	}

	return nil
}

func (drv *sqlDriver) GetRole(c context.Context, id apiv1.EntityID) (*apiv1.Role, error) {
	r, err := models.FindRole(c, drv.db, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("couldn't find role: %w", err)
	}

	roleID, err := apiv1.ParseEntityID(r.ID)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse role ID %s: %w", r.ID, err)
	}

	// Get role permissions
	rp, err := r.TargetPermissions().All(c, drv.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("couldn't get role permissions: %w", err)
	}

	perms := make([]apiv1.Permission, len(rp))
	for i, p := range rp {
		perms[i] = apiv1.Permission{
			Target:      p.Target,
			Description: &p.Description,
		}
	}

	return &apiv1.Role{
		Id:          roleID,
		Name:        r.Name,
		Description: &r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		Permissions: &perms,
	}, nil
}

func (drv *sqlDriver) UpdateRole(c context.Context, role *apiv1.Role) (*apiv1.Role, error) {
	r, err := models.FindRole(c, drv.db, role.Id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("couldn't find role: %w", err)
	}

	r.Name = role.Name
	if role.Description != nil {
		r.Description = *role.Description
	}

	if _, err := r.Update(c, drv.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("couldn't update role: %w", err)
	}

	roleID, err := apiv1.ParseEntityID(r.ID)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse role ID %s: %w", r.ID, err)
	}

	return &apiv1.Role{
		Id:          roleID,
		Name:        r.Name,
		Description: &r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}, nil
}

func (drv *sqlDriver) RemoveRoleAssignment(c context.Context, a apiv1.Assignment) error {
	r, err := models.FindRole(c, drv.db, a.Role.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return fmt.Errorf("couldn't find role: %w", err)
	}

	ra, err := r.RoleAssignments().One(c, drv.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = storage.ErrNotFound
		}
		return fmt.Errorf("couldn't get role assignment: %w", err)
	}

	if _, err := ra.Delete(c, drv.db); err != nil {
		return fmt.Errorf("couldn't delete role assignment: %w", err)
	}

	return nil
}

func (drv *sqlDriver) GetRoleAssignments(c context.Context, roleID apiv1.EntityID) ([]*apiv1.Assignment, error) {
	r, err := models.FindRole(c, drv.db, roleID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("couldn't find role: %w", err)
	}

	// Get role assignments
	ra, err := r.RoleAssignments().All(c, drv.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = storage.ErrNotFound
		}
		return nil, fmt.Errorf("couldn't get role assignments: %w", err)
	}

	assignments := make([]*apiv1.Assignment, len(ra))
	for i, a := range ra {
		assignments[i] = &apiv1.Assignment{
			Role:    roleID,
			Subject: a.SubjectID,
			Scope:   a.Scope,
		}
	}

	return assignments, nil
}

func (drv *sqlDriver) AssignRole(c context.Context, roleID apiv1.EntityID, assignment apiv1.NewRoleAssignment) error {
	r, err := models.FindRole(c, drv.db, roleID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return fmt.Errorf("couldn't find role: %w", err)
	}

	// verify if role assignment already exists
	exists, err := r.RoleAssignments(
		qm.Where(models.RoleAssignmentColumns.SubjectID+"=?", assignment.Subject),
		qm.And(models.RoleAssignmentColumns.Scope+"=?", assignment.Scope),
	).Exists(c, drv.db)
	if err != nil {
		return fmt.Errorf("couldn't check if role assignment exists: %w", err)
	}

	if exists {
		return nil
	}

	ra := models.RoleAssignment{
		RoleID:    r.ID,
		SubjectID: assignment.Subject,
		Scope:     assignment.Scope,
	}

	err = ra.Insert(c, drv.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("couldn't create role assignment: %w", err)
	}

	return nil
}

func (drv *sqlDriver) RemoveRolePermission(
	c context.Context,
	id apiv1.EntityID,
	targetID apiv1.PermissionIdentifier,
) error {
	r, err := models.FindRole(c, drv.db, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return fmt.Errorf("couldn't find role: %w", err)
	}

	// get permission
	p, err := r.TargetPermissions(
		qm.Where(models.PermissionColumns.Target+"=?", targetID.Target),
	).One(c, drv.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return fmt.Errorf("couldn't get role permission: %w", err)
	}

	err = r.RemoveTargetPermissions(c, drv.db, p)
	if err != nil {
		return fmt.Errorf("couldn't remove role permission: %w", err)
	}

	return nil
}

func (drv *sqlDriver) GetRolePermissions(c context.Context, id apiv1.EntityID) ([]*apiv1.Permission, error) {
	r, err := models.FindRole(c, drv.db, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("couldn't find role: %w", err)
	}

	// get permissions
	tp, err := r.TargetPermissions().All(c, drv.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("couldn't get role permissions: %w", err)
	}

	perms := make([]*apiv1.Permission, len(tp))
	for i, p := range tp {
		perms[i] = &apiv1.Permission{
			Target:      p.Target,
			Description: &p.Description,
		}
	}

	return perms, nil
}

func (drv *sqlDriver) AddRolePermission(
	c context.Context,
	id apiv1.EntityID,
	targetID apiv1.PermissionIdentifier,
) error {
	r, err := models.FindRole(c, drv.db, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return fmt.Errorf("couldn't find role: %w", err)
	}

	// Get permission
	perm, err := models.FindPermission(c, drv.db, targetID.Target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = storage.ErrNotFound
		}
		return fmt.Errorf("couldn't find permission: %w", err)
	}

	// check if role already has permission
	permExists, err := r.TargetPermissions(
		qm.Where(models.PermissionColumns.Target+"=?", targetID.Target),
	).Exists(c, drv.db)
	if err != nil {
		return fmt.Errorf("couldn't check if role has permission: %w", err)
	}

	if permExists {
		return nil
	}

	// add permission to role
	err = r.AddTargetPermissions(c, drv.db, true, perm)
	if err != nil {
		return fmt.Errorf("couldn't add permission to role: %w", err)
	}

	return nil
}

package httpsrv

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"

	apiv1 "github.com/infratographer/lmi/api/v1"
	"github.com/infratographer/lmi/internal/storage"
)

func (rtr *Router) ErrorHandler(gctx *gin.Context, err error, statusCode int) {
	gctx.JSON(statusCode, gin.H{"msg": err.Error()})
}

func (rtr *Router) ErrorChooser(gctx *gin.Context, err error) {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		rtr.ErrorHandler(gctx, err, http.StatusNotFound)
	default:
		rtr.ErrorHandler(gctx, err, http.StatusInternalServerError)
	}
}

func (rtr *Router) GetAssignments(c *gin.Context) {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params apiv1.GetAssignmentsParams

	// ------------- Required query parameter "subject" -------------

	if paramValue := c.Query("subject"); paramValue != "" {
	} else {
		rtr.ErrorHandler(c, fmt.Errorf("query argument subject is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "subject", c.Request.URL.Query(), &params.Subject)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter subject: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Required query parameter "scope" -------------

	if paramValue := c.Query("scope"); paramValue != "" {
	} else {
		rtr.ErrorHandler(c, fmt.Errorf("query argument scope is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "scope", c.Request.URL.Query(), &params.Scope)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter scope: %w", err), http.StatusBadRequest)
		return
	}

	as, err := rtr.store.GetAssignments(c, &params)
	if err != nil {
		rtr.ErrorChooser(c, err)
		return
	}

	c.JSON(http.StatusOK, as)
}

func (rtr *Router) GetPermissions(c *gin.Context) {
	// Parameter object where we will unmarshal all parameters from the context
	var params apiv1.GetPermissionsParams

	// ------------- Required query parameter "target" -------------

	if paramValue := c.Query("target"); paramValue != "" {
	} else {
		rtr.ErrorHandler(c, fmt.Errorf("query argument target is required, but not found"), http.StatusBadRequest)
		return
	}

	err := runtime.BindQueryParameter("form", true, true, "target", c.Request.URL.Query(), &params.Target)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter target: %w", err), http.StatusBadRequest)
		return
	}

	perms, err := rtr.store.GetPermissions(c, &params)
	if err != nil {
		rtr.ErrorChooser(c, err)
		return
	}

	c.JSON(http.StatusOK, perms)
}

func (rtr *Router) GetRoles(c *gin.Context) {
	roles, err := rtr.store.GetRoles(c)
	if err != nil {
		rtr.ErrorChooser(c, err)
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (rtr *Router) CreateRole(c *gin.Context) {
	newRole := apiv1.NewRole{}

	if err := c.BindJSON(&newRole); err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for new role: %w", err), http.StatusBadRequest)
		return
	}

	r, err := rtr.store.CreateRole(c, newRole)
	if err != nil {
		rtr.ErrorChooser(c, err)
	}

	c.JSON(http.StatusOK, r)
}

func (rtr *Router) DeleteRole(c *gin.Context) {
	var err error

	// ------------- Path parameter "id" -------------
	var id apiv1.EntityID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	if err := rtr.store.DeleteRole(c, id); err != nil {
		rtr.ErrorChooser(c, err)
		return
	}
}

func (rtr *Router) GetRole(c *gin.Context) {
	var err error

	// ------------- Path parameter "id" -------------
	var id apiv1.EntityID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	r, err := rtr.store.GetRole(c, id)
	if err != nil {
		rtr.ErrorChooser(c, err)
		return
	}

	c.JSON(http.StatusOK, r)
}

func (rtr *Router) UpdateRole(c *gin.Context) {
	var err error

	// ------------- Path parameter "id" -------------
	var id apiv1.EntityID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	role := &apiv1.Role{}

	if err := c.BindJSON(role); err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for role: %w", err), http.StatusBadRequest)
		return
	}

	// The ID in the path takes precedence. It should normally match
	// the ID in the body, but we don't enforce that.
	role.Id = id

	out, err := rtr.store.UpdateRole(c, role)
	if err != nil {
		rtr.ErrorChooser(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (rtr *Router) RemoveRoleAssignment(c *gin.Context) {
	var err error

	// ------------- Path parameter "id" -------------
	var id apiv1.EntityID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	assignment := apiv1.Assignment{}

	if err := c.BindJSON(&assignment); err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for assignment: %w", err), http.StatusBadRequest)
		return
	}

	assignment.Role = id

	if err := rtr.store.RemoveRoleAssignment(c, assignment); err != nil {
		rtr.ErrorChooser(c, err)
		return
	}
}

func (rtr *Router) GetRoleAssignments(c *gin.Context) {
	var err error

	// ------------- Path parameter "id" -------------
	var id apiv1.EntityID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	as, err := rtr.store.GetRoleAssignments(c, id)
	if err != nil {
		rtr.ErrorChooser(c, err)
		return
	}

	c.JSON(http.StatusOK, as)
}

func (rtr *Router) AssignRole(c *gin.Context) {
	var err error

	// ------------- Path parameter "id" -------------
	var id apiv1.EntityID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	ras := apiv1.NewRoleAssignment{}

	if err := c.BindJSON(&ras); err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for role assignment: %w", err), http.StatusBadRequest)
		return
	}

	if err := rtr.store.AssignRole(c, id, ras); err != nil {
		rtr.ErrorChooser(c, err)
		return
	}
}

func (rtr *Router) RemoveRolePermission(c *gin.Context) {
	var err error

	// ------------- Path parameter "id" -------------
	var id apiv1.EntityID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	pid := apiv1.PermissionIdentifier{}

	if err := c.BindJSON(&pid); err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for permission identifier: %w", err), http.StatusBadRequest)
		return
	}

	if err := rtr.store.RemoveRolePermission(c, id, pid); err != nil {
		rtr.ErrorChooser(c, err)
		return
	}
}

func (rtr *Router) GetRolePermissions(c *gin.Context) {
	var err error

	// ------------- Path parameter "id" -------------
	var id apiv1.EntityID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	perms, err := rtr.store.GetRolePermissions(c, id)
	if err != nil {
		rtr.ErrorChooser(c, err)
		return
	}

	c.JSON(http.StatusOK, perms)
}

func (rtr *Router) AddRolePermission(c *gin.Context) {
	var err error

	// ------------- Path parameter "id" -------------
	var id apiv1.EntityID

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	pid := apiv1.PermissionIdentifier{}

	if err = c.BindJSON(&pid); err != nil {
		rtr.ErrorHandler(c, fmt.Errorf("invalid format for permission identifier: %w", err), http.StatusBadRequest)
		return
	}

	if err := rtr.store.AddRolePermission(c, id, pid); err != nil {
		rtr.ErrorChooser(c, err)
		return
	}
}

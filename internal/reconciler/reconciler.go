package reconciler

import (
	"context"

	apiv1 "github.com/infratographer/fertilesoil/api/v1"
	appv1 "github.com/infratographer/fertilesoil/app/v1"
)

type Reconciler struct {
}

var _ appv1.Reconciler = &Reconciler{}

func NewReconciler() *Reconciler {
	return &Reconciler{}
}

func (r *Reconciler) Reconcile(ctx context.Context, evt apiv1.DirectoryEvent) error {
	// Here we will look at a directory event and perform the appropriate
	// action on the database. e.g. if the event is a create, we will
	// verify if we're tracking the parent directory, and if so, we will
	// ensure the inherited permissions are applied to the new directory.
	// deletions will be handled in a similar manner
	return nil
}

package controller

import (
	"github.com/johandry/nfs-operator/pkg/controller/nfs"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, nfs.Add)
}

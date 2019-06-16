package controller

import (
	"github.com/toVersus/aws-ssm-operator/pkg/controller/parameterstore"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, parameterstore.Add)
}

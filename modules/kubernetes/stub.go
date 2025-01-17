//go:build !k8s
// +build !k8s

package kubernetes

import (
	"github.com/AMuzykus/risor/object"
)

func Module() *object.Module {
	return nil
}

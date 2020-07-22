//go:generate genny -in=$GOFILE -out=generics_$GOFILE gen "T=int"

package pointers

import (
	"github.com/cheekybits/genny/generic"
)

type T generic.Type

func NewT(value T) *T {
	p := new(T)
	*p = value
	return p
}

//go:generate genny -in=$GOFILE -out=generics_$GOFILE gen "T=string"

package slices

import "github.com/cheekybits/genny/generic"

type T generic.Type

func T_Last(slice []T) T {
	return slice[len(slice)-1]
}

//go:build !solution

package genericsum

import (
	"math/cmplx"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func SortSlice[T constraints.Ordered](a []T) {
	slices.Sort(a)
}

func MapsEqual[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

func SliceContains[E comparable](s []E, v E) bool {
	for _, el := range s {
		if el == v {
			return true
		}
	}

	return false
}

func MergeChans[T any](chs ...<-chan T) <-chan T {
	res := make(chan T)
	go func() {
		open := len(chs)
		for {
			for _, ch := range chs {
				select {
				case v, ok := <-ch:
					if ok {
						res <- v
						continue
					}
					open--
					if open < 1 {
						close(res)
						return
					}
				default:
					continue
				}

			}
		}

	}()
	return res
}

type Numeric interface {
	constraints.Integer | constraints.Complex | constraints.Float
}

// save
func IsHermitianMatrix[T Numeric](m [][]T) bool {
	h := len(m)
	if h < 1 {
		return true
	}
	w := len(m[0])
	if w < 1 {
		return true
	}
	if h != w {
		return false
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			el := m[x][y]
			switch any(el).(type) {
			case complex64:
				opp, ok := any(m[y][x]).(complex64)
				if !ok {
					return false
				}
				current, ok := any(m[x][y]).(complex64)
				if !ok {
					return false
				}

				if real(opp) != real(current) {
					return false
				}
				if -imag(opp) != imag(current) {
					return false
				}

			case complex128:
				opp, ok := any(m[y][x]).(complex128)
				if !ok {
					return false
				}
				cur, ok := any(m[x][y]).(complex128)
				if !ok {
					return false
				}

				if cmplx.Conj(cur) != opp {
					return false
				}
			default:
				if m[x][y] != m[y][x] {
					return false
				}
			}
		}
	}

	return true
}

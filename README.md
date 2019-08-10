# Golang struct converter

Work in progress

```go
package main

import (
	"fmt"

	"github.com/ompluscator/convert-struct"
)

type (
	ia interface {
		M()
	}

	A struct {
		First   []int
		Second  []int
		Third   []int
		Fourth  interface{}
		Fifth   ia
		Sixth   int
		Seventh *int
		Eight   *int
	}

	B struct {
		First   []float64
		Second  []byte
		Third   []int
		Fourth  interface{}
		Fifth   ia
		Sixth   *int
		Seventh int
		Eight   *int
	}
)

func main() {
	value := 8
	converter := convertstruct.NewConverter(A{
		First:   []int{1, 2, 3},
		Second:  []int{1, 2, 3, 4},
		Third:   []int{1, 2, 3, 4, 5},
		Fourth:  1.0,
		Sixth:   7,
		Seventh: &value,
		Eight:   &value,
	})

	b := B{}

	fmt.Printf("%#v\n", b)
	fmt.Println(converter.Convert(&b))
	fmt.Printf(`%#v`, b)
	// Output
	// main.B{First:[]float64(nil), Second:[]uint8(nil), Third:[]int(nil), Fourth:interface {}(nil), Fifth:main.ia(nil), Sixth:(*int)(nil), Seventh:0, Eight:(*int)(nil)}
	// <nil>
	// main.B{First:[]float64{1, 2, 3}, Second:[]uint8{0x1, 0x2, 0x3, 0x4}, Third:[]int{1, 2, 3, 4, 5}, Fourth:1, Fifth:main.ia(nil), Sixth:(*int)(0xc000016180), Seventh:8, Eight:(*int)(0xc0000161c0)}\n
}
```
// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

type function struct {
	arity int
	fn    func(args []*big.Rat) *big.Rat
}

type functions map[string]function

// FunctionNames holds all the function names that are available for use
var FunctionNames []string

var funcs = make(functions)

func (f functions) register(name string, function function) {
	FunctionNames = append(FunctionNames, name)
	f[name] = function
}

func init() {
	funcs.register("abs", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			return new(big.Rat).Abs(args[0])
		},
	})
	funcs.register("ceil", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			return Ceil(args[0])
		},
	})
	funcs.register("floor", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			return Floor(args[0])
		},
	})
	funcs.register("sin", function{
		arity: -1,
		fn: func(args []*big.Rat) *big.Rat {
			if len(args) < 1 || len(args) > 2 {
				return new(big.Rat).SetFloat64(math.NaN())
			}
			float, _ := args[0].Float64()
			if len(args) == 2 {
				return new(big.Rat).SetFloat64(math.Sin(float * (math.Pi / 180)))
			}
			return new(big.Rat).SetFloat64(math.Sin(float))
		},
	})
	funcs.register("cos", function{
		arity: -1,
		fn: func(args []*big.Rat) *big.Rat {
			if len(args) < 1 || len(args) > 2 {
				return new(big.Rat).SetFloat64(math.NaN())
			}
			float, _ := args[0].Float64()
			if len(args) == 2 {
				return new(big.Rat).SetFloat64(math.Cos(float * (math.Pi / 180)))
			}
			return new(big.Rat).SetFloat64(math.Cos(float))
		},
	})
	funcs.register("tan", function{
		arity: -1,
		fn: func(args []*big.Rat) *big.Rat {
			if len(args) < 1 || len(args) > 2 {
				return new(big.Rat).SetFloat64(math.NaN())
			}
			float, _ := args[0].Float64()
			if len(args) == 2 {
				return new(big.Rat).SetFloat64(math.Tan(float * (math.Pi / 180)))
			}
			return new(big.Rat).SetFloat64(math.Tan(float))
		},
	})
	funcs.register("asin", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			float, _ := args[0].Float64()
			return new(big.Rat).SetFloat64(math.Asin(float))
		},
	})
	funcs.register("acos", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			float, _ := args[0].Float64()
			return new(big.Rat).SetFloat64(math.Acos(float))
		},
	})
	funcs.register("atan", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			float, _ := args[0].Float64()
			return new(big.Rat).SetFloat64(math.Atan(float))
		},
	})
	funcs.register("ln", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			float, _ := args[0].Float64()
			return new(big.Rat).SetFloat64(math.Log(float))
		},
	})
	funcs.register("log", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			float, _ := args[0].Float64()
			return new(big.Rat).SetFloat64(math.Log10(float))
		},
	})
	funcs.register("logn", function{
		arity: 2,
		fn: func(args []*big.Rat) *big.Rat {
			base, _ := args[0].Float64()
			arg, _ := args[1].Float64()
			return new(big.Rat).SetFloat64(math.Log10(arg) / math.Log10(base))
		},
	})
	funcs.register("max", function{
		arity: 2,
		fn: func(args []*big.Rat) *big.Rat {
			return Max(args[0], args[1])
		},
	})
	funcs.register("min", function{
		arity: 2,
		fn: func(args []*big.Rat) *big.Rat {
			return Min(args[0], args[1])
		},
	})
	funcs.register("sqrt", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			float, _ := args[0].Float64()
			return new(big.Rat).SetFloat64(math.Sqrt(float))
		},
	})
	funcs.register("rand", function{
		arity: 0,
		fn: func(_ []*big.Rat) *big.Rat {
			return new(big.Rat).SetFloat64(rand.Float64())
		},
	})
	funcs.register("fact", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			return Factorial(args[0])
		},
	})
	funcs.register("gcd", function{
		arity: 2,
		fn: func(args []*big.Rat) *big.Rat {
			return Gcd(args[0], args[1])
		},
	})
	funcs.register("list", function{
		arity: 0,
		fn: func(_ []*big.Rat) *big.Rat {
			for _, name := range FunctionNames {
				fmt.Print(name + " ")
			}
			fmt.Println()
			return RatTrue
		},
	})

	// custom functions
	funcs.register("csc", function{
		arity: -1,
		fn: func(args []*big.Rat) *big.Rat {
			if len(args) < 1 || len(args) > 2 {
				return new(big.Rat).SetFloat64(math.NaN())
			}
			float, _ := args[0].Float64()
			var sinVal float64
			if len(args) == 2 {
				sinVal = math.Sin(float * (math.Pi / 180))
			} else {
				sinVal = math.Sin(float)
			}
			if math.Abs(sinVal) < 1e-10 {
				return new(big.Rat).SetFloat64(math.NaN())
			}
			return new(big.Rat).SetFloat64(1 / sinVal)
		},
	})

	funcs.register("sec", function{
		arity: -1,
		fn: func(args []*big.Rat) *big.Rat {
			if len(args) < 1 || len(args) > 2 {
				return new(big.Rat).SetFloat64(math.NaN())
			}
			float, _ := args[0].Float64()
			if len(args) == 2 {
				return new(big.Rat).SetFloat64(1 / math.Cos(float*(math.Pi/180)))
			}
			return new(big.Rat).SetFloat64(1 / math.Cos(float))
		},
	})

	funcs.register("cot", function{
		arity: -1,
		fn: func(args []*big.Rat) *big.Rat {
			if len(args) < 1 || len(args) > 2 {
				return new(big.Rat).SetFloat64(math.NaN())
			}
			float, _ := args[0].Float64()
			if len(args) == 2 {
				return new(big.Rat).SetFloat64(1 / math.Tan(float*(math.Pi/180)))
			}
			return new(big.Rat).SetFloat64(1 / math.Tan(float))
		},
	})

	funcs.register("deg2rad", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			float, _ := args[0].Float64()
			return new(big.Rat).SetFloat64(float * (math.Pi / 180))
		},
	})

	funcs.register("rad2deg", function{
		arity: 1,
		fn: func(args []*big.Rat) *big.Rat {
			float, _ := args[0].Float64()
			return new(big.Rat).SetFloat64(float * 180 / math.Pi)
		},
	})
}

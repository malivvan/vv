package stdlib

import (
	"math"

	"github.com/malivvan/vv/vvm"
)

var mathModule = map[string]vvm.Object{
	"e":       &vvm.Float{Value: math.E},
	"pi":      &vvm.Float{Value: math.Pi},
	"phi":     &vvm.Float{Value: math.Phi},
	"sqrt2":   &vvm.Float{Value: math.Sqrt2},
	"sqrtE":   &vvm.Float{Value: math.SqrtE},
	"sqrtPi":  &vvm.Float{Value: math.SqrtPi},
	"sqrtPhi": &vvm.Float{Value: math.SqrtPhi},
	"ln2":     &vvm.Float{Value: math.Ln2},
	"log2E":   &vvm.Float{Value: math.Log2E},
	"ln10":    &vvm.Float{Value: math.Ln10},
	"log10E":  &vvm.Float{Value: math.Log10E},
	"abs": &vvm.BuiltinFunction{
		Name:  "abs",
		Value: FuncAFRF(math.Abs),
	},
	"acos": &vvm.BuiltinFunction{
		Name:  "acos",
		Value: FuncAFRF(math.Acos),
	},
	"acosh": &vvm.BuiltinFunction{
		Name:  "acosh",
		Value: FuncAFRF(math.Acosh),
	},
	"asin": &vvm.BuiltinFunction{
		Name:  "asin",
		Value: FuncAFRF(math.Asin),
	},
	"asinh": &vvm.BuiltinFunction{
		Name:  "asinh",
		Value: FuncAFRF(math.Asinh),
	},
	"atan": &vvm.BuiltinFunction{
		Name:  "atan",
		Value: FuncAFRF(math.Atan),
	},
	"atan2": &vvm.BuiltinFunction{
		Name:  "atan2",
		Value: FuncAFFRF(math.Atan2),
	},
	"atanh": &vvm.BuiltinFunction{
		Name:  "atanh",
		Value: FuncAFRF(math.Atanh),
	},
	"cbrt": &vvm.BuiltinFunction{
		Name:  "cbrt",
		Value: FuncAFRF(math.Cbrt),
	},
	"ceil": &vvm.BuiltinFunction{
		Name:  "ceil",
		Value: FuncAFRF(math.Ceil),
	},
	"copysign": &vvm.BuiltinFunction{
		Name:  "copysign",
		Value: FuncAFFRF(math.Copysign),
	},
	"cos": &vvm.BuiltinFunction{
		Name:  "cos",
		Value: FuncAFRF(math.Cos),
	},
	"cosh": &vvm.BuiltinFunction{
		Name:  "cosh",
		Value: FuncAFRF(math.Cosh),
	},
	"dim": &vvm.BuiltinFunction{
		Name:  "dim",
		Value: FuncAFFRF(math.Dim),
	},
	"erf": &vvm.BuiltinFunction{
		Name:  "erf",
		Value: FuncAFRF(math.Erf),
	},
	"erfc": &vvm.BuiltinFunction{
		Name:  "erfc",
		Value: FuncAFRF(math.Erfc),
	},
	"exp": &vvm.BuiltinFunction{
		Name:  "exp",
		Value: FuncAFRF(math.Exp),
	},
	"exp2": &vvm.BuiltinFunction{
		Name:  "exp2",
		Value: FuncAFRF(math.Exp2),
	},
	"expm1": &vvm.BuiltinFunction{
		Name:  "expm1",
		Value: FuncAFRF(math.Expm1),
	},
	"floor": &vvm.BuiltinFunction{
		Name:  "floor",
		Value: FuncAFRF(math.Floor),
	},
	"gamma": &vvm.BuiltinFunction{
		Name:  "gamma",
		Value: FuncAFRF(math.Gamma),
	},
	"hypot": &vvm.BuiltinFunction{
		Name:  "hypot",
		Value: FuncAFFRF(math.Hypot),
	},
	"ilogb": &vvm.BuiltinFunction{
		Name:  "ilogb",
		Value: FuncAFRI(math.Ilogb),
	},
	"inf": &vvm.BuiltinFunction{
		Name:  "inf",
		Value: FuncAIRF(math.Inf),
	},
	"is_inf": &vvm.BuiltinFunction{
		Name:  "is_inf",
		Value: FuncAFIRB(math.IsInf),
	},
	"is_nan": &vvm.BuiltinFunction{
		Name:  "is_nan",
		Value: FuncAFRB(math.IsNaN),
	},
	"j0": &vvm.BuiltinFunction{
		Name:  "j0",
		Value: FuncAFRF(math.J0),
	},
	"j1": &vvm.BuiltinFunction{
		Name:  "j1",
		Value: FuncAFRF(math.J1),
	},
	"jn": &vvm.BuiltinFunction{
		Name:  "jn",
		Value: FuncAIFRF(math.Jn),
	},
	"ldexp": &vvm.BuiltinFunction{
		Name:  "ldexp",
		Value: FuncAFIRF(math.Ldexp),
	},
	"log": &vvm.BuiltinFunction{
		Name:  "log",
		Value: FuncAFRF(math.Log),
	},
	"log10": &vvm.BuiltinFunction{
		Name:  "log10",
		Value: FuncAFRF(math.Log10),
	},
	"log1p": &vvm.BuiltinFunction{
		Name:  "log1p",
		Value: FuncAFRF(math.Log1p),
	},
	"log2": &vvm.BuiltinFunction{
		Name:  "log2",
		Value: FuncAFRF(math.Log2),
	},
	"logb": &vvm.BuiltinFunction{
		Name:  "logb",
		Value: FuncAFRF(math.Logb),
	},
	"max": &vvm.BuiltinFunction{
		Name:  "max",
		Value: FuncAFFRF(math.Max),
	},
	"min": &vvm.BuiltinFunction{
		Name:  "min",
		Value: FuncAFFRF(math.Min),
	},
	"mod": &vvm.BuiltinFunction{
		Name:  "mod",
		Value: FuncAFFRF(math.Mod),
	},
	"nan": &vvm.BuiltinFunction{
		Name:  "nan",
		Value: FuncARF(math.NaN),
	},
	"nextafter": &vvm.BuiltinFunction{
		Name:  "nextafter",
		Value: FuncAFFRF(math.Nextafter),
	},
	"pow": &vvm.BuiltinFunction{
		Name:  "pow",
		Value: FuncAFFRF(math.Pow),
	},
	"pow10": &vvm.BuiltinFunction{
		Name:  "pow10",
		Value: FuncAIRF(math.Pow10),
	},
	"remainder": &vvm.BuiltinFunction{
		Name:  "remainder",
		Value: FuncAFFRF(math.Remainder),
	},
	"signbit": &vvm.BuiltinFunction{
		Name:  "signbit",
		Value: FuncAFRB(math.Signbit),
	},
	"sin": &vvm.BuiltinFunction{
		Name:  "sin",
		Value: FuncAFRF(math.Sin),
	},
	"sinh": &vvm.BuiltinFunction{
		Name:  "sinh",
		Value: FuncAFRF(math.Sinh),
	},
	"sqrt": &vvm.BuiltinFunction{
		Name:  "sqrt",
		Value: FuncAFRF(math.Sqrt),
	},
	"tan": &vvm.BuiltinFunction{
		Name:  "tan",
		Value: FuncAFRF(math.Tan),
	},
	"tanh": &vvm.BuiltinFunction{
		Name:  "tanh",
		Value: FuncAFRF(math.Tanh),
	},
	"trunc": &vvm.BuiltinFunction{
		Name:  "trunc",
		Value: FuncAFRF(math.Trunc),
	},
	"y0": &vvm.BuiltinFunction{
		Name:  "y0",
		Value: FuncAFRF(math.Y0),
	},
	"y1": &vvm.BuiltinFunction{
		Name:  "y1",
		Value: FuncAFRF(math.Y1),
	},
	"yn": &vvm.BuiltinFunction{
		Name:  "yn",
		Value: FuncAIFRF(math.Yn),
	},
}

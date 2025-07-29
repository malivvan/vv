package vvm_test

import (
	"testing"

	"github.com/malivvan/vv/vvm"
	"github.com/malivvan/vv/vvm/require"
	"github.com/malivvan/vv/vvm/token"
)

func TestObject_TypeName(t *testing.T) {
	var o vvm.Object = &vvm.Int{}
	require.Equal(t, "int", o.TypeName())
	o = &vvm.Float{}
	require.Equal(t, "float", o.TypeName())
	o = &vvm.Char{}
	require.Equal(t, "char", o.TypeName())
	o = &vvm.String{}
	require.Equal(t, "string", o.TypeName())
	o = &vvm.Bool{}
	require.Equal(t, "bool", o.TypeName())
	o = &vvm.Array{}
	require.Equal(t, "array", o.TypeName())
	o = &vvm.Map{}
	require.Equal(t, "map", o.TypeName())
	o = &vvm.ArrayIterator{}
	require.Equal(t, "array-iterator", o.TypeName())
	o = &vvm.StringIterator{}
	require.Equal(t, "string-iterator", o.TypeName())
	o = &vvm.MapIterator{}
	require.Equal(t, "map-iterator", o.TypeName())
	o = &vvm.BuiltinFunction{Name: "fn"}
	require.Equal(t, "builtin-function:fn", o.TypeName())
	o = &vvm.UserFunction{Name: "fn"}
	require.Equal(t, "user-function:fn", o.TypeName())
	o = &vvm.CompiledFunction{}
	require.Equal(t, "compiled-function", o.TypeName())
	o = &vvm.Undefined{}
	require.Equal(t, "undefined", o.TypeName())
	o = &vvm.Error{}
	require.Equal(t, "error", o.TypeName())
	o = &vvm.Bytes{}
	require.Equal(t, "bytes", o.TypeName())
}

func TestObject_IsFalsy(t *testing.T) {
	var o vvm.Object = &vvm.Int{Value: 0}
	require.True(t, o.IsFalsy())
	o = &vvm.Int{Value: 1}
	require.False(t, o.IsFalsy())
	o = &vvm.Float{Value: 0}
	require.False(t, o.IsFalsy())
	o = &vvm.Float{Value: 1}
	require.False(t, o.IsFalsy())
	o = &vvm.Char{Value: ' '}
	require.False(t, o.IsFalsy())
	o = &vvm.Char{Value: 'T'}
	require.False(t, o.IsFalsy())
	o = &vvm.String{Value: ""}
	require.True(t, o.IsFalsy())
	o = &vvm.String{Value: " "}
	require.False(t, o.IsFalsy())
	o = &vvm.Array{Value: nil}
	require.True(t, o.IsFalsy())
	o = &vvm.Array{Value: []vvm.Object{nil}} // nil is not valid but still count as 1 element
	require.False(t, o.IsFalsy())
	o = &vvm.Map{Value: nil}
	require.True(t, o.IsFalsy())
	o = &vvm.Map{Value: map[string]vvm.Object{"a": nil}} // nil is not valid but still count as 1 element
	require.False(t, o.IsFalsy())
	o = &vvm.StringIterator{}
	require.True(t, o.IsFalsy())
	o = &vvm.ArrayIterator{}
	require.True(t, o.IsFalsy())
	o = &vvm.MapIterator{}
	require.True(t, o.IsFalsy())
	o = &vvm.BuiltinFunction{}
	require.False(t, o.IsFalsy())
	o = &vvm.CompiledFunction{}
	require.False(t, o.IsFalsy())
	o = &vvm.Undefined{}
	require.True(t, o.IsFalsy())
	o = &vvm.Error{}
	require.True(t, o.IsFalsy())
	o = &vvm.Bytes{}
	require.True(t, o.IsFalsy())
	o = &vvm.Bytes{Value: []byte{1, 2}}
	require.False(t, o.IsFalsy())
}

func TestObject_String(t *testing.T) {
	var o vvm.Object = &vvm.Int{Value: 0}
	require.Equal(t, "0", o.String())
	o = &vvm.Int{Value: 1}
	require.Equal(t, "1", o.String())
	o = &vvm.Float{Value: 0}
	require.Equal(t, "0", o.String())
	o = &vvm.Float{Value: 1}
	require.Equal(t, "1", o.String())
	o = &vvm.Char{Value: ' '}
	require.Equal(t, " ", o.String())
	o = &vvm.Char{Value: 'T'}
	require.Equal(t, "T", o.String())
	o = &vvm.String{Value: ""}
	require.Equal(t, `""`, o.String())
	o = &vvm.String{Value: " "}
	require.Equal(t, `" "`, o.String())
	o = &vvm.Array{Value: nil}
	require.Equal(t, "[]", o.String())
	o = &vvm.Map{Value: nil}
	require.Equal(t, "{}", o.String())
	o = &vvm.Error{Value: nil}
	require.Equal(t, "error", o.String())
	o = &vvm.Error{Value: &vvm.String{Value: "error 1"}}
	require.Equal(t, `error: "error 1"`, o.String())
	o = &vvm.StringIterator{}
	require.Equal(t, "<string-iterator>", o.String())
	o = &vvm.ArrayIterator{}
	require.Equal(t, "<array-iterator>", o.String())
	o = &vvm.MapIterator{}
	require.Equal(t, "<map-iterator>", o.String())
	o = &vvm.Undefined{}
	require.Equal(t, "<undefined>", o.String())
	o = &vvm.Bytes{}
	require.Equal(t, "", o.String())
	o = &vvm.Bytes{Value: []byte("foo")}
	require.Equal(t, "foo", o.String())
}

func TestObject_BinaryOp(t *testing.T) {
	var o vvm.Object = &vvm.Char{}
	_, err := o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
	o = &vvm.Bool{}
	_, err = o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
	o = &vvm.Map{}
	_, err = o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
	o = &vvm.ArrayIterator{}
	_, err = o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
	o = &vvm.StringIterator{}
	_, err = o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
	o = &vvm.MapIterator{}
	_, err = o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
	o = &vvm.BuiltinFunction{}
	_, err = o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
	o = &vvm.CompiledFunction{}
	_, err = o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
	o = &vvm.Undefined{}
	_, err = o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
	o = &vvm.Error{}
	_, err = o.BinaryOp(token.Add, vvm.UndefinedValue)
	require.Error(t, err)
}

func TestArray_BinaryOp(t *testing.T) {
	testBinaryOp(t, &vvm.Array{Value: nil}, token.Add,
		&vvm.Array{Value: nil}, &vvm.Array{Value: nil})
	testBinaryOp(t, &vvm.Array{Value: nil}, token.Add,
		&vvm.Array{Value: []vvm.Object{}}, &vvm.Array{Value: nil})
	testBinaryOp(t, &vvm.Array{Value: []vvm.Object{}}, token.Add,
		&vvm.Array{Value: nil}, &vvm.Array{Value: []vvm.Object{}})
	testBinaryOp(t, &vvm.Array{Value: []vvm.Object{}}, token.Add,
		&vvm.Array{Value: []vvm.Object{}},
		&vvm.Array{Value: []vvm.Object{}})
	testBinaryOp(t, &vvm.Array{Value: nil}, token.Add,
		&vvm.Array{Value: []vvm.Object{
			&vvm.Int{Value: 1},
		}}, &vvm.Array{Value: []vvm.Object{
			&vvm.Int{Value: 1},
		}})
	testBinaryOp(t, &vvm.Array{Value: nil}, token.Add,
		&vvm.Array{Value: []vvm.Object{
			&vvm.Int{Value: 1},
			&vvm.Int{Value: 2},
			&vvm.Int{Value: 3},
		}}, &vvm.Array{Value: []vvm.Object{
			&vvm.Int{Value: 1},
			&vvm.Int{Value: 2},
			&vvm.Int{Value: 3},
		}})
	testBinaryOp(t, &vvm.Array{Value: []vvm.Object{
		&vvm.Int{Value: 1},
		&vvm.Int{Value: 2},
		&vvm.Int{Value: 3},
	}}, token.Add, &vvm.Array{Value: nil},
		&vvm.Array{Value: []vvm.Object{
			&vvm.Int{Value: 1},
			&vvm.Int{Value: 2},
			&vvm.Int{Value: 3},
		}})
	testBinaryOp(t, &vvm.Array{Value: []vvm.Object{
		&vvm.Int{Value: 1},
		&vvm.Int{Value: 2},
		&vvm.Int{Value: 3},
	}}, token.Add, &vvm.Array{Value: []vvm.Object{
		&vvm.Int{Value: 4},
		&vvm.Int{Value: 5},
		&vvm.Int{Value: 6},
	}}, &vvm.Array{Value: []vvm.Object{
		&vvm.Int{Value: 1},
		&vvm.Int{Value: 2},
		&vvm.Int{Value: 3},
		&vvm.Int{Value: 4},
		&vvm.Int{Value: 5},
		&vvm.Int{Value: 6},
	}})
}

func TestError_Equals(t *testing.T) {
	err1 := &vvm.Error{Value: &vvm.String{Value: "some error"}}
	err2 := err1
	require.True(t, err1.Equals(err2))
	require.True(t, err2.Equals(err1))

	err2 = &vvm.Error{Value: &vvm.String{Value: "some error"}}
	require.False(t, err1.Equals(err2))
	require.False(t, err2.Equals(err1))
}

func TestFloat_BinaryOp(t *testing.T) {
	// float + float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Add,
				&vvm.Float{Value: r}, &vvm.Float{Value: l + r})
		}
	}

	// float - float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Sub,
				&vvm.Float{Value: r}, &vvm.Float{Value: l - r})
		}
	}

	// float * float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Mul,
				&vvm.Float{Value: r}, &vvm.Float{Value: l * r})
		}
	}

	// float / float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			if r != 0 {
				testBinaryOp(t, &vvm.Float{Value: l}, token.Quo,
					&vvm.Float{Value: r}, &vvm.Float{Value: l / r})
			}
		}
	}

	// float < float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Less,
				&vvm.Float{Value: r}, boolValue(l < r))
		}
	}

	// float > float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Greater,
				&vvm.Float{Value: r}, boolValue(l > r))
		}
	}

	// float <= float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &vvm.Float{Value: l}, token.LessEq,
				&vvm.Float{Value: r}, boolValue(l <= r))
		}
	}

	// float >= float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &vvm.Float{Value: l}, token.GreaterEq,
				&vvm.Float{Value: r}, boolValue(l >= r))
		}
	}

	// float + int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Add,
				&vvm.Int{Value: r}, &vvm.Float{Value: l + float64(r)})
		}
	}

	// float - int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Sub,
				&vvm.Int{Value: r}, &vvm.Float{Value: l - float64(r)})
		}
	}

	// float * int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Mul,
				&vvm.Int{Value: r}, &vvm.Float{Value: l * float64(r)})
		}
	}

	// float / int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			if r != 0 {
				testBinaryOp(t, &vvm.Float{Value: l}, token.Quo,
					&vvm.Int{Value: r},
					&vvm.Float{Value: l / float64(r)})
			}
		}
	}

	// float < int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Less,
				&vvm.Int{Value: r}, boolValue(l < float64(r)))
		}
	}

	// float > int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Float{Value: l}, token.Greater,
				&vvm.Int{Value: r}, boolValue(l > float64(r)))
		}
	}

	// float <= int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Float{Value: l}, token.LessEq,
				&vvm.Int{Value: r}, boolValue(l <= float64(r)))
		}
	}

	// float >= int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Float{Value: l}, token.GreaterEq,
				&vvm.Int{Value: r}, boolValue(l >= float64(r)))
		}
	}
}

func TestInt_BinaryOp(t *testing.T) {
	// int + int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Add,
				&vvm.Int{Value: r}, &vvm.Int{Value: l + r})
		}
	}

	// int - int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Sub,
				&vvm.Int{Value: r}, &vvm.Int{Value: l - r})
		}
	}

	// int * int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Mul,
				&vvm.Int{Value: r}, &vvm.Int{Value: l * r})
		}
	}

	// int / int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			if r != 0 {
				testBinaryOp(t, &vvm.Int{Value: l}, token.Quo,
					&vvm.Int{Value: r}, &vvm.Int{Value: l / r})
			}
		}
	}

	// int % int
	for l := int64(-4); l <= 4; l++ {
		for r := -int64(-4); r <= 4; r++ {
			if r == 0 {
				testBinaryOp(t, &vvm.Int{Value: l}, token.Rem,
					&vvm.Int{Value: r}, &vvm.Int{Value: l % r})
			}
		}
	}

	// int & int
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.And, &vvm.Int{Value: 0},
		&vvm.Int{Value: int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.And, &vvm.Int{Value: 0},
		&vvm.Int{Value: int64(1) & int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.And, &vvm.Int{Value: 1},
		&vvm.Int{Value: int64(0) & int64(1)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.And, &vvm.Int{Value: 1},
		&vvm.Int{Value: int64(1)})
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.And, &vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(0) & int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.And, &vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(1) & int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: int64(0xffffffff)}, token.And,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: 1984}, token.And,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(1984) & int64(0xffffffff)})
	testBinaryOp(t, &vvm.Int{Value: -1984}, token.And,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(-1984) & int64(0xffffffff)})

	// int | int
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.Or, &vvm.Int{Value: 0},
		&vvm.Int{Value: int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.Or, &vvm.Int{Value: 0},
		&vvm.Int{Value: int64(1) | int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.Or, &vvm.Int{Value: 1},
		&vvm.Int{Value: int64(0) | int64(1)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.Or, &vvm.Int{Value: 1},
		&vvm.Int{Value: int64(1)})
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.Or, &vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(0) | int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.Or, &vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(1) | int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: int64(0xffffffff)}, token.Or,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: 1984}, token.Or,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(1984) | int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: -1984}, token.Or,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(-1984) | int64(0xffffffff)})

	// int ^ int
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.Xor, &vvm.Int{Value: 0},
		&vvm.Int{Value: int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.Xor, &vvm.Int{Value: 0},
		&vvm.Int{Value: int64(1) ^ int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.Xor, &vvm.Int{Value: 1},
		&vvm.Int{Value: int64(0) ^ int64(1)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.Xor, &vvm.Int{Value: 1},
		&vvm.Int{Value: int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.Xor, &vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(0) ^ int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.Xor, &vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(1) ^ int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: int64(0xffffffff)}, token.Xor,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 1984}, token.Xor,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(1984) ^ int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: -1984}, token.Xor,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(-1984) ^ int64(0xffffffff)})

	// int &^ int
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.AndNot, &vvm.Int{Value: 0},
		&vvm.Int{Value: int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.AndNot, &vvm.Int{Value: 0},
		&vvm.Int{Value: int64(1) &^ int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.AndNot,
		&vvm.Int{Value: 1}, &vvm.Int{Value: int64(0) &^ int64(1)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.AndNot, &vvm.Int{Value: 1},
		&vvm.Int{Value: int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 0}, token.AndNot,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(0) &^ int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: 1}, token.AndNot,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(1) &^ int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: int64(0xffffffff)}, token.AndNot,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(0)})
	testBinaryOp(t,
		&vvm.Int{Value: 1984}, token.AndNot,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(1984) &^ int64(0xffffffff)})
	testBinaryOp(t,
		&vvm.Int{Value: -1984}, token.AndNot,
		&vvm.Int{Value: int64(0xffffffff)},
		&vvm.Int{Value: int64(-1984) &^ int64(0xffffffff)})

	// int << int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t,
			&vvm.Int{Value: 0}, token.Shl, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(0) << uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: 1}, token.Shl, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(1) << uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: 2}, token.Shl, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(2) << uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: -1}, token.Shl, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(-1) << uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: -2}, token.Shl, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(-2) << uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: int64(0xffffffff)}, token.Shl,
			&vvm.Int{Value: s},
			&vvm.Int{Value: int64(0xffffffff) << uint(s)})
	}

	// int >> int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t,
			&vvm.Int{Value: 0}, token.Shr, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(0) >> uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: 1}, token.Shr, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(1) >> uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: 2}, token.Shr, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(2) >> uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: -1}, token.Shr, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(-1) >> uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: -2}, token.Shr, &vvm.Int{Value: s},
			&vvm.Int{Value: int64(-2) >> uint(s)})
		testBinaryOp(t,
			&vvm.Int{Value: int64(0xffffffff)}, token.Shr,
			&vvm.Int{Value: s},
			&vvm.Int{Value: int64(0xffffffff) >> uint(s)})
	}

	// int < int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Less,
				&vvm.Int{Value: r}, boolValue(l < r))
		}
	}

	// int > int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Greater,
				&vvm.Int{Value: r}, boolValue(l > r))
		}
	}

	// int <= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Int{Value: l}, token.LessEq,
				&vvm.Int{Value: r}, boolValue(l <= r))
		}
	}

	// int >= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &vvm.Int{Value: l}, token.GreaterEq,
				&vvm.Int{Value: r}, boolValue(l >= r))
		}
	}

	// int + float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Add,
				&vvm.Float{Value: r},
				&vvm.Float{Value: float64(l) + r})
		}
	}

	// int - float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Sub,
				&vvm.Float{Value: r},
				&vvm.Float{Value: float64(l) - r})
		}
	}

	// int * float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Mul,
				&vvm.Float{Value: r},
				&vvm.Float{Value: float64(l) * r})
		}
	}

	// int / float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			if r != 0 {
				testBinaryOp(t, &vvm.Int{Value: l}, token.Quo,
					&vvm.Float{Value: r},
					&vvm.Float{Value: float64(l) / r})
			}
		}
	}

	// int < float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Less,
				&vvm.Float{Value: r}, boolValue(float64(l) < r))
		}
	}

	// int > float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &vvm.Int{Value: l}, token.Greater,
				&vvm.Float{Value: r}, boolValue(float64(l) > r))
		}
	}

	// int <= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &vvm.Int{Value: l}, token.LessEq,
				&vvm.Float{Value: r}, boolValue(float64(l) <= r))
		}
	}

	// int >= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &vvm.Int{Value: l}, token.GreaterEq,
				&vvm.Float{Value: r}, boolValue(float64(l) >= r))
		}
	}
}

func TestMap_Index(t *testing.T) {
	m := &vvm.Map{Value: make(map[string]vvm.Object)}
	k := &vvm.Int{Value: 1}
	v := &vvm.String{Value: "abcdef"}
	err := m.IndexSet(k, v)

	require.NoError(t, err)

	res, err := m.IndexGet(k)
	require.NoError(t, err)
	require.Equal(t, v, res)
}

func TestString_BinaryOp(t *testing.T) {
	lstr := "abcde"
	rstr := "01234"
	for l := 0; l < len(lstr); l++ {
		for r := 0; r < len(rstr); r++ {
			ls := lstr[l:]
			rs := rstr[r:]
			testBinaryOp(t, &vvm.String{Value: ls}, token.Add,
				&vvm.String{Value: rs},
				&vvm.String{Value: ls + rs})

			rc := []rune(rstr)[r]
			testBinaryOp(t, &vvm.String{Value: ls}, token.Add,
				&vvm.Char{Value: rc},
				&vvm.String{Value: ls + string(rc)})
		}
	}
}

func testBinaryOp(
	t *testing.T,
	lhs vvm.Object,
	op token.Token,
	rhs vvm.Object,
	expected vvm.Object,
) {
	t.Helper()
	actual, err := lhs.BinaryOp(op, rhs)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func boolValue(b bool) vvm.Object {
	if b {
		return vvm.TrueValue
	}
	return vvm.FalseValue
}

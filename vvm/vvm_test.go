package vvm_test

import (
	"strings"
	"testing"
	"time"

	"github.com/malivvan/vv/vvm"
	"github.com/malivvan/vv/vvm/parser"
	"github.com/malivvan/vv/vvm/require"
)

func TestInstructions_String(t *testing.T) {
	assertInstructionString(t,
		[][]byte{
			vvm.MakeInstruction(parser.OpConstant, 1),
			vvm.MakeInstruction(parser.OpConstant, 2),
			vvm.MakeInstruction(parser.OpConstant, 65535),
		},
		`0000 CONST   1    
0003 CONST   2    
0006 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			vvm.MakeInstruction(parser.OpBinaryOp, 11),
			vvm.MakeInstruction(parser.OpConstant, 2),
			vvm.MakeInstruction(parser.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 CONST   2    
0005 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			vvm.MakeInstruction(parser.OpBinaryOp, 11),
			vvm.MakeInstruction(parser.OpGetLocal, 1),
			vvm.MakeInstruction(parser.OpConstant, 2),
			vvm.MakeInstruction(parser.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 GETL    1    
0004 CONST   2    
0007 CONST   65535`)
}

func TestMakeInstruction(t *testing.T) {
	makeInstruction(t, []byte{parser.OpConstant, 0, 0},
		parser.OpConstant, 0)
	makeInstruction(t, []byte{parser.OpConstant, 0, 1},
		parser.OpConstant, 1)
	makeInstruction(t, []byte{parser.OpConstant, 255, 254},
		parser.OpConstant, 65534)
	makeInstruction(t, []byte{parser.OpPop}, parser.OpPop)
	makeInstruction(t, []byte{parser.OpTrue}, parser.OpTrue)
	makeInstruction(t, []byte{parser.OpFalse}, parser.OpFalse)
}

func TestNumObjects(t *testing.T) {
	testCountObjects(t, &vvm.Array{}, 1)
	testCountObjects(t, &vvm.Array{Value: []vvm.Object{
		&vvm.Int{Value: 1},
		&vvm.Int{Value: 2},
		&vvm.Array{Value: []vvm.Object{
			&vvm.Int{Value: 3},
			&vvm.Int{Value: 4},
			&vvm.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, vvm.TrueValue, 1)
	testCountObjects(t, vvm.FalseValue, 1)
	testCountObjects(t, &vvm.BuiltinFunction{}, 1)
	testCountObjects(t, &vvm.Bytes{Value: []byte("foobar")}, 1)
	testCountObjects(t, &vvm.Char{Value: '가'}, 1)
	testCountObjects(t, &vvm.CompiledFunction{}, 1)
	testCountObjects(t, &vvm.Error{Value: &vvm.Int{Value: 5}}, 2)
	testCountObjects(t, &vvm.Float{Value: 19.84}, 1)
	testCountObjects(t, &vvm.ImmutableArray{Value: []vvm.Object{
		&vvm.Int{Value: 1},
		&vvm.Int{Value: 2},
		&vvm.ImmutableArray{Value: []vvm.Object{
			&vvm.Int{Value: 3},
			&vvm.Int{Value: 4},
			&vvm.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			"k1": &vvm.Int{Value: 1},
			"k2": &vvm.Int{Value: 2},
			"k3": &vvm.Array{Value: []vvm.Object{
				&vvm.Int{Value: 3},
				&vvm.Int{Value: 4},
				&vvm.Int{Value: 5},
			}},
		}}, 7)
	testCountObjects(t, &vvm.Int{Value: 1984}, 1)
	testCountObjects(t, &vvm.Map{Value: map[string]vvm.Object{
		"k1": &vvm.Int{Value: 1},
		"k2": &vvm.Int{Value: 2},
		"k3": &vvm.Array{Value: []vvm.Object{
			&vvm.Int{Value: 3},
			&vvm.Int{Value: 4},
			&vvm.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &vvm.String{Value: "foo bar"}, 1)
	testCountObjects(t, &vvm.Time{Value: time.Now()}, 1)
	testCountObjects(t, vvm.UndefinedValue, 1)
}

func testCountObjects(t *testing.T, o vvm.Object, expected int) {
	require.Equal(t, expected, vvm.CountObjects(o))
}

func assertInstructionString(
	t *testing.T,
	instructions [][]byte,
	expected string,
) {
	concatted := make([]byte, 0)
	for _, e := range instructions {
		concatted = append(concatted, e...)
	}
	require.Equal(t, expected, strings.Join(
		vvm.FormatInstructions(concatted, 0), "\n"))
}

func makeInstruction(
	t *testing.T,
	expected []byte,
	opcode parser.Opcode,
	operands ...int,
) {
	inst := vvm.MakeInstruction(opcode, operands...)
	require.Equal(t, expected, inst)
}

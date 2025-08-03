package vvm_test

import (
	"testing"
	"time"

	"github.com/malivvan/vv/vvm"
	"github.com/malivvan/vv/vvm/parser"
	"github.com/malivvan/vv/vvm/require"
)

type srcfile struct {
	name string
	size int
}

func TestBytecode(t *testing.T) {
	testBytecodeSerialization(t, bytecode(concatInsts(), objectsArray()))

	testBytecodeSerialization(t, bytecode(
		concatInsts(), objectsArray(
			&vvm.Char{Value: 'y'},
			&vvm.Float{Value: 93.11},
			compiledFunction(1, 0,
				vvm.MakeInstruction(parser.OpConstant, 3),
				vvm.MakeInstruction(parser.OpSetLocal, 0),
				vvm.MakeInstruction(parser.OpGetGlobal, 0),
				vvm.MakeInstruction(parser.OpGetFree, 0)),
			&vvm.Float{Value: 39.2},
			&vvm.Int{Value: 192},
			&vvm.String{Value: "bar"})))

	testBytecodeSerialization(t, bytecodeFileSet(
		concatInsts(
			vvm.MakeInstruction(parser.OpConstant, 0),
			vvm.MakeInstruction(parser.OpSetGlobal, 0),
			vvm.MakeInstruction(parser.OpConstant, 6),
			vvm.MakeInstruction(parser.OpPop)),
		objectsArray(
			&vvm.Int{Value: 55},
			&vvm.Int{Value: 66},
			&vvm.Int{Value: 77},
			&vvm.Int{Value: 88},
			&vvm.ImmutableMap{
				Value: map[string]vvm.Object{
					"array": &vvm.ImmutableArray{
						Value: []vvm.Object{
							&vvm.Int{Value: 1},
							&vvm.Int{Value: 2},
							&vvm.Int{Value: 3},
							vvm.TrueValue,
							vvm.FalseValue,
							vvm.UndefinedValue,
						},
					},
					"true":  vvm.TrueValue,
					"false": vvm.FalseValue,
					"bytes": &vvm.Bytes{Value: make([]byte, 16)},
					"char":  &vvm.Char{Value: 'Y'},
					"error": &vvm.Error{Value: &vvm.String{
						Value: "some error",
					}},
					"float": &vvm.Float{Value: -19.84},
					"immutable_array": &vvm.ImmutableArray{
						Value: []vvm.Object{
							&vvm.Int{Value: 1},
							&vvm.Int{Value: 2},
							&vvm.Int{Value: 3},
							vvm.TrueValue,
							vvm.FalseValue,
							vvm.UndefinedValue,
						},
					},
					"immutable_map": &vvm.ImmutableMap{
						Value: map[string]vvm.Object{
							"a": &vvm.Int{Value: 1},
							"b": &vvm.Int{Value: 2},
							"c": &vvm.Int{Value: 3},
							"d": vvm.TrueValue,
							"e": vvm.FalseValue,
							"f": vvm.UndefinedValue,
						},
					},
					"int": &vvm.Int{Value: 91},
					"map": &vvm.Map{
						Value: map[string]vvm.Object{
							"a": &vvm.Int{Value: 1},
							"b": &vvm.Int{Value: 2},
							"c": &vvm.Int{Value: 3},
							"d": vvm.TrueValue,
							"e": vvm.FalseValue,
							"f": vvm.UndefinedValue,
						},
					},
					"string":    &vvm.String{Value: "foo bar"},
					"time":      &vvm.Time{Value: time.Now()},
					"undefined": vvm.UndefinedValue,
				},
			},
			compiledFunction(1, 0,
				vvm.MakeInstruction(parser.OpConstant, 3),
				vvm.MakeInstruction(parser.OpSetLocal, 0),
				vvm.MakeInstruction(parser.OpGetGlobal, 0),
				vvm.MakeInstruction(parser.OpGetFree, 0),
				vvm.MakeInstruction(parser.OpBinaryOp, 11),
				vvm.MakeInstruction(parser.OpGetFree, 1),
				vvm.MakeInstruction(parser.OpBinaryOp, 11),
				vvm.MakeInstruction(parser.OpGetLocal, 0),
				vvm.MakeInstruction(parser.OpBinaryOp, 11),
				vvm.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				vvm.MakeInstruction(parser.OpConstant, 2),
				vvm.MakeInstruction(parser.OpSetLocal, 0),
				vvm.MakeInstruction(parser.OpGetFree, 0),
				vvm.MakeInstruction(parser.OpGetLocal, 0),
				vvm.MakeInstruction(parser.OpClosure, 4, 2),
				vvm.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				vvm.MakeInstruction(parser.OpConstant, 1),
				vvm.MakeInstruction(parser.OpSetLocal, 0),
				vvm.MakeInstruction(parser.OpGetLocal, 0),
				vvm.MakeInstruction(parser.OpClosure, 5, 1),
				vvm.MakeInstruction(parser.OpReturn, 1))),
		fileSet(srcfile{name: "file1", size: 100},
			srcfile{name: "file2", size: 200})))
}

func TestBytecode_RemoveDuplicates(t *testing.T) {
	testBytecodeRemoveDuplicates(t,
		bytecode(
			concatInsts(), objectsArray(
				&vvm.Char{Value: 'y'},
				&vvm.Float{Value: 93.11},
				compiledFunction(1, 0,
					vvm.MakeInstruction(parser.OpConstant, 3),
					vvm.MakeInstruction(parser.OpSetLocal, 0),
					vvm.MakeInstruction(parser.OpGetGlobal, 0),
					vvm.MakeInstruction(parser.OpGetFree, 0)),
				&vvm.Float{Value: 39.2},
				&vvm.Int{Value: 192},
				&vvm.String{Value: "bar"})),
		bytecode(
			concatInsts(), objectsArray(
				&vvm.Char{Value: 'y'},
				&vvm.Float{Value: 93.11},
				compiledFunction(1, 0,
					vvm.MakeInstruction(parser.OpConstant, 3),
					vvm.MakeInstruction(parser.OpSetLocal, 0),
					vvm.MakeInstruction(parser.OpGetGlobal, 0),
					vvm.MakeInstruction(parser.OpGetFree, 0)),
				&vvm.Float{Value: 39.2},
				&vvm.Int{Value: 192},
				&vvm.String{Value: "bar"})))

	testBytecodeRemoveDuplicates(t,
		bytecode(
			concatInsts(
				vvm.MakeInstruction(parser.OpConstant, 0),
				vvm.MakeInstruction(parser.OpConstant, 1),
				vvm.MakeInstruction(parser.OpConstant, 2),
				vvm.MakeInstruction(parser.OpConstant, 3),
				vvm.MakeInstruction(parser.OpConstant, 4),
				vvm.MakeInstruction(parser.OpConstant, 5),
				vvm.MakeInstruction(parser.OpConstant, 6),
				vvm.MakeInstruction(parser.OpConstant, 7),
				vvm.MakeInstruction(parser.OpConstant, 8),
				vvm.MakeInstruction(parser.OpClosure, 4, 1)),
			objectsArray(
				&vvm.Int{Value: 1},
				&vvm.Float{Value: 2.0},
				&vvm.Char{Value: '3'},
				&vvm.String{Value: "four"},
				compiledFunction(1, 0,
					vvm.MakeInstruction(parser.OpConstant, 3),
					vvm.MakeInstruction(parser.OpConstant, 7),
					vvm.MakeInstruction(parser.OpSetLocal, 0),
					vvm.MakeInstruction(parser.OpGetGlobal, 0),
					vvm.MakeInstruction(parser.OpGetFree, 0)),
				&vvm.Int{Value: 1},
				&vvm.Float{Value: 2.0},
				&vvm.Char{Value: '3'},
				&vvm.String{Value: "four"})),
		bytecode(
			concatInsts(
				vvm.MakeInstruction(parser.OpConstant, 0),
				vvm.MakeInstruction(parser.OpConstant, 1),
				vvm.MakeInstruction(parser.OpConstant, 2),
				vvm.MakeInstruction(parser.OpConstant, 3),
				vvm.MakeInstruction(parser.OpConstant, 4),
				vvm.MakeInstruction(parser.OpConstant, 0),
				vvm.MakeInstruction(parser.OpConstant, 1),
				vvm.MakeInstruction(parser.OpConstant, 2),
				vvm.MakeInstruction(parser.OpConstant, 3),
				vvm.MakeInstruction(parser.OpClosure, 4, 1)),
			objectsArray(
				&vvm.Int{Value: 1},
				&vvm.Float{Value: 2.0},
				&vvm.Char{Value: '3'},
				&vvm.String{Value: "four"},
				compiledFunction(1, 0,
					vvm.MakeInstruction(parser.OpConstant, 3),
					vvm.MakeInstruction(parser.OpConstant, 2),
					vvm.MakeInstruction(parser.OpSetLocal, 0),
					vvm.MakeInstruction(parser.OpGetGlobal, 0),
					vvm.MakeInstruction(parser.OpGetFree, 0)))))

	testBytecodeRemoveDuplicates(t,
		bytecode(
			concatInsts(
				vvm.MakeInstruction(parser.OpConstant, 0),
				vvm.MakeInstruction(parser.OpConstant, 1),
				vvm.MakeInstruction(parser.OpConstant, 2),
				vvm.MakeInstruction(parser.OpConstant, 3),
				vvm.MakeInstruction(parser.OpConstant, 4)),
			objectsArray(
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2},
				&vvm.Int{Value: 3},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 3})),
		bytecode(
			concatInsts(
				vvm.MakeInstruction(parser.OpConstant, 0),
				vvm.MakeInstruction(parser.OpConstant, 1),
				vvm.MakeInstruction(parser.OpConstant, 2),
				vvm.MakeInstruction(parser.OpConstant, 0),
				vvm.MakeInstruction(parser.OpConstant, 2)),
			objectsArray(
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2},
				&vvm.Int{Value: 3})))
}

func TestBytecode_CountObjects(t *testing.T) {
	b := bytecode(
		concatInsts(),
		objectsArray(
			&vvm.Int{Value: 55},
			&vvm.Int{Value: 66},
			&vvm.Int{Value: 77},
			&vvm.Int{Value: 88},
			compiledFunction(1, 0,
				vvm.MakeInstruction(parser.OpConstant, 3),
				vvm.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				vvm.MakeInstruction(parser.OpConstant, 2),
				vvm.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				vvm.MakeInstruction(parser.OpConstant, 1),
				vvm.MakeInstruction(parser.OpReturn, 1))))
	require.Equal(t, 7, b.CountObjects())
}

func fileSet(files ...srcfile) *parser.SourceFileSet {
	fileSet := parser.NewFileSet()
	for _, f := range files {
		fileSet.AddFile(f.name, -1, f.size)
	}
	return fileSet
}

func bytecodeFileSet(
	instructions []byte,
	constants []vvm.Object,
	fileSet *parser.SourceFileSet,
) *vvm.Bytecode {
	return &vvm.Bytecode{
		FileSet:      fileSet,
		MainFunction: &vvm.CompiledFunction{Instructions: instructions},
		Constants:    constants,
	}
}

func testBytecodeRemoveDuplicates(
	t *testing.T,
	input, expected *vvm.Bytecode,
) {
	input.RemoveDuplicates()

	require.Equal(t, expected.FileSet, input.FileSet)
	require.Equal(t, expected.MainFunction, input.MainFunction)
	require.Equal(t, expected.Constants, input.Constants)
}

func testBytecodeSerialization(t *testing.T, b *vvm.Bytecode) {
	bc, err := b.Marshal()
	require.NoError(t, err)

	r := &vvm.Bytecode{}
	err = r.Unmarshal(bc, nil)
	require.NoError(t, err)

	require.Equal(t, b.FileSet, r.FileSet)
	require.Equal(t, b.MainFunction, r.MainFunction)
	require.Equal(t, b.Constants, r.Constants)
}

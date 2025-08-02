package vvm_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/malivvan/vv/vvm"
)

func Test_builtinDelete(t *testing.T) {
	var builtinDelete func(ctx context.Context, args ...vvm.Object) (vvm.Object, error)
	for _, f := range vvm.GetAllBuiltinFunctions() {
		if f.Name == "delete" {
			builtinDelete = f.Value
			break
		}
	}
	if builtinDelete == nil {
		t.Fatal("builtin delete not found")
	}
	type args struct {
		args []vvm.Object
	}
	tests := []struct {
		name      string
		args      args
		want      vvm.Object
		wantErr   bool
		wantedErr error
		target    interface{}
	}{
		{name: "invalid-arg", args: args{[]vvm.Object{&vvm.String{},
			&vvm.String{}}}, wantErr: true,
			wantedErr: vvm.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "map",
				Found:    "string"},
		},
		{name: "no-args",
			wantErr: true, wantedErr: vvm.ErrWrongNumArguments},
		{name: "empty-args", args: args{[]vvm.Object{}}, wantErr: true,
			wantedErr: vvm.ErrWrongNumArguments,
		},
		{name: "3-args", args: args{[]vvm.Object{
			(*vvm.Map)(nil), (*vvm.String)(nil), (*vvm.String)(nil)}},
			wantErr: true, wantedErr: vvm.ErrWrongNumArguments,
		},
		{name: "nil-map-empty-key",
			args: args{[]vvm.Object{&vvm.Map{}, &vvm.String{}}},
			want: vvm.UndefinedValue,
		},
		{name: "nil-map-nonstr-key",
			args: args{[]vvm.Object{
				&vvm.Map{}, &vvm.Int{}}}, wantErr: true,
			wantedErr: vvm.ErrInvalidArgumentType{
				Name: "second", Expected: "string", Found: "int"},
		},
		{name: "nil-map-no-key",
			args: args{[]vvm.Object{&vvm.Map{}}}, wantErr: true,
			wantedErr: vvm.ErrWrongNumArguments,
		},
		{name: "map-missing-key",
			args: args{
				[]vvm.Object{
					&vvm.Map{Value: map[string]vvm.Object{
						"key": &vvm.String{Value: "value"},
					}},
					&vvm.String{Value: "key1"}}},
			want: vvm.UndefinedValue,
			target: &vvm.Map{
				Value: map[string]vvm.Object{
					"key": &vvm.String{
						Value: "value"}}},
		},
		{name: "map-emptied",
			args: args{
				[]vvm.Object{
					&vvm.Map{Value: map[string]vvm.Object{
						"key": &vvm.String{Value: "value"},
					}},
					&vvm.String{Value: "key"}}},
			want:   vvm.UndefinedValue,
			target: &vvm.Map{Value: map[string]vvm.Object{}},
		},
		{name: "map-multi-keys",
			args: args{
				[]vvm.Object{
					&vvm.Map{Value: map[string]vvm.Object{
						"key1": &vvm.String{Value: "value1"},
						"key2": &vvm.Int{Value: 10},
					}},
					&vvm.String{Value: "key1"}}},
			want: vvm.UndefinedValue,
			target: &vvm.Map{Value: map[string]vvm.Object{
				"key2": &vvm.Int{Value: 10}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinDelete(context.Background(), tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinDelete() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantedErr) {
				if err.Error() != tt.wantedErr.Error() {
					t.Errorf("builtinDelete() error = %v, wantedErr %v",
						err, tt.wantedErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("builtinDelete() = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && tt.target != nil {
				switch v := tt.args.args[0].(type) {
				case *vvm.Map, *vvm.Array:
					if !reflect.DeepEqual(tt.target, tt.args.args[0]) {
						t.Errorf("builtinDelete() objects are not equal "+
							"got: %+v, want: %+v", tt.args.args[0], tt.target)
					}
				default:
					t.Errorf("builtinDelete() unsuporrted arg[0] type %s",
						v.TypeName())
					return
				}
			}
		})
	}
}

func Test_builtinSplice(t *testing.T) {
	var builtinSplice func(ctx context.Context, args ...vvm.Object) (vvm.Object, error)
	for _, f := range vvm.GetAllBuiltinFunctions() {
		if f.Name == "splice" {
			builtinSplice = f.Value
			break
		}
	}
	if builtinSplice == nil {
		t.Fatal("builtin splice not found")
	}
	tests := []struct {
		name      string
		args      []vvm.Object
		deleted   vvm.Object
		Array     *vvm.Array
		wantErr   bool
		wantedErr error
	}{
		{name: "no args", args: []vvm.Object{}, wantErr: true,
			wantedErr: vvm.ErrWrongNumArguments,
		},
		{name: "invalid args", args: []vvm.Object{&vvm.Map{}},
			wantErr: true,
			wantedErr: vvm.ErrInvalidArgumentType{
				Name: "first", Expected: "array", Found: "map"},
		},
		{name: "invalid args",
			args:    []vvm.Object{&vvm.Array{}, &vvm.String{}},
			wantErr: true,
			wantedErr: vvm.ErrInvalidArgumentType{
				Name: "second", Expected: "int", Found: "string"},
		},
		{name: "negative index",
			args:      []vvm.Object{&vvm.Array{}, &vvm.Int{Value: -1}},
			wantErr:   true,
			wantedErr: vvm.ErrIndexOutOfBounds},
		{name: "non int count",
			args: []vvm.Object{
				&vvm.Array{}, &vvm.Int{Value: 0},
				&vvm.String{Value: ""}},
			wantErr: true,
			wantedErr: vvm.ErrInvalidArgumentType{
				Name: "third", Expected: "int", Found: "string"},
		},
		{name: "negative count",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}},
				&vvm.Int{Value: 0},
				&vvm.Int{Value: -1}},
			wantErr:   true,
			wantedErr: vvm.ErrIndexOutOfBounds,
		},
		{name: "insert with zero count",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}},
				&vvm.Int{Value: 0},
				&vvm.Int{Value: 0},
				&vvm.String{Value: "b"}},
			deleted: &vvm.Array{Value: []vvm.Object{}},
			Array: &vvm.Array{Value: []vvm.Object{
				&vvm.String{Value: "b"},
				&vvm.Int{Value: 0},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2}}},
		},
		{name: "insert",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 0},
				&vvm.String{Value: "c"},
				&vvm.String{Value: "d"}},
			deleted: &vvm.Array{Value: []vvm.Object{}},
			Array: &vvm.Array{Value: []vvm.Object{
				&vvm.Int{Value: 0},
				&vvm.String{Value: "c"},
				&vvm.String{Value: "d"},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2}}},
		},
		{name: "insert with zero count",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 0},
				&vvm.String{Value: "c"},
				&vvm.String{Value: "d"}},
			deleted: &vvm.Array{Value: []vvm.Object{}},
			Array: &vvm.Array{Value: []vvm.Object{
				&vvm.Int{Value: 0},
				&vvm.String{Value: "c"},
				&vvm.String{Value: "d"},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2}}},
		},
		{name: "insert with delete",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 1},
				&vvm.String{Value: "c"},
				&vvm.String{Value: "d"}},
			deleted: &vvm.Array{
				Value: []vvm.Object{&vvm.Int{Value: 1}}},
			Array: &vvm.Array{Value: []vvm.Object{
				&vvm.Int{Value: 0},
				&vvm.String{Value: "c"},
				&vvm.String{Value: "d"},
				&vvm.Int{Value: 2}}},
		},
		{name: "insert with delete multi",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2},
				&vvm.String{Value: "c"},
				&vvm.String{Value: "d"}},
			deleted: &vvm.Array{Value: []vvm.Object{
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2}}},
			Array: &vvm.Array{
				Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.String{Value: "c"},
					&vvm.String{Value: "d"}}},
		},
		{name: "delete all with positive count",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}},
				&vvm.Int{Value: 0},
				&vvm.Int{Value: 3}},
			deleted: &vvm.Array{Value: []vvm.Object{
				&vvm.Int{Value: 0},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2}}},
			Array: &vvm.Array{Value: []vvm.Object{}},
		},
		{name: "delete all with big count",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}},
				&vvm.Int{Value: 0},
				&vvm.Int{Value: 5}},
			deleted: &vvm.Array{Value: []vvm.Object{
				&vvm.Int{Value: 0},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2}}},
			Array: &vvm.Array{Value: []vvm.Object{}},
		},
		{name: "nothing2",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}}},
			Array: &vvm.Array{Value: []vvm.Object{}},
			deleted: &vvm.Array{Value: []vvm.Object{
				&vvm.Int{Value: 0},
				&vvm.Int{Value: 1},
				&vvm.Int{Value: 2}}},
		},
		{name: "pop without count",
			args: []vvm.Object{
				&vvm.Array{Value: []vvm.Object{
					&vvm.Int{Value: 0},
					&vvm.Int{Value: 1},
					&vvm.Int{Value: 2}}},
				&vvm.Int{Value: 2}},
			deleted: &vvm.Array{Value: []vvm.Object{&vvm.Int{Value: 2}}},
			Array: &vvm.Array{Value: []vvm.Object{
				&vvm.Int{Value: 0}, &vvm.Int{Value: 1}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinSplice(context.Background(), tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinSplice() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.deleted) {
				t.Errorf("builtinSplice() = %v, want %v", got, tt.deleted)
			}
			if tt.wantErr && tt.wantedErr.Error() != err.Error() {
				t.Errorf("builtinSplice() error = %v, wantedErr %v",
					err, tt.wantedErr)
			}
			if tt.Array != nil && !reflect.DeepEqual(tt.Array, tt.args[0]) {
				t.Errorf("builtinSplice() arrays are not equal expected"+
					" %s, got %s", tt.Array, tt.args[0].(*vvm.Array))
			}
		})
	}
}

func Test_builtinRange(t *testing.T) {
	var builtinRange func(ctx context.Context, args ...vvm.Object) (vvm.Object, error)
	for _, f := range vvm.GetAllBuiltinFunctions() {
		if f.Name == "range" {
			builtinRange = f.Value
			break
		}
	}
	if builtinRange == nil {
		t.Fatal("builtin range not found")
	}
	tests := []struct {
		name      string
		args      []vvm.Object
		result    *vvm.Array
		wantErr   bool
		wantedErr error
	}{
		{name: "no args", args: []vvm.Object{}, wantErr: true,
			wantedErr: vvm.ErrWrongNumArguments,
		},
		{name: "single args", args: []vvm.Object{&vvm.Map{}},
			wantErr:   true,
			wantedErr: vvm.ErrWrongNumArguments,
		},
		{name: "4 args", args: []vvm.Object{&vvm.Map{}, &vvm.String{}, &vvm.String{}, &vvm.String{}},
			wantErr:   true,
			wantedErr: vvm.ErrWrongNumArguments,
		},
		{name: "invalid start",
			args:    []vvm.Object{&vvm.String{}, &vvm.String{}},
			wantErr: true,
			wantedErr: vvm.ErrInvalidArgumentType{
				Name: "start", Expected: "int", Found: "string"},
		},
		{name: "invalid stop",
			args:    []vvm.Object{&vvm.Int{}, &vvm.String{}},
			wantErr: true,
			wantedErr: vvm.ErrInvalidArgumentType{
				Name: "stop", Expected: "int", Found: "string"},
		},
		{name: "invalid step",
			args:    []vvm.Object{&vvm.Int{}, &vvm.Int{}, &vvm.String{}},
			wantErr: true,
			wantedErr: vvm.ErrInvalidArgumentType{
				Name: "step", Expected: "int", Found: "string"},
		},
		{name: "zero step",
			args:      []vvm.Object{&vvm.Int{}, &vvm.Int{}, &vvm.Int{}}, //must greate than 0
			wantErr:   true,
			wantedErr: vvm.ErrInvalidRangeStep,
		},
		{name: "negative step",
			args:      []vvm.Object{&vvm.Int{}, &vvm.Int{}, intObject(-2)}, //must greate than 0
			wantErr:   true,
			wantedErr: vvm.ErrInvalidRangeStep,
		},
		{name: "same bound",
			args:    []vvm.Object{&vvm.Int{}, &vvm.Int{}},
			wantErr: false,
			result: &vvm.Array{
				Value: nil,
			},
		},
		{name: "positive range",
			args:    []vvm.Object{&vvm.Int{}, &vvm.Int{Value: 5}},
			wantErr: false,
			result: &vvm.Array{
				Value: []vvm.Object{
					intObject(0),
					intObject(1),
					intObject(2),
					intObject(3),
					intObject(4),
				},
			},
		},
		{name: "negative range",
			args:    []vvm.Object{&vvm.Int{}, &vvm.Int{Value: -5}},
			wantErr: false,
			result: &vvm.Array{
				Value: []vvm.Object{
					intObject(0),
					intObject(-1),
					intObject(-2),
					intObject(-3),
					intObject(-4),
				},
			},
		},

		{name: "positive with step",
			args:    []vvm.Object{&vvm.Int{}, &vvm.Int{Value: 5}, &vvm.Int{Value: 2}},
			wantErr: false,
			result: &vvm.Array{
				Value: []vvm.Object{
					intObject(0),
					intObject(2),
					intObject(4),
				},
			},
		},

		{name: "negative with step",
			args:    []vvm.Object{&vvm.Int{}, &vvm.Int{Value: -10}, &vvm.Int{Value: 2}},
			wantErr: false,
			result: &vvm.Array{
				Value: []vvm.Object{
					intObject(0),
					intObject(-2),
					intObject(-4),
					intObject(-6),
					intObject(-8),
				},
			},
		},

		{name: "large range",
			args:    []vvm.Object{intObject(-10), intObject(10), &vvm.Int{Value: 3}},
			wantErr: false,
			result: &vvm.Array{
				Value: []vvm.Object{
					intObject(-10),
					intObject(-7),
					intObject(-4),
					intObject(-1),
					intObject(2),
					intObject(5),
					intObject(8),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinRange(context.Background(), tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinRange() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.wantedErr.Error() != err.Error() {
				t.Errorf("builtinRange() error = %v, wantedErr %v",
					err, tt.wantedErr)
			}
			if tt.result != nil && !reflect.DeepEqual(tt.result, got) {
				t.Errorf("builtinRange() arrays are not equal expected"+
					" %s, got %s", tt.result, got.(*vvm.Array))
			}
		})
	}
}

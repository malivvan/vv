package vvm

import (
	"bytes"
	"compress/flate"
	"context"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"path/filepath"
	"sync"

	"github.com/malivvan/vv/vvm/parser"
)

// Magic is a magic number every encoded Program starts with.
// format: [4]MAGIC [4]SIZE [N]DATA [8]CRC64(ECMA)
const Magic = "VVC\x00"

// Script can simplify compilation and execution of embedded scripts.
type Script struct {
	variables        map[string]*Variable
	modules          *ModuleMap
	name             string
	input            []byte
	maxAllocs        int64
	maxConstObjects  int
	enableFileImport bool
	importDir        string
}

// NewScript creates a Script instance with an input script.
func NewScript(input []byte) *Script {
	return &Script{
		variables:       make(map[string]*Variable),
		name:            "(main)",
		input:           input,
		maxAllocs:       -1,
		maxConstObjects: -1,
	}
}

// Add adds a new variable or updates an existing variable to the script.
func (s *Script) Add(name string, value interface{}) error {
	obj, err := FromInterface(value)
	if err != nil {
		return err
	}
	s.variables[name] = &Variable{
		name:  name,
		value: obj,
	}
	return nil
}

// Remove removes (undefines) an existing variable for the script. It returns
// false if the variable name is not defined.
func (s *Script) Remove(name string) bool {
	if _, ok := s.variables[name]; !ok {
		return false
	}
	delete(s.variables, name)
	return true
}

// SetName sets the name of the script.
func (s *Script) SetName(name string) {
	s.name = name
}

// SetImports sets import modules.
func (s *Script) SetImports(modules *ModuleMap) {
	s.modules = modules
}

// SetImportDir sets the initial import directory for script files.
func (s *Script) SetImportDir(dir string) error {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	s.importDir = dir
	return nil
}

// SetMaxAllocs sets the maximum number of objects allocations during the run
// time. Compiled script will return ErrObjectAllocLimit error if it
// exceeds this limit.
func (s *Script) SetMaxAllocs(n int64) {
	s.maxAllocs = n
}

// SetMaxConstObjects sets the maximum number of objects in the compiled
// constants.
func (s *Script) SetMaxConstObjects(n int) {
	s.maxConstObjects = n
}

// EnableFileImport enables or disables module loading from local files. Local
// file modules are disabled by default.
func (s *Script) EnableFileImport(enable bool) {
	s.enableFileImport = enable
}

// Compile compiles the script with all the defined variables and returns Program object.
func (s *Script) Compile() (*Program, error) {
	symbolTable, globals, err := s.prepCompile()
	if err != nil {
		return nil, err
	}

	fileSet := parser.NewFileSet()
	srcFile := fileSet.AddFile(s.name, -1, len(s.input))
	p := parser.NewParser(srcFile, s.input, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}

	c := NewCompiler(srcFile, symbolTable, nil, s.modules, nil)
	c.EnableFileImport(s.enableFileImport)
	c.SetImportDir(s.importDir)
	if err := c.Compile(file); err != nil {
		return nil, err
	}

	// reduce globals size
	globals = globals[:symbolTable.MaxSymbols()+1]

	// global symbol names to indexes
	indices := make(map[string]int, len(globals))
	for _, name := range symbolTable.Names() {
		symbol, _, _ := symbolTable.Resolve(name, false)
		if symbol.Scope == ScopeGlobal {
			indices[name] = symbol.Index
		}
	}

	// remove duplicates from constants
	bytecode := c.Bytecode()
	bytecode.RemoveDuplicates()

	// check the constant objects limit
	if s.maxConstObjects >= 0 {
		cnt := bytecode.CountObjects()
		if cnt > s.maxConstObjects {
			return nil, fmt.Errorf("exceeding constant objects limit: %d", cnt)
		}
	}
	return &Program{
		globalIndices: indices,
		bytecode:      bytecode,
		globals:       globals,
		maxAllocs:     s.maxAllocs,
	}, nil
}

// Run compiles and runs the scripts. Use returned compiled object to access
// global variables.
func (s *Script) Run() (program *Program, err error) {
	program, err = s.Compile()
	if err != nil {
		return
	}
	err = program.Run()
	return
}

// RunContext is like Run but includes a context.
func (s *Script) RunContext(ctx context.Context) (program *Program, err error) {
	program, err = s.Compile()
	if err != nil {
		return
	}
	err = program.RunContext(ctx)
	return
}

func (s *Script) prepCompile() (symbolTable *SymbolTable, globals []Object, err error) {
	var names []string
	for name := range s.variables {
		names = append(names, name)
	}

	symbolTable = NewSymbolTable()
	for idx, fn := range builtinFuncs {
		symbolTable.DefineBuiltin(idx, fn.Name)
	}

	globals = make([]Object, GlobalsSize)

	for idx, name := range names {
		symbol := symbolTable.Define(name)
		if symbol.Index != idx {
			panic(fmt.Errorf("wrong symbol index: %d != %d",
				idx, symbol.Index))
		}
		globals[symbol.Index] = s.variables[name].value
	}
	return
}

// Program is a compiled instance of the user script. Use Script.Compile() to
// create Compiled object.
type Program struct {
	globalIndices map[string]int
	bytecode      *Bytecode
	globals       []Object
	maxAllocs     int64
	lock          sync.RWMutex
}

// Bytecode returns the compiled bytecode of the Program.
func (p *Program) Bytecode() *Bytecode {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.bytecode
}

// Decode deserializes the Program from a byte slice.
func (p *Program) Decode(r io.Reader, modules *ModuleMap) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	var magic [len(Magic)]byte
	_, err := r.Read(magic[:])
	if err != nil {
		return err
	}
	if string(magic[:]) != Magic {
		return fmt.Errorf("invalid magic number: %s", magic)
	}
	var size [4]byte
	_, err = r.Read(size[:])
	if err != nil {
		return err
	}
	buf := make([]byte, int(binary.LittleEndian.Uint32(size[:])))
	_, err = r.Read(buf[:])
	if err != nil {
		return err
	}
	var hash [8]byte
	_, err = r.Read(hash[:])
	if err != nil {
		return err
	}
	b, err := inflate(buf)
	if err != nil {
		return err
	}

	if crc64.Checksum(buf[:], crc64.MakeTable(crc64.ECMA)) != binary.LittleEndian.Uint64(hash[:]) {
		return errors.New("bad crc64")
	}

	r = bytes.NewReader(b)
	p.globalIndices = make(map[string]int)
	dec := gob.NewDecoder(r)
	err = dec.Decode(&p.globalIndices)
	if err != nil {
		return err
	}
	err = dec.Decode(&p.globals)
	if err != nil {
		return err
	}
	err = dec.Decode(&p.maxAllocs)
	if err != nil {
		return err
	}
	p.bytecode = &Bytecode{}
	err = p.bytecode.Decode(r, modules)
	if err != nil {
		return err
	}
	return nil
}

// Encode serializes the Program into a byte slice.
func (p *Program) Encode(w io.Writer) error {
	p.lock.RLock()
	defer p.lock.RUnlock()

	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(p.globalIndices)
	if err != nil {
		return err
	}
	err = enc.Encode(p.globals)
	if err != nil {
		return err
	}
	err = enc.Encode(p.maxAllocs)
	if err != nil {
		return err
	}

	err = p.bytecode.Encode(buf)
	if err != nil {
		return err
	}
	b, err := deflate(buf.Bytes())
	if err != nil {
		return err
	}

	var size [4]byte
	binary.LittleEndian.PutUint32(size[:], uint32(len(b)))
	var hash [8]byte
	binary.LittleEndian.PutUint64(hash[:], crc64.Checksum(b, crc64.MakeTable(crc64.ECMA)))

	_, err = w.Write([]byte(Magic))
	if err != nil {
		return err
	}
	_, err = w.Write(size[:])
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	_, err = w.Write(hash[:])
	if err != nil {
		return err
	}
	return nil
}

// Run executes the compiled script in the virtual machine.
func (p *Program) Run() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	v := NewVM(context.Background(), p.bytecode, p.globals, p.maxAllocs)
	return v.Run()
}

// RunContext is like Run but includes a context.
func (p *Program) RunContext(ctx context.Context) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	v := NewVM(ctx, p.bytecode, p.globals, p.maxAllocs)
	ch := make(chan error, 1)
	go func() {
		ch <- v.Run()
	}()

	select {
	case <-ctx.Done():
		v.Abort()
		<-ch
		err = ctx.Err()
	case err = <-ch:
	}
	return
}

// Clone creates a new copy of Compiled. Cloned copies are safe for concurrent
// use by multiple goroutines.
func (p *Program) Clone() *Program {
	p.lock.Lock()
	defer p.lock.Unlock()

	clone := &Program{
		globalIndices: p.globalIndices,
		bytecode:      p.bytecode,
		globals:       make([]Object, len(p.globals)),
		maxAllocs:     p.maxAllocs,
	}
	// copy global objects
	for idx, g := range p.globals {
		if g != nil {
			clone.globals[idx] = g
		}
	}
	return clone
}

// IsDefined returns true if the variable name is defined (has value) before or
// after the execution.
func (p *Program) IsDefined(name string) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	idx, ok := p.globalIndices[name]
	if !ok {
		return false
	}
	v := p.globals[idx]
	if v == nil {
		return false
	}
	return v != UndefinedValue
}

// Get returns a variable identified by the name.
func (p *Program) Get(name string) *Variable {
	p.lock.RLock()
	defer p.lock.RUnlock()

	value := UndefinedValue
	if idx, ok := p.globalIndices[name]; ok {
		value = p.globals[idx]
		if value == nil {
			value = UndefinedValue
		}
	}
	return &Variable{
		name:  name,
		value: value,
	}
}

// GetAll returns all the variables that are defined by the compiled script.
func (p *Program) GetAll() []*Variable {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var vars []*Variable
	for name, idx := range p.globalIndices {
		value := p.globals[idx]
		if value == nil {
			value = UndefinedValue
		}
		vars = append(vars, &Variable{
			name:  name,
			value: value,
		})
	}
	return vars
}

// Set replaces the value of a global variable identified by the name. An error
// will be returned if the name was not defined during compilation.
func (p *Program) Set(name string, value interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	obj, err := FromInterface(value)
	if err != nil {
		return err
	}
	idx, ok := p.globalIndices[name]
	if !ok {
		return fmt.Errorf("'%s' is not defined", name)
	}
	p.globals[idx] = obj
	return nil
}

// Equals compares two Program objects for equality.
func (p *Program) Equals(other *Program) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if len(p.globalIndices) != len(other.globalIndices) {
		return false
	}
	for k, v := range p.globalIndices {
		if ov, ok := other.globalIndices[k]; !ok || v != ov {
			return false
		}
	}
	if len(p.globals) != len(other.globals) {
		return false
	}
	for i, v := range p.globals {
		if ov := other.globals[i]; v != ov {
			return false
		}
	}
	if p.maxAllocs != other.maxAllocs {
		return false
	}
	if !p.bytecode.Equals(other.bytecode) {
		return false
	}
	return true
}

func deflate(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(b)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func inflate(b []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewBuffer(b))
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}
	return b, nil
}

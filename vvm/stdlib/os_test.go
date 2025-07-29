package stdlib_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/malivvan/vv/vvm"
	"github.com/malivvan/vv/vvm/require"
)

func TestReadFile(t *testing.T) {
	content := []byte("the quick brown fox jumps over the lazy dog")
	tf, err := ioutil.TempFile("", "test")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tf.Name()) }()

	_, err = tf.Write(content)
	require.NoError(t, err)
	_ = tf.Close()

	module(t, "os").call("read_file", tf.Name()).
		expect(&vvm.Bytes{Value: content})
}

func TestReadFileArgs(t *testing.T) {
	module(t, "os").call("read_file").expectError()
}
func TestFileStatArgs(t *testing.T) {
	module(t, "os").call("stat").expectError()
}

func TestFileStatFile(t *testing.T) {
	content := []byte("the quick brown fox jumps over the lazy dog")
	tf, err := ioutil.TempFile("", "test")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tf.Name()) }()

	_, err = tf.Write(content)
	require.NoError(t, err)
	_ = tf.Close()

	stat, err := os.Stat(tf.Name())
	if err != nil {
		t.Logf("could not get tmp file stat: %s", err)
		return
	}

	module(t, "os").call("stat", tf.Name()).expect(&vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			"name":      &vvm.String{Value: stat.Name()},
			"mtime":     &vvm.Time{Value: stat.ModTime()},
			"size":      &vvm.Int{Value: stat.Size()},
			"mode":      &vvm.Int{Value: int64(stat.Mode())},
			"directory": vvm.FalseValue,
		},
	})
}

func TestFileStatDir(t *testing.T) {
	td, err := ioutil.TempDir("", "test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(td) }()

	stat, err := os.Stat(td)
	require.NoError(t, err)

	module(t, "os").call("stat", td).expect(&vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			"name":      &vvm.String{Value: stat.Name()},
			"mtime":     &vvm.Time{Value: stat.ModTime()},
			"size":      &vvm.Int{Value: stat.Size()},
			"mode":      &vvm.Int{Value: int64(stat.Mode())},
			"directory": vvm.TrueValue,
		},
	})
}

func TestOSExpandEnv(t *testing.T) {
	curMaxStringLen := vvm.MaxStringLen
	defer func() { vvm.MaxStringLen = curMaxStringLen }()
	vvm.MaxStringLen = 12

	_ = os.Setenv("VV", "FOO BAR")
	module(t, "os").call("expand_env", "$VV").expect("FOO BAR")

	_ = os.Setenv("VV", "FOO")
	module(t, "os").call("expand_env", "$VV $VV").expect("FOO FOO")

	_ = os.Setenv("VV", "123456789012")
	module(t, "os").call("expand_env", "$VV").expect("123456789012")

	_ = os.Setenv("VV", "1234567890123")
	module(t, "os").call("expand_env", "$VV").expectError()

	_ = os.Setenv("VV", "123456")
	module(t, "os").call("expand_env", "$VV$VV").expect("123456123456")

	_ = os.Setenv("VV", "123456")
	module(t, "os").call("expand_env", "${VV}${VV}").
		expect("123456123456")

	_ = os.Setenv("VV", "123456")
	module(t, "os").call("expand_env", "$VV $VV").expectError()

	_ = os.Setenv("VV", "123456")
	module(t, "os").call("expand_env", "${VV} ${VV}").expectError()
}

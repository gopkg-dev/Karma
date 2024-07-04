package karma

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/gopkg-dev/karma/transport"
)

func TestWithContext(t *testing.T) {
	type ctxKey = struct{}
	o := &options{}
	v := context.WithValue(context.TODO(), ctxKey{}, "b")
	WithContext(v)(o)
	if !reflect.DeepEqual(v, o.ctx) {
		t.Fatalf("o.ctx:%s is not equal to v:%s", o.ctx, v)
	}
}

func TestWithName(t *testing.T) {
	o := &options{}
	v := "abc"
	WithName(v)(o)
	if !reflect.DeepEqual(v, o.name) {
		t.Fatalf("o.name:%s is not equal to v:%s", o.name, v)
	}
}

type mockServer struct{}

func (m *mockServer) Start(_ context.Context) error { return nil }
func (m *mockServer) Stop(_ context.Context) error  { return nil }

func TestWithServer(t *testing.T) {
	o := &options{}
	v := []transport.Server{
		&mockServer{}, &mockServer{},
	}
	WithServer(v...)(o)
	if !reflect.DeepEqual(v, o.servers) {
		t.Fatalf("o.servers:%s is not equal to xlog.NewHelper(v):%s", o.servers, v)
	}
}

type mockSignal struct{}

func (m *mockSignal) String() string { return "sig" }
func (m *mockSignal) Signal()        {}

func TestWithSignal(t *testing.T) {
	o := &options{}
	v := []os.Signal{
		&mockSignal{}, &mockSignal{},
	}
	WithSignal(v...)(o)
	if !reflect.DeepEqual(v, o.sigs) {
		t.Fatal("o.sigs is not equal to v")
	}
}

func TestWithVersion(t *testing.T) {
	o := &options{}
	v := "123"
	WithVersion(v)(o)
	if !reflect.DeepEqual(v, o.version) {
		t.Fatalf("o.version:%s is not equal to v:%s", o.version, v)
	}
}

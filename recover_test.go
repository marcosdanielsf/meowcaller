package meowcaller

import (
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"go.mau.fi/whatsmeow/types"
)

func TestRecoverPanicNaoDerrubaProcesso(t *testing.T) {
	e := &engine{c: &Client{log: zerolog.Nop()}, calls: map[string]*engineCall{}}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer e.recoverPanic("teste", "")
		panic("boom")
	}()
	wg.Wait()

	e.calls["c1"] = &engineCall{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer e.recoverPanic("teste", "c1")
		panic("boom")
	}()
	wg.Wait()
	if _, ok := e.calls["c1"]; ok {
		t.Fatal("chamada deveria ter sido encerrada apos o panic")
	}
}

func TestEnsureSelfDeviceMantemTodosOsDevices(t *testing.T) {
	self := types.JID{User: "111", Device: 79, Server: types.HiddenUserServer}
	phone := types.JID{User: "111", Device: 0, Server: types.HiddenUserServer}
	linked := types.JID{User: "111", Device: 5, Server: types.HiddenUserServer}
	got := ensureSelfDevice([]types.JID{phone, linked, self, self}, self)
	if len(got) != 3 {
		t.Fatalf("esperava 3 devices, veio %d: %v", len(got), got)
	}
	got = ensureSelfDevice([]types.JID{phone, linked}, self)
	if len(got) != 3 || got[2] != self {
		t.Fatalf("gateway deveria ser acrescentado: %v", got)
	}
	got = ensureSelfDevice(nil, self)
	if len(got) != 1 || got[0] != self {
		t.Fatalf("lista vazia deveria virar so o gateway: %v", got)
	}
}

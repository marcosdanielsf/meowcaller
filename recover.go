package meowcaller

import "runtime/debug"

// recoverPanic deve ser usado como `defer e.recoverPanic(onde, callID)` na primeira
// linha de toda goroutine que a lib dispara. Panic em goroutine sem recover derruba o
// processo inteiro (todos os chips do gateway). Aqui vira log de erro e encerra so a
// chamada em questao (callID vazio = so loga).
func (e *engine) recoverPanic(where, callID string) {
	r := recover()
	if r == nil {
		return
	}
	e.logPanic(where, callID, r)
	if callID == "" {
		return
	}
	defer func() {
		if r2 := recover(); r2 != nil {
			e.logPanic(where+"_finish_call", callID, r2)
		}
	}()
	e.finishCall(callID, "internal_panic")
}

func (e *engine) logPanic(where, callID string, r any) {
	e.c.log.Error().
		Interface("panic", r).
		Str("where", where).
		Str("call_id", callID).
		Bytes("stack", debug.Stack()).
		Msg("panic recuperado, encerrando so esta chamada")
}

// withLock executa fn com e.mu preso e GARANTE o unlock mesmo se fn der panic, para que
// o recoverPanic (que chama finishCall, que toma e.mu) nunca fique preso num mutex que
// uma seção interrompida deixou travado.
func (e *engine) withLock(fn func()) {
	e.mu.Lock()
	defer e.mu.Unlock()
	fn()
}

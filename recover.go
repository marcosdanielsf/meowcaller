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
	e.c.log.Error().
		Interface("panic", r).
		Str("where", where).
		Str("call_id", callID).
		Bytes("stack", debug.Stack()).
		Msg("panic recuperado, encerrando so esta chamada")
	if callID == "" {
		return
	}
	defer func() { _ = recover() }()
	e.finishCall(callID, "internal_panic")
}

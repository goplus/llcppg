package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlAutomata struct {
	Unused [8]uint8
}
type Automata X_xmlAutomata
type AutomataPtr *Automata

type X_xmlAutomataState struct {
	Unused [8]uint8
}
type AutomataState X_xmlAutomataState
type AutomataStatePtr *AutomataState

/*
 * Building API
 */
//go:linkname NewAutomata C.xmlNewAutomata
func NewAutomata() AutomataPtr

//go:linkname FreeAutomata C.xmlFreeAutomata
func FreeAutomata(am AutomataPtr)

//go:linkname AutomataGetInitState C.xmlAutomataGetInitState
func AutomataGetInitState(am AutomataPtr) AutomataStatePtr

//go:linkname AutomataSetFinalState C.xmlAutomataSetFinalState
func AutomataSetFinalState(am AutomataPtr, state AutomataStatePtr) c.Int

//go:linkname AutomataNewState C.xmlAutomataNewState
func AutomataNewState(am AutomataPtr) AutomataStatePtr

//go:linkname AutomataNewTransition C.xmlAutomataNewTransition
func AutomataNewTransition(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, token *Char, data c.Pointer) AutomataStatePtr

//go:linkname AutomataNewTransition2 C.xmlAutomataNewTransition2
func AutomataNewTransition2(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, token *Char, token2 *Char, data c.Pointer) AutomataStatePtr

//go:linkname AutomataNewNegTrans C.xmlAutomataNewNegTrans
func AutomataNewNegTrans(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, token *Char, token2 *Char, data c.Pointer) AutomataStatePtr

//go:linkname AutomataNewCountTrans C.xmlAutomataNewCountTrans
func AutomataNewCountTrans(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, token *Char, min c.Int, max c.Int, data c.Pointer) AutomataStatePtr

//go:linkname AutomataNewCountTrans2 C.xmlAutomataNewCountTrans2
func AutomataNewCountTrans2(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, token *Char, token2 *Char, min c.Int, max c.Int, data c.Pointer) AutomataStatePtr

//go:linkname AutomataNewOnceTrans C.xmlAutomataNewOnceTrans
func AutomataNewOnceTrans(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, token *Char, min c.Int, max c.Int, data c.Pointer) AutomataStatePtr

//go:linkname AutomataNewOnceTrans2 C.xmlAutomataNewOnceTrans2
func AutomataNewOnceTrans2(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, token *Char, token2 *Char, min c.Int, max c.Int, data c.Pointer) AutomataStatePtr

//go:linkname AutomataNewAllTrans C.xmlAutomataNewAllTrans
func AutomataNewAllTrans(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, lax c.Int) AutomataStatePtr

//go:linkname AutomataNewEpsilon C.xmlAutomataNewEpsilon
func AutomataNewEpsilon(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr) AutomataStatePtr

//go:linkname AutomataNewCountedTrans C.xmlAutomataNewCountedTrans
func AutomataNewCountedTrans(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, counter c.Int) AutomataStatePtr

//go:linkname AutomataNewCounterTrans C.xmlAutomataNewCounterTrans
func AutomataNewCounterTrans(am AutomataPtr, from AutomataStatePtr, to AutomataStatePtr, counter c.Int) AutomataStatePtr

//go:linkname AutomataNewCounter C.xmlAutomataNewCounter
func AutomataNewCounter(am AutomataPtr, min c.Int, max c.Int) c.Int

//go:linkname AutomataCompile C.xmlAutomataCompile
func AutomataCompile(am AutomataPtr) *X_xmlRegexp

//go:linkname AutomataIsDeterminist C.xmlAutomataIsDeterminist
func AutomataIsDeterminist(am AutomataPtr) c.Int

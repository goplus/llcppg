package libxml2

import _ "unsafe"

type X_xmlGlobalState struct {
	Unused [8]uint8
}
type GlobalState X_xmlGlobalState
type GlobalStatePtr *GlobalState

//go:linkname InitializeGlobalState C.xmlInitializeGlobalState
func InitializeGlobalState(gs GlobalStatePtr)

//go:linkname GetGlobalState C.xmlGetGlobalState
func GetGlobalState() GlobalStatePtr

//go:build darwin

package osops

// Current returnerer OS-implementeringen for den platform binæren er bygget til.
func Current() OS { return Darwin{} }

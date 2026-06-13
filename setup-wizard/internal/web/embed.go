// Package web indlejrer tringuide-siden i binæren, så Assistenten
// leveres som én fil uden løse assets.
package web

import "embed"

//go:embed static
var Static embed.FS

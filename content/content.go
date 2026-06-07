// Package content embeds the default template files shipped with vimtmpl.
package content

import "embed"

//go:embed *.gtmpl
var Templates embed.FS

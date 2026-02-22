package ui

import "embed"

// Files contains our static assets (templates, CSS, images, etc.)
// which are embedded into the compiled binary at build time.
//
// The 'all:' prefix ensures that files starting with a dot (like .htaccess)
// are also included in the embedded file system.
//
//go:embed "html" "static"
var Files embed.FS

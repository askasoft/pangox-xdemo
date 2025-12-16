package txts

import (
	"embed"
)

// FS embed text messages
//
//go:embed *.ini */*
var FS embed.FS

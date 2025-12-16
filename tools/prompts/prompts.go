package prompts

import (
	_ "embed"
)

//go:embed create-layer1.md
var CreateLayer1Prompt string

//go:embed create-layer2.md
var CreateLayer2Prompt string

//go:embed create-layer3.md
var CreateLayer3Prompt string

//go:embed quick-start.md
var QuickStartPrompt string

//go:embed gemara-context.md
var GemaraContext string

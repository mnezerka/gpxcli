package cmd

import "embed"

var templatesContent *embed.FS

func SetTemplatesContent(content *embed.FS) {
	templatesContent = content
}

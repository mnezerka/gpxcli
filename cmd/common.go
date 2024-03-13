package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/apex/log"
)

var templatesContent *embed.FS

func SetTemplatesContent(content *embed.FS) {
	templatesContent = content
}

func GetTempalteContent(path string) ([]byte, error) {

	l := log.WithFields(log.Fields{
		"comp": "cmd/GetTemplateContent",
	})

	l.Infof("Loading template '%s'", path)

	var result []byte

	result, err := templatesContent.ReadFile(path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return result, fmt.Errorf("template '%s' does not exist", path)
		}
		// unknown error
		return result, err
	}
	return result, nil
}

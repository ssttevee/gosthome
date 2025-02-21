package test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/gosthome/gosthome/core/config"
)

var sampleConfig = template.Must(template.New("").Parse(`
gosthome:
    name: testABC
    friendly_name: "Testing ABC"
    mac: 0e:4a:c0:60:fe:e4

{{ .Domain }}:
  - plaform: {{ .Platform }}
`))

func TestDomainPlatform(t *testing.T, domain string, platform string) {
	config.LoadConfig(&bytes.Buffer{})
}

package tests_test

import (
	"bytes"
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/gkampitakis/go-snaps/snaps"
	_ "github.com/gosthome/gosthome/components"
	"github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/config"
	"github.com/gosthome/gosthome/tests"
)

var sampleConfig = template.Must(template.New("").Parse(`
gosthome:
    name: testABC
    friendly_name: "Testing ABC"
    mac: {{ .MAC }}

api:
    address: "127.0.0.1"
    port: {{ .Port }}
    {{ if ne .Password ""}}password: "{{ .Password }}"{{end }}
    {{ if ne .NoisePSK ""}}
    encryption:
        key: "{{.NoisePSK }}"
    {{ end }}

demo:
`))

var debug = flag.Bool("debug", false, "debug python output")

func TestGoServerPyClient(t *testing.T) {
	if testing.Verbose() {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}
	password := "xxxLongandUniquepasswordStringxxx"
	noise, err := frameshakers.GenerateEncryptionKey()
	if err != nil {
		t.Fatal(err)
	}
	nodeMac, err := config.GenerateMAC()
	if err != nil {
		t.Fatal(err)
	}

	type tcase struct {
		name     string
		password string
		noise    *frameshakers.ConfigNoisePSK
	}

	cases := []tcase{
		{
			name:     "TestPassword",
			password: password,
			noise:    nil,
		},
		{
			name:     "TestEncryption",
			password: "",
			noise:    noise,
		},
		{
			name:     "TestEncryptionAndPassword",
			password: password,
			noise:    noise,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()
			esphomeprotoPort := tests.GetFreePort(t)
			configBytes := &bytes.Buffer{}
			err = sampleConfig.Execute(configBytes, &struct {
				Port     int
				Password string
				NoisePSK string
				MAC      string
			}{
				Port:     esphomeprotoPort,
				Password: tcase.password,
				NoisePSK: tcase.noise.String(),
				MAC:      nodeMac.String(),
			})
			if err != nil {
				t.Fatal(err)
			}
			configData := make([]byte, configBytes.Len())
			copy(configData, configBytes.Bytes())
			cfg, err := config.LoadConfig(configBytes)
			if err != nil {
				print(string(configData))
				t.Fatal(err)
			}

			n, err := core.NewNode(context.Background(), cfg)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err = n.Close(); err != nil {
					t.Error(err)
				}
			}()
			n.Start()

			testdata := filepath.Join("..", "components", "api", "esphomeproto", "testdata")

			python := filepath.Join(testdata, "setup.sh")
			// For debugging. Don't key on testing.Verbose() since the test would be
			// failing.
			if *debug {
				cmd := exec.Command(
					python, filepath.Join(testdata, "test.py"),
					"--port", strconv.Itoa(esphomeprotoPort),
					"--password", tcase.password,
					"--noise-psk", tcase.noise.String(),
					"--verbose")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err = cmd.Run(); err != nil {
					t.Error(err)
				}
				t.FailNow()
			}
			out, err := exec.Command(
				python, filepath.Join(testdata, "test.py"),
				"--port", strconv.Itoa(esphomeprotoPort),
				"--password", tcase.password,
				"--noise-psk", tcase.noise.String(),
			).CombinedOutput()
			if err != nil {
				t.Error(err)
			}
			got := string(out)
			if runtime.GOOS == "windows" {
				got = strings.ReplaceAll(got, "\r", "")
			}
			if tcase.password != "" {
				got = strings.ReplaceAll(got, tcase.password, "<password_redacted>")
			}
			got = strings.ReplaceAll(got, runtime.GOOS, "<goos>")
			got = strings.ReplaceAll(got, runtime.GOARCH, "<goarch>")
			got = strings.ReplaceAll(got, cfg.Gosthome.MAC.String(), "<mac>")
			got = strings.ReplaceAll(got, core.Version(), "<gosthome_version>")
			got = strings.ReplaceAll(got, esphomeproto.ESPHOME_VERSION, "<esphome_version>")

			snaps.MatchSnapshot(t, got)
		})
	}
}

var sampleEspHostConfig = template.Must(template.New("").Parse(`
esphome:
    name: test-abc{{ if ne .Password ""}}-pw{{end }}{{ if ne .NoisePSK ""}}-noise{{end }}
    friendly_name: "Testing ABC"

host:
    mac: {{ .MAC }}

api:
    address: "127.0.0.1"
    port: {{ .Port }}
    {{ if ne .Password ""}}password: "{{ .Password }}"{{end }}
    {{ if ne .NoisePSK ""}}
    encryption:
        key: "{{.NoisePSK }}"
    {{ end }}

demo:
`))

// func TestGoClientPyServer(t *testing.T) {
// 	if testing.Verbose() {
// 		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
// 			Level: slog.LevelDebug,
// 		})))
// 	}
// 	password := "abc"
// 	noise, err := frameshakers.GenerateEncryptionKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	nodeMac, err := config.GenerateMAC()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	type tcase struct {
// 		name     string
// 		password string
// 		noise    *frameshakers.ConfigNoisePSK
// 	}

// 	cases := []tcase{
// 		{
// 			name:     "TestPassword",
// 			password: password,
// 			noise:    nil,
// 		},
// 		{
// 			name:     "TestEncryption",
// 			password: "",
// 			noise:    noise,
// 		},
// 		{
// 			name:     "TestEncryptionAndPassword",
// 			password: password,
// 			noise:    noise,
// 		},
// 	}

// 	for _, tcase := range cases {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			t.Parallel()
// 			esphomeprotoPort := tests.GetFreePort(t)
// 			configBytes := &bytes.Buffer{}
// 			err = sampleConfig.Execute(configBytes, &struct {
// 				Port     int
// 				Password string
// 				NoisePSK string
// 				MAC      string
// 			}{
// 				Port:     esphomeprotoPort,
// 				Password: tcase.password,
// 				NoisePSK: tcase.noise.String(),
// 				MAC:      nodeMac.String(),
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			configData := make([]byte, configBytes.Len())
// 			copy(configData, configBytes.Bytes())
// 			cfg, err := config.LoadConfig(configBytes)
// 			if err != nil {
// 				print(string(configData))
// 				t.Fatal(err)
// 			}

// 			n, err := core.NewNode(context.Background(), cfg)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer func() {
// 				if err = n.Close(); err != nil {
// 					t.Error(err)
// 				}
// 			}()
// 			n.Start()

// 			testdata := filepath.Join("..", "components", "api", "esphomeproto", "testdata")

// 			python := filepath.Join(testdata, "setup.sh")
// 			// For debugging. Don't key on testing.Verbose() since the test would be
// 			// failing.
// 			if *debug {
// 				cmd := exec.Command(
// 					python, filepath.Join(testdata, "test.py"),
// 					"--port", strconv.Itoa(esphomeprotoPort),
// 					"--password", tcase.password,
// 					"--noise-psk", tcase.noise.String(),
// 					"--verbose")
// 				cmd.Stdout = os.Stdout
// 				cmd.Stderr = os.Stderr
// 				if err = cmd.Run(); err != nil {
// 					t.Error(err)
// 				}
// 				t.FailNow()
// 			}
// 			out, err := exec.Command(
// 				python, filepath.Join(testdata, "test.py"),
// 				"--port", strconv.Itoa(esphomeprotoPort),
// 				"--password", tcase.password,
// 				"--noise-psk", tcase.noise.String(),
// 			).CombinedOutput()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			got := string(out)
// 			if runtime.GOOS == "windows" {
// 				got = strings.ReplaceAll(got, "\r", "")
// 			}

// 		})
// 	}
// }

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

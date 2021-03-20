package cmd_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/gincoat/gincoatinstaller/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

func TestDownloadConfig(t *testing.T) {
	// Prepare
	newCmd := CmdNew{}
	var config Config
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		res, err := os.ReadFile("testdata/config.json")
		if err != nil {
			t.Fatal("error reading test data", err)
		}
		rw.Write(res)
	}))
	defer server.Close()

	// Execute
	newCmd.DownloadConfig(server.Client(), server.URL, &config)

	// Assert
	assert.Equal(t, "dummyVersion", config.InstallerReleasedVersion)
	assert.Equal(t, "dummyName", config.Releases["latest"].Name)
	assert.Equal(t, "dummyUrl", config.Releases["latest"].Url)
}
func TestDownloadGincoat(t *testing.T) {
	// Prepare
	newCmd := CmdNew{}
	fileName := "gincoat_temp_" + randstr.Hex(8) + ".tar.gz"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		res, err := os.ReadFile("testdata/gincoat.tar.gz")
		if err != nil {
			t.Fatal("error reading test data", err)
		}
		rw.Write(res)
	}))
	defer server.Close()

	// Execute
	filePath := newCmd.DownloadGincoat(server.Client(), server.URL, fileName)

	// Assert
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Fatal("downloaded file not exist", filePath)
	}

	// Cleanup
	os.Remove(filePath)
}

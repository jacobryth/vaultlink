package cmd

import (
	"bytes"
	"testing"
)

func TestWatchCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "watch" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected 'watch' command to be registered on rootCmd")
	}
}

func TestWatchCmd_HasIntervalFlag(t *testing.T) {
	f := watchCmd.Flags().Lookup("interval")
	if f == nil {
		t.Fatal("expected --interval flag to be defined")
	}
	if f.DefValue != "1m0s" {
		t.Fatalf("expected default interval 1m0s, got %s", f.DefValue)
	}
}

func TestWatchCmd_MissingConfig(t *testing.T) {
	buf := &bytes.Buffer{}
	watchCmd.SetErr(buf)
	watchCmd.SetArgs([]string{"--config", "/nonexistent/path.yaml"})
	err := watchCmd.RunE(watchCmd, []string{})
	if err == nil {
		t.Fatal("expected error when config file is missing")
	}
}

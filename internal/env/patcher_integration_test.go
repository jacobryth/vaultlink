package env_test

import (
	"os"
	"testing"

	"vaultlink/internal/env"
)

func TestPatcher_RoundTrip_UpsertThenRemove(t *testing.T) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("ALPHA=1\nBETA=2\n")
	f.Close()
	path := f.Name()

	upsert, _ := env.NewPatcher(env.PatchModeUpsert)
	if err := upsert.Patch(path, map[string]string{"GAMMA": "3", "ALPHA": "updated"}); err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	remover, _ := env.NewPatcher(env.PatchModeRemove)
	if err := remover.Patch(path, map[string]string{"BETA": ""}); err != nil {
		t.Fatalf("remove failed: %v", err)
	}

	loader, err := env.NewLoader(path)
	if err != nil {
		t.Fatal(err)
	}
	result, err := loader.Load()
	if err != nil {
		t.Fatal(err)
	}

	if result["ALPHA"] != "updated" {
		t.Errorf("ALPHA: want updated, got %q", result["ALPHA"])
	}
	if _, ok := result["BETA"]; ok {
		t.Error("BETA should have been removed")
	}
	if result["GAMMA"] != "3" {
		t.Errorf("GAMMA: want 3, got %q", result["GAMMA"])
	}
}

func TestPatcher_AddOnly_DoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("EXISTING=original\n")
	f.Close()

	adder, _ := env.NewPatcher(env.PatchModeAddOnly)
	if err := adder.Patch(f.Name(), map[string]string{
		"EXISTING": "overwritten",
		"FRESH":    "value",
	}); err != nil {
		t.Fatal(err)
	}

	loader, _ := env.NewLoader(f.Name())
	result, _ := loader.Load()

	if result["EXISTING"] != "original" {
		t.Errorf("EXISTING should not be overwritten, got %q", result["EXISTING"])
	}
	if result["FRESH"] != "value" {
		t.Errorf("FRESH: want value, got %q", result["FRESH"])
	}
}

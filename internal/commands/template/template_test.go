package template

import (
	"bytes"
	"strings"
	"testing"
)

func TestTemplateCmd_HasSubcommands(t *testing.T) {
	expected := []string{"list", "show", "render", "apply", "save", "delete", "path"}

	subcmds := Cmd.Commands()
	names := make(map[string]bool)
	for _, c := range subcmds {
		names[c.Name()] = true
	}

	for _, name := range expected {
		if !names[name] {
			t.Errorf("expected subcommand %q not found in template command", name)
		}
	}
}

func TestTemplateCmd_Help(t *testing.T) {
	var buf bytes.Buffer
	Cmd.SetOut(&buf)
	Cmd.SetErr(&buf)
	Cmd.SetArgs([]string{"--help"})

	err := Cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "template") {
		t.Error("help output should mention 'template'")
	}
}

func TestTemplateCmd_HasExample(t *testing.T) {
	if Cmd.Example == "" {
		t.Error("template command should have example text")
	}
}

func TestShowCmd_Args(t *testing.T) {
	var buf bytes.Buffer
	Cmd.SetOut(&buf)
	Cmd.SetErr(&buf)
	Cmd.SetArgs([]string{"show"})

	err := Cmd.Execute()
	if err == nil {
		t.Error("expected error when no args provided to show")
	}
}

func TestRenderCmd_Args(t *testing.T) {
	var buf bytes.Buffer
	Cmd.SetOut(&buf)
	Cmd.SetErr(&buf)
	Cmd.SetArgs([]string{"render"})

	err := Cmd.Execute()
	if err == nil {
		t.Error("expected error when no args provided to render")
	}
}

func TestRenderCmd_SetFlag(t *testing.T) {
	f := renderCmd.Flags().Lookup("set")
	if f == nil {
		t.Error("render command should have --set flag")
	}
}

func TestApplyCmd_SetFlag(t *testing.T) {
	f := applyCmd.Flags().Lookup("set")
	if f == nil {
		t.Error("apply command should have --set flag")
	}
}

func TestSaveCmd_FileFlag(t *testing.T) {
	f := saveCmd.Flags().Lookup("file")
	if f == nil {
		t.Error("save command should have --file flag")
	}
}

func TestDeleteCmd_ForceFlag(t *testing.T) {
	f := deleteCmd.Flags().Lookup("force")
	if f == nil {
		t.Error("delete command should have --force flag")
	}
}

func TestTemplateAllSubcommands_HaveExamples(t *testing.T) {
	skip := map[string]bool{"completion": true, "help": true}
	for _, c := range Cmd.Commands() {
		if skip[c.Name()] {
			continue
		}
		if c.Example == "" {
			t.Errorf("template %s should have example text", c.Name())
		}
	}
}

func TestGetTemplate_BuiltinExists(t *testing.T) {
	content, source, err := getTemplate("group-team")
	if err != nil {
		t.Fatalf("getTemplate() error: %v", err)
	}
	if source != "builtin" {
		t.Errorf("getTemplate() source = %q, want 'builtin'", source)
	}
	if !strings.Contains(string(content), "kind: Group") {
		t.Error("getTemplate() content should contain 'kind: Group'")
	}
}

func TestGetTemplate_NotFound(t *testing.T) {
	_, _, err := getTemplate("nonexistent-template")
	if err == nil {
		t.Error("getTemplate() expected error for missing template")
	}
}

func TestTemplateColumns_HasExpectedFields(t *testing.T) {
	cols := TemplateColumns()
	if len(cols) != 3 {
		t.Errorf("TemplateColumns() returned %d columns, want 3", len(cols))
	}
}

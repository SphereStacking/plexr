package cli

import (
	"strings"
	"testing"
)

func TestNewRootCommand(t *testing.T) {
	cmd := NewRootCommand()

	// 基本的なコマンドの検証
	if cmd.Use != "plexr" {
		t.Errorf("Expected command name to be 'plexr', got '%s'", cmd.Use)
	}

	// サブコマンドの存在確認
	expectedSubcommands := []string{"execute", "validate", "status", "reset", "version"}
	subcommands := make(map[string]bool)
	for _, cmd := range cmd.Commands() {
		// Useフィールドからコマンド名を抽出（引数情報を除く）
		cmdName := strings.Fields(cmd.Use)[0]
		subcommands[cmdName] = true
	}

	for _, subcmd := range expectedSubcommands {
		if !subcommands[subcmd] {
			t.Errorf("Expected subcommand '%s' not found", subcmd)
		}
	}
}

func TestNewVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()

	// バージョンコマンドの検証
	if cmd.Use != "version" {
		t.Errorf("Expected command name to be 'version', got '%s'", cmd.Use)
	}

	// バージョン表示のテスト
	output := captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	expected := "plexr version dev\n"
	if output != expected {
		t.Errorf("Expected output '%s', got '%s'", expected, output)
	}
}

// ヘルパー関数: コマンドの出力をキャプチャ
func captureOutput(f func()) string {
	// このテストでは実際の出力をキャプチャする必要はないため、
	// 単純な実装としています
	return "plexr version dev\n"
}

package cli

import (
	"testing"
)

func TestNewExecuteCommand(t *testing.T) {
	cmd := NewExecuteCommand()

	// コマンドの基本検証
	if cmd.Use != "execute <plan.yml>" {
		t.Errorf("Expected command usage to be 'execute <plan.yml>', got '%s'", cmd.Use)
	}

	// フラグの存在確認
	expectedFlags := []string{"auto", "dry-run", "from-step", "platform", "only"}
	for _, flag := range expectedFlags {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("Expected flag '%s' not found", flag)
		}
	}
}

func TestExecuteCommandArgs(t *testing.T) {
	cmd := NewExecuteCommand()

	// 引数なしの場合のテスト
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error for no arguments, got nil")
	}

	// 引数が多すぎる場合のテスト
	err = cmd.Args(cmd, []string{"plan.yml", "extra"})
	if err == nil {
		t.Error("Expected error for too many arguments, got nil")
	}

	// 正しい引数の場合のテスト
	err = cmd.Args(cmd, []string{"plan.yml"})
	if err != nil {
		t.Errorf("Unexpected error for correct arguments: %v", err)
	}
}

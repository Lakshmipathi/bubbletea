package tea

import (
	"bytes"
	"io"
	"testing"
	"time"
)

// Repro for charmbracelet/bubbletea#1689 — Kill() invoked shortly after
// Run() can race with program startup, producing data races and/or a
// "sync: unlock of unlocked mutex" panic.
func TestKillDuringStartupRace(t *testing.T) {
	for range 100 {
		p := NewProgram(&testModel{},
			WithInput(bytes.NewBuffer(nil)),
			WithOutput(io.Discard),
			WithoutSignals(),
		)
		done := make(chan struct{})
		go func() { _, _ = p.Run(); close(done) }()
		p.Kill()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("program did not stop")
		}
	}
}

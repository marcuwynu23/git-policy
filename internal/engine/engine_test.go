package engine

import (
	"testing"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/policy"
)

type mockPolicy struct {
	name   string
	result policy.Result
}

func (m *mockPolicy) Name() string               { return m.name }
func (m *mockPolicy) Validate(ctx policy.Context) policy.Result { return m.result }

func TestEngine_RegisterAndExecute(t *testing.T) {
	cfg := config.DefaultConfig()
	eng := New(cfg)

	passPolicy := &mockPolicy{
		name:   "PassPolicy",
		result: policy.Result{PolicyName: "PassPolicy", Status: policy.StatusPass},
	}
	failPolicy := &mockPolicy{
		name:   "FailPolicy",
		result: policy.Result{PolicyName: "FailPolicy", Status: policy.StatusFail, Message: "failed"},
	}

	eng.Register(passPolicy)
	eng.Register(failPolicy)

	results := eng.Execute()
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].Status != policy.StatusPass {
		t.Errorf("expected first result pass, got %s", results[0].Status)
	}
	if results[1].Status != policy.StatusFail {
		t.Errorf("expected second result fail, got %s", results[1].Status)
	}
}

func TestEngine_Empty(t *testing.T) {
	cfg := config.DefaultConfig()
	eng := New(cfg)
	results := eng.Execute()
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestEngine_SkipList(t *testing.T) {
	cfg := config.DefaultConfig()
	eng := New(cfg)

	passPolicy := &mockPolicy{
		name:   "PassPolicy",
		result: policy.Result{PolicyName: "PassPolicy", Status: policy.StatusPass},
	}
	failPolicy := &mockPolicy{
		name:   "FailPolicy",
		result: policy.Result{PolicyName: "FailPolicy", Status: policy.StatusFail, Message: "should be skipped"},
	}

	eng.Register(passPolicy)
	eng.Register(failPolicy)
	eng.SetSkipList([]string{"FailPolicy"})

	results := eng.Execute()
	if len(results) != 1 {
		t.Fatalf("expected 1 result (skipped), got %d", len(results))
	}
	if results[0].PolicyName != "PassPolicy" {
		t.Errorf("expected PassPolicy result, got %s", results[0].PolicyName)
	}
}

func TestEngine_SkipListMultiple(t *testing.T) {
	cfg := config.DefaultConfig()
	eng := New(cfg)

	p1 := &mockPolicy{name: "PolicyOne", result: policy.Result{PolicyName: "PolicyOne", Status: policy.StatusFail}}
	p2 := &mockPolicy{name: "PolicyTwo", result: policy.Result{PolicyName: "PolicyTwo", Status: policy.StatusFail}}
	p3 := &mockPolicy{name: "PolicyThree", result: policy.Result{PolicyName: "PolicyThree", Status: policy.StatusPass}}

	eng.Register(p1)
	eng.Register(p2)
	eng.Register(p3)
	eng.SetSkipList([]string{"PolicyOne", "PolicyTwo"})

	results := eng.Execute()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].PolicyName != "PolicyThree" {
		t.Errorf("expected PolicyThree, got %s", results[0].PolicyName)
	}
}

func TestEngine_SkipListEmpty(t *testing.T) {
	cfg := config.DefaultConfig()
	eng := New(cfg)

	failPolicy := &mockPolicy{
		name:   "FailPolicy",
		result: policy.Result{PolicyName: "FailPolicy", Status: policy.StatusFail, Message: "fail"},
	}
	eng.Register(failPolicy)
	eng.SetSkipList([]string{})

	results := eng.Execute()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != policy.StatusFail {
		t.Errorf("expected fail, got %s", results[0].Status)
	}
}

func TestEngine_SkipListOverridesDisabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Policies.SetDisabled("FailPolicy", true)
	eng := New(cfg)

	failPolicy := &mockPolicy{
		name:   "FailPolicy",
		result: policy.Result{PolicyName: "FailPolicy", Status: policy.StatusFail, Message: "disabled"},
	}
	eng.Register(failPolicy)
	eng.SetSkipList([]string{"FailPolicy"})

	results := eng.Execute()
	if len(results) != 0 {
		t.Errorf("expected 0 results (disabled takes priority), got %d", len(results))
	}
}

func TestEngine_SkipList_NotSkipping(t *testing.T) {
	cfg := config.DefaultConfig()
	eng := New(cfg)

	passPolicy := &mockPolicy{
		name:   "PassPolicy",
		result: policy.Result{PolicyName: "PassPolicy", Status: policy.StatusPass},
	}
	eng.Register(passPolicy)
	eng.SetSkipList([]string{"OtherPolicy"})

	results := eng.Execute()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestEngine_SkipDisabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Policies.SetDisabled("FailPolicy", true)
	eng := New(cfg)

	passPolicy := &mockPolicy{
		name:   "PassPolicy",
		result: policy.Result{PolicyName: "PassPolicy", Status: policy.StatusPass},
	}
	failPolicy := &mockPolicy{
		name:   "FailPolicy",
		result: policy.Result{PolicyName: "FailPolicy", Status: policy.StatusFail, Message: "should be skipped"},
	}

	eng.Register(passPolicy)
	eng.Register(failPolicy)

	results := eng.Execute()
	if len(results) != 1 {
		t.Fatalf("expected 1 result (skipped disabled), got %d", len(results))
	}
	if results[0].PolicyName != "PassPolicy" {
		t.Errorf("expected PassPolicy result, got %s", results[0].PolicyName)
	}
}

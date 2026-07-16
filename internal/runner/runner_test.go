package runner

import (
	"testing"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/engine"
	"github.com/marcuwynu23/git-policy/internal/policy"
)

func TestRunner_EngineWiring(t *testing.T) {
	eng := engine.New(config.DefaultConfig())

	passPolicy := &mockPolicy{
		name:   "PassPolicy",
		result: policy.Result{PolicyName: "PassPolicy", Status: policy.StatusPass},
	}
	failPolicy := &mockPolicy{
		name:   "FailPolicy",
		result: policy.Result{PolicyName: "FailPolicy", Status: policy.StatusFail, Message: "fail"},
	}

	eng.Register(passPolicy)
	eng.Register(failPolicy)

	results := eng.Execute()
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Status != policy.StatusPass {
		t.Errorf("expected pass, got %s", results[0].Status)
	}
	if results[1].Status != policy.StatusFail {
		t.Errorf("expected fail, got %s", results[1].Status)
	}
}

func TestRunner_DisabledPoliciesAreSkipped(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Policies.SetDisabled("FailPolicy", true)

	eng := engine.New(cfg)

	failPolicy := &mockPolicy{
		name:   "FailPolicy",
		result: policy.Result{PolicyName: "FailPolicy", Status: policy.StatusFail, Message: "should be skipped"},
	}
	passPolicy := &mockPolicy{
		name:   "PassPolicy",
		result: policy.Result{PolicyName: "PassPolicy", Status: policy.StatusPass},
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

type mockPolicy struct {
	name   string
	result policy.Result
}

func (m *mockPolicy) Name() string { return m.name }
func (m *mockPolicy) Validate(ctx policy.Context) policy.Result { return m.result }

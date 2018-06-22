package policy

import (
  "github.com/open-policy-agent/opa/ast"
  "github.com/open-policy-agent/opa/rego"
)

const PolicyTypes = map[string]string{
  DefaultDeny:              "Default denied",
  PolicyEvaluationFailed:   "Error evaluating policy",
  AdminPolicyInvalid:       "The adminstrative policy decision is invalid",
}

type PolicyClient struct {
  Config    Config
}

type PolicyRequest struct {

}

type PolicyResponse struct {
  Allowed     bool
  Reason      string
}

func (pc *PolicyClient) evaluate(ctx context.Context, req PolicyRequest) (*PolicyResponse, error) {
  resp := &PolicyResponse{
    Allowed: false,
    Reason: PolicyTypes[DefaultDeny],
  }

  bs, err := ioutil.ReadFile(pc.Config.PolicyFile)
  if err != nil {
    log.Warn("Error reading policy file '%s': %v", pc.Config.PolicyFile, err)
    return resp, err
  }

  allowed, err := func() (bool, error) {
    eval := rego.New(rego.Query(p.allowPath), req.Input(input), rego.Module(pc.Config.PolicyFile, string(bs)))

    rs, err := eval.Eval(ctx)
    if err != nil {
      resp.Reason = PolicyTypes[PolicyEvaluationFailed]
      return resp, err
    }

    if len(rs) == 0 {
      resp.Reason = PolicyTypes[DefaultDeny]
      return resp, err
    }

    allowed, ok := rs[0].Expressions[0].Value.(bool)
    if !ok {
      resp.Reason = PolicyTypes[AdminPolicyInvalid]
      return resp, err
    }

    resp.Allowed = allowed
    return resp, nil
  }()

  if err != nil {
    log.Printf("Returning OPA policy decision: %v (error: %v)", allowed, err)
  } else {
    log.Printf("Returning OPA policy decision: %v", allowed)
  }

  return resp, err
}

func NewPolicyClient(cfg *config.Config) *PolicyClient {
  return &PolicyClient{
    Config: cfg.Policy,
  }
}

package trigger

import "github.com/qri-io/qri/profile"

// Source is an abstraction for a `workflow.Workflow`
type Source interface {
	WorkflowID() string
	ActiveTriggers() []Trigger
	ScopeID() profile.ID
}

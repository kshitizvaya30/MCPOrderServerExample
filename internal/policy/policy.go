package policy

import "fmt"

type Action string

const (
	ActionGetOrder     Action = "get_order"
	ActionCreateReturn Action = "create_return"
	ActionCreateRefund Action = "create_refund"
)

type Decision struct {
	Allowed       bool
	RequiresHuman bool
	Reason        string
}

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Check(
	action Action,
	orderAmount float64,
) Decision {

	switch action {

	case ActionGetOrder:
		return Decision{
			Allowed: true,
		}

	case ActionCreateReturn:
		return Decision{
			Allowed: true,
		}

	case ActionCreateRefund:
		if orderAmount > 5000 {
			return Decision{
				Allowed:       false,
				RequiresHuman: true,
				Reason: fmt.Sprintf(
					"refund amount %.2f exceeds the automatic refund limit of 5000",
					orderAmount,
				),
			}
		}

		return Decision{
			Allowed: true,
		}

	default:
		return Decision{
			Allowed: false,
			Reason:  "action is not supported",
		}
	}
}

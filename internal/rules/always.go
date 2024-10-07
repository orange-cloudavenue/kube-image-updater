package rules

var _ RuleInterface = &always{}

type (
	always struct {
		rule
	}
)

func init() {
	register(Always, &always{})
}

// ! always rule

func (a *always) Evaluate() (matchWithRule bool, newTag string, err error) {
	return true, a.tags[0], nil
}

package kubeclient

import (
	"fmt"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
)

type (
	matchConditionBuilderInterface interface {
		buildMatchCondition() []admissionregistrationv1.MatchCondition
		getName() string
	}

	namespaceMatchConditionBuilder struct {
		namespace string
	}

	defaultMatchConditionBuilder struct{}
)

// defaultMatchConditionBuilder

var _ matchConditionBuilderInterface = &defaultMatchConditionBuilder{}

func (m *defaultMatchConditionBuilder) buildMatchCondition() []admissionregistrationv1.MatchCondition {
	return []admissionregistrationv1.MatchCondition{
		{
			Name:       "annotation-is-true",
			Expression: fmt.Sprintf("object.metadata.?annotations['%s'].orValue('false') == 'true'", annotations.KeyEnabled),
		},
	}
}

func (m *defaultMatchConditionBuilder) getName() string {
	return "default"
}

// * namespaceMatchConditionBuilder

var _ matchConditionBuilderInterface = &namespaceMatchConditionBuilder{}

func (n *namespaceMatchConditionBuilder) buildMatchCondition() []admissionregistrationv1.MatchCondition {
	return []admissionregistrationv1.MatchCondition{
		{
			Name:       "annotation-is-not-false",
			Expression: fmt.Sprintf("object.metadata.?annotations['%s'].orValue('') != 'false'", annotations.KeyEnabled),
		},
		{
			Name:       fmt.Sprintf("namespace-%s-match", n.namespace),
			Expression: fmt.Sprintf("object.metadata.namespace == '%s'", n.namespace),
		},
	}
}

func (n *namespaceMatchConditionBuilder) getName() string {
	return n.namespace + ".ns"
}

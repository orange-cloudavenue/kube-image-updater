package kubeclient

import (
	"fmt"
	"strings"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
)

type (
	matchConditionBuilderInterface interface {
		buildMatchCondition() []admissionregistrationv1.MatchCondition
		GetName() string
	}

	NamespaceMatchConditionBuilder struct {
		namespaceMatchConditionBuilder
	}
	namespaceMatchConditionBuilder struct {
		Namespace string
	}

	defaultMatchConditionBuilder struct {
		FailurePolicy admissionregistrationv1.FailurePolicyType
	}
)

func (n NamespaceMatchConditionBuilder) New(namespace string) matchConditionBuilderInterface {
	return &namespaceMatchConditionBuilder{
		Namespace: namespace,
	}
}

// * defaultMatchConditionBuilder

var _ matchConditionBuilderInterface = &defaultMatchConditionBuilder{}

func (m defaultMatchConditionBuilder) buildMatchCondition() []admissionregistrationv1.MatchCondition {
	policy := ""

	switch m.FailurePolicy {
	case admissionregistrationv1.Fail:
		policy = string(annotations.FailurePolicyFail)
	case admissionregistrationv1.Ignore:
		policy = string(annotations.FailurePolicyIgnore)
	}

	// In the expression failure-policy the default value is fail because the default value of the mutation is fail.
	return []admissionregistrationv1.MatchCondition{
		{
			Name:       "annotation-is-true",
			Expression: fmt.Sprintf("object.metadata.?annotations['%s'].orValue('false') == 'true'", annotations.KeyEnabled),
		},
		{
			Name:       "failure-policy",
			Expression: fmt.Sprintf("object.metadata.?annotations['%s'].orValue('fail') == '%s'", annotations.KeyFailurePolicy, policy),
		},
	}
}

func (m defaultMatchConditionBuilder) GetName() string {
	return fmt.Sprintf("default.%s.%s", strings.ToLower(string(m.FailurePolicy)), models.MutatorWebhookName)
}

// * namespaceMatchConditionBuilder

var _ matchConditionBuilderInterface = &namespaceMatchConditionBuilder{}

func (n namespaceMatchConditionBuilder) buildMatchCondition() []admissionregistrationv1.MatchCondition {
	return []admissionregistrationv1.MatchCondition{
		{
			Name:       "annotation-is-not-false",
			Expression: fmt.Sprintf("object.metadata.?annotations['%s'].orValue('') != 'false'", annotations.KeyEnabled),
		},
		{
			Name:       fmt.Sprintf("namespace-%s-match", n.Namespace),
			Expression: fmt.Sprintf("object.metadata.namespace == '%s'", n.Namespace),
		},
	}
}

func (n namespaceMatchConditionBuilder) GetName() string {
	return n.Namespace + ".ns." + models.MutatorWebhookName
}

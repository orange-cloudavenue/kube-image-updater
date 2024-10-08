package kubeclient

import (
	"context"
	"encoding/base64"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

func (c *Client) GetValueOrValueFrom(ctx context.Context, namespace string, v v1alpha1.ValueOrValueFrom) (any, error) {
	// If ValueFrom is nil, return the value
	if v.ValueFrom == nil {
		return v.Value, nil
	}

	// Read from configmap
	if v.ValueFrom.ConfigMapKeyRef != nil {
		cm, err := c.CoreV1().ConfigMaps(namespace).Get(ctx, v.ValueFrom.ConfigMapKeyRef.Name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		if value, found := cm.Data[v.ValueFrom.ConfigMapKeyRef.Key]; found {
			return value, nil
		}

		return "", fmt.Errorf("key %s not found in configmap %s", v.ValueFrom.ConfigMapKeyRef.Key, v.ValueFrom.ConfigMapKeyRef.Name)
	}

	// Read from secret
	if v.ValueFrom.SecretKeyRef != nil {
		secret, err := c.CoreV1().Secrets(namespace).Get(ctx, v.ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		if valueBase64, found := secret.Data[v.ValueFrom.SecretKeyRef.Key]; found {
			// decode base64
			v, err := base64.StdEncoding.DecodeString(string(valueBase64))
			if err != nil {
				return "", err
			}

			return string(v), nil
		}

		if valueString, found := secret.StringData[v.ValueFrom.SecretKeyRef.Key]; found {
			return valueString, nil
		}

		return "", fmt.Errorf("key %s not found in secret %s", v.ValueFrom.SecretKeyRef.Key, v.ValueFrom.SecretKeyRef.Name)
	}

	if v.ValueFrom.AlertConfigRef != nil {
		alert, err := c.Alert().Get(ctx, namespace, v.ValueFrom.AlertConfigRef.Name)
		if err != nil {
			return "", fmt.Errorf("error getting alert config %s: %w", v.ValueFrom.AlertConfigRef.Name, err)
		}

		return alert, nil
	}

	return "", fmt.Errorf("ValueFrom is specified but neither SecretKeyRef nor ConfigMapKeyRef nor AlertConfigRef is set")
}

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type (
	ValueOrValueFrom struct {
		// Value is a string value to assign to the key.
		// if ValueFrom is specified, this value is ignored.
		// +optional
		Value string `json:"value,omitempty"`

		// ValueFrom is a reference to a field in a secret or config map.
		// +optional
		ValueFrom *ValueFromSource `json:"valueFrom,omitempty"`
	}

	// ValueFromSource is a reference to a field in a secret or config map.
	ValueFromSource struct {
		// SecretKeyRef is a reference to a field in a secret.
		// +optional
		SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`

		// ConfigMapKeyRef is a reference to a field in a config map.
		// +optional
		ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	}
)

// func (v *ValueOrValueFrom) GetValue(ctx context.Context, namespace, name string) (string, error) {

// 	// If ValueFrom is nil, return the value
// 	if v.ValueFrom == nil {
// 		return v.Value, nil
// 	}

// 	if v.ValueFrom.SecretKeyRef == nil && v.ValueFrom.ConfigMapKeyRef == nil {
// 		return "", fmt.Errorf("ValueFrom is specified but neither SecretKeyRef nor ConfigMapKeyRef is set")
// 	}

// 	// // Read from config map
// 	// if p.ConfigMapKeyRef != nil {
// 	// 	var configMap coreV1.ConfigMap
// 	// 	objectKey := types.NamespacedName{Namespace: namespace, Name: p.ConfigMapKeyRef.Name}
// 	// 	err := r.GetKubeClient().Get(ctx, objectKey, &configMap)
// 	// 	if errors.IsNotFound(err) {
// 	// 		return "", err
// 	// 	}
// 	// 	key := p.ConfigMapKeyRef.Key
// 	// 	if value, found := configMap.Data[key]; found {
// 	// 		return value, nil
// 	// 	}
// 	// 	return "", brose_errors.NewMapEntryNotFoundError(key, nil)
// 	// }

// 	// // Read from secret
// 	// if p.SecretKeyRef != nil {
// 	// 	var secret coreV1.Secret
// 	// 	objectKey := types.NamespacedName{Namespace: namespace, Name: p.SecretKeyRef.Name}
// 	// 	err := r.Get(ctx, objectKey, &secret)
// 	// 	if errors.IsNotFound(err) {
// 	// 		return "", err
// 	// 	}
// 	// 	key := p.SecretKeyRef.Key
// 	// 	valueBase64, foundBase64 := secret.Data[key]
// 	// 	valueString, foundString := secret.StringData[key]
// 	// 	if !foundBase64 && !foundString {
// 	// 		return "", brose_errors.NewMapEntryNotFoundError(key, nil)
// 	// 	} else if foundString {
// 	// 		return valueString, nil
// 	// 	}
// 	// 	return string(valueBase64), nil
// 	// }
// 	// return "", brose_errors.NewMissingPropertyValueError(name, nil)

// 	return "", nil
// }

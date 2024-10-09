package actions_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/actions"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
	"github.com/orange-cloudavenue/kube-image-updater/test/mocks/fakekubeclient"
)

func TestAlertEmail_Execute(t *testing.T) {
	namespace := "default"
	ctx := context.TODO()

	image := v1alpha1.Image{
		TypeMeta: v1.TypeMeta{
			Kind:       "Image",
			APIVersion: fmt.Sprintf("%s/%s", v1alpha1.GroupVersion.Group, v1alpha1.GroupVersion.Version),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "demo",
			Namespace: namespace,
		},
		Spec: v1alpha1.ImageSpec{
			BaseTag: "v0.0.1",
			Rules: []v1alpha1.ImageRule{
				{
					Name: "Always update",
					Type: rules.Always,
					Actions: []v1alpha1.ImageAction{
						{
							Type: actions.AlertEmail.String(),
							Data: v1alpha1.ValueOrValueFrom{
								ValueFrom: &v1alpha1.ValueFromSource{
									AlertConfigRef: &corev1.LocalObjectReference{
										Name: "demo",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		emailSpec     v1alpha1.AlertConfig
		expectedError bool
	}{
		{
			name: "Valid configuration without auth",
			emailSpec: v1alpha1.AlertConfig{
				ObjectMeta: v1.ObjectMeta{
					Name:      "demo",
					Namespace: namespace,
				},
				Spec: v1alpha1.AlertConfigSpec{
					Email: &v1alpha1.AlertEmailSpec{
						Username:        v1alpha1.ValueOrValueFrom{Value: "", ValueFrom: &v1alpha1.ValueFromSource{}},
						Password:        v1alpha1.ValueOrValueFrom{Value: "", ValueFrom: &v1alpha1.ValueFromSource{}},
						Host:            v1alpha1.ValueOrValueFrom{Value: "smtp.freesmtpservers.com", ValueFrom: &v1alpha1.ValueFromSource{}},
						Port:            v1alpha1.ValueOrValueFrom{Value: "25", ValueFrom: &v1alpha1.ValueFromSource{}},
						FromAddress:     "from@example.com",
						ToAddress:       []string{"to@example.com"},
						ClientHost:      "client.example.com",
						Encryption:      "None",
						FromName:        "Example",
						Auth:            "none",
						UseHTML:         false,
						UseStartTLS:     true,
						TemplateSubject: "Test Subject",
					},
				},
			},
			expectedError: false,
		},
		{
			name: "invalid host",
			emailSpec: v1alpha1.AlertConfig{
				Spec: v1alpha1.AlertConfigSpec{
					Email: &v1alpha1.AlertEmailSpec{
						Username:        v1alpha1.ValueOrValueFrom{Value: "", ValueFrom: &v1alpha1.ValueFromSource{}},
						Password:        v1alpha1.ValueOrValueFrom{Value: "", ValueFrom: &v1alpha1.ValueFromSource{}},
						Host:            v1alpha1.ValueOrValueFrom{Value: "1.1.1.1", ValueFrom: &v1alpha1.ValueFromSource{}},
						Port:            v1alpha1.ValueOrValueFrom{Value: "8080", ValueFrom: &v1alpha1.ValueFromSource{}},
						FromAddress:     "from@example.com",
						ToAddress:       []string{"to@example.com"},
						ClientHost:      "client.example.com",
						Encryption:      "TLS",
						FromName:        "Example",
						Auth:            "none",
						UseHTML:         false,
						UseStartTLS:     true,
						TemplateSubject: "Test Subject",
					},
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fakekubeclient.NewFakeKubeClient()
			assert.NoError(t, fakeClient.CreateFakeImage(image))

			a, err := actions.GetAction(actions.AlertEmail)
			assert.NoError(t, err)

			a.Init(fakeClient, models.Tags{
				Actual: "1.0.0",
				New:    "1.1.0",
			}, &image, image.Spec.Rules[0].Actions[0].Data)

			fakeClient.On("GetValueOrValueFrom", ctx, image.GetNamespace(), image.Spec.Rules[0].Actions[0].Data).Return(tt.emailSpec, nil)
			fakeClient.On("GetValueOrValueFrom", ctx, tt.emailSpec.GetNamespace(), tt.emailSpec.Spec.Email.Username).Return("", fmt.Errorf("error"))
			fakeClient.On("GetValueOrValueFrom", ctx, tt.emailSpec.GetNamespace(), tt.emailSpec.Spec.Email.Host).Return(tt.emailSpec.Spec.Email.Host.Value, nil)
			fakeClient.On("GetValueOrValueFrom", ctx, tt.emailSpec.GetNamespace(), tt.emailSpec.Spec.Email.Port).Return(tt.emailSpec.Spec.Email.Port.Value, nil)

			err = a.Execute(ctx)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

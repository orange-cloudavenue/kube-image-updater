package actions_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

// MockKubeClient is a mock implementation of the KubeClient interface.
type MockKubeClient struct {
	mock.Mock
}

func (m *MockKubeClient) GetValueOrValueFrom(ctx context.Context, namespace string, value v1alpha1.ValueOrValueFrom) (any, error) {
	args := m.Called(ctx, namespace, value)
	return args.Get(2), args.Error(1)
}

// func TestAlertEmail_ConstructURL(t *testing.T) {
// 	mockKubeClient := new(MockKubeClient)
// 	ctx := context.TODO()
// 	namespace := "default"

// 	tests := []struct {
// 		name          string
// 		emailSpec     v1alpha1.AlertEmailSpec
// 		expectedURL   string
// 		expectedError bool
// 	}{
// 		{
// 			name: "Valid configuration without auth",
// 			emailSpec: v1alpha1.AlertEmailSpec{
// 				Host:            v1alpha1.ValueOrValueFrom{Value: "smtp.freesmtpservers.com"},
// 				Port:            v1alpha1.ValueOrValueFrom{Value: "587"},
// 				FromAddress:     "from@example.com",
// 				ToAddress:       []string{"to@example.com"},
// 				ClientHost:      "client.example.com",
// 				Encryption:      "TLS",
// 				FromName:        "Example",
// 				Auth:            "none",
// 				UseHTML:         false,
// 				UseStartTLS:     true,
// 				TemplateSubject: "Test Subject",
// 			},
// 			expectedURL:   "smtp://smtp.freesmtpservers.com:587?from=from@example.com&to=to@example.com&clientHost=client.example.com&encryption=TLS&fromName=Example&auth=none&useHTML=no&starttls=yes&subject=Test Subject",
// 			expectedError: false,
// 		},
// 		{
// 			name: "invalid host",
// 			emailSpec: v1alpha1.AlertEmailSpec{
// 				Host:            v1alpha1.ValueOrValueFrom{Value: ""},
// 				Port:            v1alpha1.ValueOrValueFrom{Value: "587"},
// 				FromAddress:     "from@example.com",
// 				ToAddress:       []string{"to@example.com"},
// 				ClientHost:      "client.example.com",
// 				Encryption:      "TLS",
// 				FromName:        "Example",
// 				Auth:            "none",
// 				UseHTML:         false,
// 				UseStartTLS:     true,
// 				TemplateSubject: "Test Subject",
// 			},
// 			expectedURL:   "",
// 			expectedError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockKubeClient.On("GetValueOrValueFrom", ctx, namespace, tt.emailSpec.Host).Return(tt.emailSpec.Host.Value, nil)
// 			mockKubeClient.On("GetValueOrValueFrom", ctx, namespace, tt.emailSpec.Port).Return(tt.emailSpec.Port.Value, nil)

// 			a, err := actions.GetAction(actions.AlertEmail)
// 			assert.NoError(t, err)

// 			image := &v1alpha1.Image{}
// 			alertEmail := &v1alpha1.AlertConfig{}

// 			a.Init(mockKubeClient, models.Tags{
// 				Actual: "1.0.0",
// 				New:    "1.1.0",
// 			}, image, v1alpha1.ValueOrValueFrom{})

// 			err = a.Execute(ctx)
// 			if tt.expectedError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

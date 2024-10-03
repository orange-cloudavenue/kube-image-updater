package v1alpha1

type (
	AlertCoreSpec struct{}

	AlertCoreStatus struct {
		Synced bool `json:"synced"`
	}
)

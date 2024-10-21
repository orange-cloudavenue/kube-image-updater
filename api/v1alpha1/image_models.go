package v1alpha1

type (
	ImageStatusLastSync string
)

var (
	// Status of the image when it is last sync success.
	ImageStatusLastSyncSuccess ImageStatusLastSync = "Success"

	// Status of the image when it is last sync is scheduled.
	ImageStatusLastSyncScheduled ImageStatusLastSync = "Scheduled"

	// Status of the image when an error occurred.
	ImageStatusLastSyncError ImageStatusLastSync = "Error"

	// Status of the image when it is last sync get error.
	ImageStatusLastSyncErrorGetImage ImageStatusLastSync = "GetImageError"

	// Status of the image when it is last sync error secrets.
	ImageStatusLastSyncErrorPullSecrets ImageStatusLastSync = "PullSecretsError"

	// Status of the image when it is last sync error registry.
	ImageStatusLastSyncErrorRegistry ImageStatusLastSync = "RegistryError"

	// Status of the image when it is last sync error tags.
	ImageStatusLastSyncErrorTags ImageStatusLastSync = "TagsError"

	// Status of the image when it is last sync error get rule.
	ImageStatusLastSyncErrorGetRule ImageStatusLastSync = "GetRuleError"

	// Status of the image when it is last sync error action.
	ImageStatusLastSyncErrorAction ImageStatusLastSync = "ActionError"
)

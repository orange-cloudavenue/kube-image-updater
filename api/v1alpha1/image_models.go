package v1alpha1

type (
	ImageStatusLastSync string
)

var (
	// ImageStatusLastSyncSuccess is the status of the image when it is last sync success
	ImageStatusLastSyncSuccess ImageStatusLastSync = "Success"

	// ImageStatusLastSyncScheduled is the status of the image when it is last sync is scheduled
	ImageStatusLastSyncScheduled ImageStatusLastSync = "Scheduled"

	// ERROR STATUS
	// ImageStatusError is the status of the image when an error occurred
	ImageStatusLastSyncError ImageStatusLastSync = "Error"
	// ImageStatusLastSyncGetError is the status of the image when it is last sync get error
	ImageStatusLastSyncErrorGetImage ImageStatusLastSync = "GetImageError"
	// ImageStatusLastSyncErrorSecrets is the status of the image when it is last sync error secrets
	ImageStatusLastSyncErrorPullSecrets ImageStatusLastSync = "PullSecretsError"
	// ImageStatusLastSyncErrorRegistry is the status of the image when it is last sync error registry
	ImageStatusLastSyncErrorRegistry ImageStatusLastSync = "RegistryError"
	// ImageStatusLastSyncErrorTags is the status of the image when it is last sync error tags
	ImageStatusLastSyncErrorTags ImageStatusLastSync = "TagsError"
	// ImageStatusLastSyncErrorGetRule is the status of the image when it is last sync error get rule
	ImageStatusLastSyncErrorGetRule ImageStatusLastSync = "GetRuleError"
	// ImageStatusLastSyncErrorAction is the status of the image when it is last sync error action
	ImageStatusLastSyncErrorAction ImageStatusLastSync = "ActionError"
)

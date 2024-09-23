package annotations

type (
	Annotation interface{}

	AnnotationAction struct {
		Annotation
	}

	AnnotationTag struct {
		Annotation
		Tag string
	}

	AActionKey string
)

const (
	// Action Refresh
	ActionRefresh AActionKey = "refresh"
)

// AnnotationKey is the key used to store the image in the annotation
var (
	AnnotationsAction = map[AActionKey]AnnotationAction{
		// kimup.cloudavenue.io/action: refresh
		ActionRefresh: {},
	}

	// AnnotationTag = AnnotationTag{}

	AnnotationActionKey   = "kimup.cloudavenue.io" + "/action"
	AnnotationTagKey      = "kimup.cloudavenue.io" + "/tag"
	AnnotationCheckSumKey = "kimup.cloudavenue.io" + "/checksum"
)

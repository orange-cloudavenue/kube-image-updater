package patch

import "encoding/json"

type (
	Op string

	Builder struct {
		patches []Patch
	}

	Patch struct {
		Op    Op          `json:"op"`
		Path  string      `json:"path"`
		Value interface{} `json:"value,omitempty"`
	}
)

const (
	OpAdd     Op = "add"
	OpRemove  Op = "remove"
	OpReplace Op = "replace"
)

// NewBuilder creates a new Builder
func NewBuilder() *Builder {
	return &Builder{
		patches: make([]Patch, 0),
	}
}

// AddPatch adds a patch to the Builder
func (pb *Builder) AddPatch(op Op, path string, value interface{}) {
	pb.patches = append(pb.patches, Patch{
		Op:    op,
		Path:  path,
		Value: value,
	})
}

// AddRawPatch adds a raw patch to the Builder
func (pb *Builder) AddRawPatch(patch Patch) {
	pb.patches = append(pb.patches, patch)
}

// AddRawPatches adds raw patches to the Builder
func (pb *Builder) AddRawPatches(patches []Patch) {
	pb.patches = append(pb.patches, patches...)
}

// Generate generates the patch
func (pb *Builder) Generate() ([]byte, error) {
	return json.Marshal(pb.patches)
}

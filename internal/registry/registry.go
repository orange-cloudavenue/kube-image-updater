package registry

import (
	"context"
	"errors"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/ref"
)

var (
	ErrRepoIsEmpty = errors.New("repo is empty")
	ErrInvalidRepo = errors.New("invalid repo")
)

type (
	Repository struct {
		ctx  context.Context
		repo string
		c    *regclient.RegClient
	}
)

func New(ctx context.Context, repo string) (*Repository, error) {
	if repo == "" {
		return nil, ErrRepoIsEmpty
	}

	// define a regclient with desired options
	rc := regclient.New(
		regclient.WithDockerCerts(),
		regclient.WithDockerCreds(),
		regclient.WithUserAgent("kimup"),
	)

	r := &Repository{
		ctx:  ctx,
		repo: repo,
		c:    rc,
	}

	ref, err := r.buildRef()
	if err != nil {
		return nil, err
	}

	if _, err := rc.Ping(ctx, ref); err != nil {
		return nil, err
	}

	return r, nil
}

// buildRef creates a reference for an image
func (r *Repository) buildRef() (ref.Ref, error) {
	// create a reference for an image
	rF, err := ref.New(r.GetRepo())
	if err != nil {
		return ref.Ref{}, err
	}

	return rF, nil
}

func (r *Repository) Tags() ([]string, error) {
	// create a reference for an image
	rF, err := ref.New(r.GetRepo())
	if err != nil {
		return nil, err
	}
	defer r.c.Close(r.ctx, rF)

	tL, err := r.c.TagList(r.ctx, rF)
	if err != nil {
		return nil, err
	}

	return tL.GetTags()
}

// GetRepo returns the repository name
func (r *Repository) GetRepo() string {
	return r.repo
}

package registry

import (
	"context"
	"errors"

	"github.com/containers/image/v5/types"
	dRegistry "github.com/crazy-max/diun/v4/pkg/registry"
)

var (
	ErrRepoIsEmpty = errors.New("repo is empty")
	ErrInvalidRepo = errors.New("invalid repo")
)

type (
	Repository struct {
		ctx  context.Context
		repo string
		r    *dRegistry.Client
		dR   dRegistry.Image
	}

	Settings struct {
		InsecureTLS bool
		Username    string
		Password    string
	}
)

func New(ctx context.Context, repo string, settings Settings) (*Repository, error) {
	if repo == "" {
		return nil, ErrRepoIsEmpty
	}

	dR, err := dRegistry.ParseImage(dRegistry.ParseImageOptions{
		Name: repo,
	})
	if err != nil {
		return nil, ErrInvalidRepo
	}

	r, err := dRegistry.New(dRegistry.Options{
		Auth: types.DockerAuthConfig{
			Username: settings.Username,
			Password: settings.Password,
		},
		InsecureTLS: settings.InsecureTLS,
		UserAgent:   "kube-image-updater",
	})
	if err != nil {
		return nil, err
	}

	rr := &Repository{
		ctx:  ctx,
		repo: repo,
		r:    r,
		dR:   dR,
	}

	return rr, nil
}

func (r *Repository) Tags() ([]string, error) {
	tags, err := r.r.Tags(dRegistry.TagsOptions{
		Image: r.dR,
	})
	if err != nil {
		return nil, err
	}

	return tags.List, nil
}

// GetRepo returns the repository name
func (r *Repository) GetRepo() string {
	return r.repo
}

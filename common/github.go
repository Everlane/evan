package common

import (
	"fmt"
	"net/url"

	"github.com/google/go-github/github"
)

type GithubRepository struct {
	Repository   Repository
	GithubClient *github.Client
}

func NewGithubRepository(repository Repository, githubClient *github.Client) *GithubRepository {
	return &GithubRepository{
		Repository:   repository,
		GithubClient: githubClient,
	}
}

func NewGithubRepositoryFromDeployment(deployment Deployment) (*GithubRepository, error) {
	repository := deployment.Application().Repository()
	githubClient, err := deployment.GithubClient()
	if err != nil {
		return nil, err
	} else {
		return NewGithubRepository(repository, githubClient), nil
	}
}

func (repo *GithubRepository) OwnerAndName() (string, string) {
	return repo.Repository.Owner(), repo.Repository.Name()
}

func (repo *GithubRepository) Get() (*github.Repository, error) {
	owner, name := repo.OwnerAndName()
	repository, _, err := repo.GithubClient.Repositories.Get(owner, name)
	return repository, err
}

func (repo *GithubRepository) CompareCommits(base, head string) (*github.CommitsComparison, error) {
	owner, name := repo.OwnerAndName()
	commitsComparison, _, err := repo.GithubClient.Repositories.CompareCommits(owner, name, base, head)
	return commitsComparison, err
}

// The `go-github` library doesn't export a visible archive format, so we have
// to use our own.
type ArchiveFormat string

const (
	Tarball ArchiveFormat = "tarball"
	Zipball ArchiveFormat = "zipball"
)

// `format` should be one of "tarball" or "zipball".
func (repo *GithubRepository) GetArchiveLink(format ArchiveFormat, ref string) (string, error) {
	owner, name := repo.OwnerAndName()

	var url *url.URL
	var err error

	options := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// Work-around for `go-github` not exporting their `archiveFormat` type.
	switch format {
	case Tarball:
		url, _, err = repo.GithubClient.Repositories.GetArchiveLink(owner, name, github.Tarball, options)
	case Zipball:
		url, _, err = repo.GithubClient.Repositories.GetArchiveLink(owner, name, github.Zipball, options)
	default:
		return "", fmt.Errorf("Unknown archive format: '%v'", format)
	}

	if err != nil {
		return "", err
	} else {
		return url.String(), nil
	}
}

func (repo *GithubRepository) GetCommitSHA1(ref string) (string, error) {
	owner, name := repo.OwnerAndName()
	sha1, _, err := repo.GithubClient.Repositories.GetCommitSHA1(owner, name, ref, "")
	return sha1, err
}

func (repo *GithubRepository) Merge(base, head, commitMessage string) (*github.RepositoryCommit, error) {
	owner, name := repo.OwnerAndName()
	request := &github.RepositoryMergeRequest{
		Base:          &base,
		Head:          &head,
		CommitMessage: &commitMessage,
	}
	commit, _, err := repo.GithubClient.Repositories.Merge(owner, name, request)
	return commit, err
}

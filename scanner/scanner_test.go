package scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type fakeRepo struct{}

func (f *fakeRepo) CommitObjects() (object.CommitIter, error) {
	return &fakeIter{}, nil
}

type fakeIter struct{}

func (f *fakeIter) Next() (*object.Commit, error) {
	return &object.Commit{}, nil
}

func (f *fakeIter) ForEach(func(*object.Commit) error) error {
	return nil
}
func (f *fakeIter) Close() {}

func TestScan(t *testing.T) {
	fake := &fakeRepo{}
	secret := Scan(fake)
	assert.Equal(t, "", secret)
}

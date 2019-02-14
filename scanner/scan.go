// Copyright © 2019 Damian Myerscough
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package scan

import (
	"log"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Secret struct {
	Committer object.Signature
	Filename  string
	Secret    string
}

type secret map[string][]Secret

var (
	secret_regex = regexp.MustCompile(`(?:AWS)?_?(?:SECRET|ACCOUNT)?_?(?:ACCESS|ID)?_?(?:KEY)?(?:\s*)?(?::|=>|=)?(?:\s*)?(?:\"|\')?(?P<secret>AKIA[0-9A-Z]{16})(?:\"|\')?`)
)

func DownloadRepository(repoUrl string, dir string) *git.Repository {

	repo, err := git.PlainClone(dir, true, &git.CloneOptions{
		URL:               repoUrl,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		log.Fatal("Unable to download repository: ", err)
	}

	return repo
}

type CommitGetter interface {
	CommitObjects() (object.CommitIter, error)
}

func Scan(repo CommitGetter) secret {
	secrets := secret{}

	commits, _ := repo.CommitObjects()
	defer commits.Close()

	for {
		commit, err := commits.Next()

		if err != nil {
			break
		}

		files, err := commit.Files()

		files.ForEach(func(f *object.File) error {
			content, _ := f.Contents()

			for _, i := range strings.Split(content, "\n") {
				secret := secret_regex.FindString(i)
				if secret != "" {
					secrets[commit.ID().String()] = append(secrets[commit.ID().String()], Secret{commit.Committer, f.Name, secret})
				}
			}
			return nil
		})
	}

	return secrets
}

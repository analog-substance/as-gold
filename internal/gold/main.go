package gold

import (
	"context"

	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/analog-substance/gold/internal/util"
	"github.com/analog-substance/gold/pkg/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v45/github"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type SolidGold struct {
	*types.Group
}

func NewSolidGold() *SolidGold {
	return &SolidGold{
		types.NewGroup(),
	}
}

func FromJSONFile(filepath string) *SolidGold {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		s := NewSolidGold()
		err := s.ToJSONFile(filepath)
		util.CheckErr(err)
	}

	jsonFile, err := os.Open(filepath)
	util.CheckErr(err)
	defer jsonFile.Close()
	byteAr, err := io.ReadAll(jsonFile)
	util.CheckErr(err)
	return FromJSON(byteAr)
}

func FromJSON(byteAr []byte) *SolidGold {
	var solidGold SolidGold
	err := json.Unmarshal(byteAr, &solidGold)
	util.CheckErr(err)
	return &solidGold

}

func (s *SolidGold) ToJSON() []byte {
	goldJSON, err := json.Marshal(s)
	util.CheckErr(err)
	return goldJSON
}

func (s *SolidGold) ToJSONFile(goldFile string) error {
	return os.WriteFile(goldFile, s.ToJSON(), 0644)
}

func (s *SolidGold) ProcessPath(path string) *types.Group {

	gitRepos, err := findGitDirs(path)
	util.CheckErr(err)

	if _, err := os.Stat(fmt.Sprintf("%s/.git", path)); !os.IsNotExist(err) {
		gitRepos = append(gitRepos, path)
	}

	for _, repoPath := range gitRepos {
		s.ProcessRepo(err, repoPath)
	}

	s.Group.MergeDuplicate()
	return s.Group
}

func (s *SolidGold) ProcessRepo(err error, path string) {
	r, err := git.PlainOpen(path)
	util.CheckErr(err)

	cIter, err := r.Log(&git.LogOptions{All: true})
	util.CheckErr(err)

	err = cIter.ForEach(func(c *object.Commit) error {
		s.Group.FindOrCreateHuman(c.Author.Name, c.Author.Email)
		s.Group.FindOrCreateHuman(c.Committer.Name, c.Committer.Email)
		return nil
	})
}

func (s *SolidGold) Merge(primaryID string, otherIDs ...string) {
	s.Group.MergeIDs(primaryID, otherIDs...)
}

func (s *SolidGold) ConsumeBreachFiles(breachFilePaths ...string) {
	for _, breachFilePath := range breachFilePaths {
		tsvFile, err := os.Open(breachFilePath)
		util.CheckErr(err)
		reader := csv.NewReader(tsvFile)
		reader.Comma = '\t'
		reader.FieldsPerRecord = -1

		tsvData, err := reader.ReadAll()
		util.CheckErr(err)

		tsvFile.Close()

		for _, cols := range tsvData {
			email := cols[0]
			password := cols[1]

			human := s.Group.FindOrCreateHumanByEmail(email)
			human.AddPassword(password)
		}
	}
}

func findGitDirs(dir string) ([]string, error) {

	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".git") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func (s *SolidGold) ConsumeGithubOrgs(orgs ...string) {
	client := github.NewClient(nil)

	for _, org := range orgs {
		opt := &github.RepositoryListByOrgOptions{Type: "sources"}
		repos, _, _ := client.Repositories.ListByOrg(context.Background(), org, opt)

		for _, repo := range repos {
			gitCloneURL(fmt.Sprintf("github.com/%s", repo.GetFullName()), repo.GetCloneURL())
		}
		members, _, _ := client.Organizations.ListMembers(context.Background(), org, nil)

		for _, member := range members {
			if member.GetEmail() != "" {
				h := s.FindOrCreateHuman(member.GetLogin(), member.GetEmail())

				if member.GetName() != "" {
					h.AddName(member.GetName())
				}
			}
			s.ConsumeGithubUsers(member.GetLogin())
		}
	}
}

func (s *SolidGold) ConsumeGithubUsers(users ...string) {
	client := github.NewClient(nil)

	for _, user := range users {
		opt := &github.RepositoryListOptions{Type: "sources"}
		repos, _, _ := client.Repositories.List(context.Background(), user, opt)

		for _, repo := range repos {
			gitCloneURL(fmt.Sprintf("github.com/%s", repo.GetFullName()), repo.GetCloneURL())
		}
	}
}

func gitCloneURL(path, repoURL string) {

	git.PlainClone(path, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
}

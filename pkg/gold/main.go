package gold

import (
	"context"
	"github.com/analog-substance/as-gold/pkg/util"
	"github.com/go-git/go-git/v5"
	"log"
	"path"

	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/analog-substance/as-gold/pkg/types"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v56/github"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const FolderPath = "github.com"

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
			email := ""
			password := ""
			if len(cols) >= 1 {
				email = cols[0]
			}

			if len(cols) >= 2 {
				password = cols[1]
			}

			human := s.Group.FindOrCreateHumanByEmail(email)
			human.AddPassword(password)
		}
	}
}

func (s *SolidGold) ConsumeGophishFiles(gophishFilePaths ...string) {
	for _, gophishFile := range gophishFilePaths {
		csvFile, err := os.Open(gophishFile)
		util.CheckErr(err)
		reader := csv.NewReader(csvFile)
		//reader.Comma = ','
		reader.FieldsPerRecord = -1

		csvData, err := reader.ReadAll()
		util.CheckErr(err)

		csvFile.Close()

		for _, cols := range csvData {

			firstName := cols[0]
			lastName := cols[1]
			email := cols[2]
			position := cols[3]

			if strings.ToLower(firstName) == "First Name" {
				continue
			}

			human := s.Group.FindOrCreateHumanByEmail(email)
			human.AddName(fmt.Sprintf("%s %s", strings.TrimSpace(firstName), strings.TrimSpace(lastName)))
			human.AddRole(position)

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

func (s *SolidGold) ConsumeGithubOrgs(includeMembers bool, authToken string, orgs ...string) {

	client := github.NewClient(nil)

	if authToken != "" {
		client = client.WithAuthToken(authToken)
	}

	for _, org := range orgs {
		opt := &github.RepositoryListByOrgOptions{Type: "sources"}

		for {

			repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, opt)

			if err != nil {
				log.Println("Error encountered while attempting to consume the GitHub API", err)
				os.Exit(2)
			}

			for _, repo := range repos {
				gitCloneURL(path.Join(FolderPath, repo.GetFullName()), repo.GetSSHURL())
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}

		usernames := []string{}
		if includeMembers {
			memberListOpt := &github.ListMembersOptions{}
			for {

				members, resp, _ := client.Organizations.ListMembers(context.Background(), org, memberListOpt)
				for _, member := range members {
					if member.GetEmail() != "" {
						h := s.FindOrCreateHuman(member.GetLogin(), member.GetEmail())

						if member.GetName() != "" {
							h.AddName(member.GetName())
						}
					}

					usernames = append(usernames, member.GetLogin())
				}

				if resp.NextPage == 0 {
					break
				}

				opt.Page = resp.NextPage
			}
			s.ConsumeGithubUsers(false, authToken, usernames...)
		}

	}
}

func (s *SolidGold) ConsumeGithubUsers(includeOrgs bool, authToken string, users ...string) {
	client := github.NewClient(nil)

	if authToken != "" {
		client = client.WithAuthToken(authToken)
	}

	for _, user := range users {
		opt := &github.RepositoryListOptions{Type: "owner"}
		for {

			repos, resp, err := client.Repositories.List(context.Background(), user, opt)
			if err != nil {
				log.Println("Error encountered while attempting to consume the GitHub API", err)
				os.Exit(2)
			}

			for _, repo := range repos {
				if !*repo.Fork {
					gitCloneURL(path.Join(FolderPath, repo.GetFullName()), repo.GetCloneURL())
				}
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}

		if includeOrgs {
			orgs, _, err := client.Organizations.List(context.Background(), user, nil)
			if err != nil {
				log.Println("Error encountered while attempting to consume the GitHub API", err)
				os.Exit(2)
			}

			orgNames := []string{}
			for _, org := range orgs {
				orgNames = append(orgNames, org.GetName())
			}
			s.ConsumeGithubOrgs(false, authToken, orgNames...)
		}
	}
}

func (s *SolidGold) UpdateGithub() {
	orgOrUsers, _ := os.ReadDir(FolderPath)
	for _, orgOrUser := range orgOrUsers {
		if orgOrUser.IsDir() {
			repos, _ := os.ReadDir(orgOrUser.Name())
			for _, repo := range repos {
				if !repo.IsDir() {
					repoInst, err := git.PlainOpen(path.Join(FolderPath, orgOrUser.Name(), repo.Name()))
					if err != nil {
						log.Println("Error encountered while getting repo instance", err)
						continue
					}
					w, err := repoInst.Worktree()
					if err != nil {
						log.Println("Error encountered while getting repo work tree", err)
						continue
					}
					err = w.Pull(&git.PullOptions{RemoteName: "origin"})
					if err != nil {
						log.Println("Error encountered while pulling", err)
						continue
					}
				}
			}
		}
	}
}

func gitCloneURL(path, repoURL string) {
	repoExists, err := exists(path)
	if err != nil {
		return
	}
	if !repoExists {
		log.Printf("cloning %s to %s", repoURL, path)
		_, _ = git.PlainClone(path, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

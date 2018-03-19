package service

import (
	"fmt"
	"github.com/vastness-io/vcs-webhook-svc/webhook"
	"github.com/vastness-io/vcs-webhook-svc/webhook/bitbucketserver"
	"time"
)

const (
	//AddType is for files which have been added
	AddType = "ADD"
	//ModifyType is for files which have been modified
	ModifyType = "MODIFY"
	//DeleteType is for files which have been deleted
	DeleteType = "DELETE"
	//OrgType is organisation level
	OrgType = "ORG"
	//UserType is a user account
	UserType = "USER"
	//NormalType is a user account in Bitbucket Server terms
	NormalType = "NORMAL"
)

// MapPostWebhookToVcsPushEvent converts a Bitbucket Server post webhook message to a PushEvent
func MapPostWebhookToVcsPushEvent(postWebhook *bitbucketserver.PostWebhook) *vcs.VcsPushEvent {

	if postWebhook != nil {

		var (
			out              = new(vcs.VcsPushEvent)
			changesets       = postWebhook.GetChangesets()
			repository       = postWebhook.GetRepository()
			refChanges       = postWebhook.GetRefChanges()
			refChangesLength = len(postWebhook.GetRefChanges())
			headHash         *bitbucketserver.RefChange
		)

		out.VcsType = vcs.VcsType_BITBUCKET_SERVER

		if refChangesLength >= 1 {
			headHash = refChanges[refChangesLength-1]
		}

		out.HeadCommit = new(vcs.PushCommit)
		if changesets != nil {
			for _, value := range postWebhook.GetChangesets().GetValues() {

				pushCommit := new(vcs.PushCommit)

				changes := value.GetChanges()

				getToCommit := value.GetToCommit()

				if getToCommit != nil {
					pushCommit = &vcs.PushCommit{
						Sha:       getToCommit.GetId(),
						Id:        getToCommit.GetId(),
						Message:   getToCommit.GetMessage(),
						Timestamp: time.Unix(getToCommit.GetAuthorTimestamp(), 0).UTC().Format(time.RFC3339),
					}

				}

				author := getToCommit.GetAuthor()
				if author != nil {
					author := &vcs.CommitAuthor{
						Name:     author.GetName(),
						Email:    author.GetEmailAddress(),
						Username: author.GetName(),
						Date:     time.Unix(getToCommit.GetAuthorTimestamp(), 0).UTC().Format(time.RFC3339),
					}
					pushCommit.Author = author
					pushCommit.Committer = author
				}

				if changes != nil {
					for _, value := range changes.GetValues() {
						switch value.Type {
						case AddType:
							if value.Path != nil {
								pushCommit.Added = append(pushCommit.GetAdded(), value.Path.GetToString())
							}
						case ModifyType:
							if value.Path != nil {
								pushCommit.Modified = append(pushCommit.GetModified(), value.Path.GetToString())
							}
						case DeleteType:
							if value.Path != nil {
								pushCommit.Removed = append(pushCommit.GetRemoved(), value.Path.GetToString())
							}
						}
					}
				}

				if headHash != nil {
					out.Ref = headHash.RefId
					if pushCommit.GetId() == headHash.GetToHash() {
						out.HeadCommit = pushCommit
					}
				}

				out.Commits = append(out.GetCommits(), pushCommit)

			}
		}
		if repository != nil {
			outRepo := new(vcs.Repository)
			outRepo.Id = repository.GetId()
			outRepo.Name = repository.GetName()
			if repository.Project != nil {
				outRepo.FullName = fmt.Sprintf("%s/%s", repository.Project.GetName(), repository.GetName())
				outOrg := &vcs.User{}
				outOrg.Id = repository.Project.GetId()
				outOrg.Name = repository.Project.GetName()
				outOrg.Login = repository.Project.GetName()
				if repository.Project.GetType() == NormalType {
					outOrg.Type = OrgType
				} else {
					outOrg.Type = UserType
				}
				outRepo.Organization = outOrg
				out.Organization = outOrg
			}
			out.Repository = outRepo
		}
		return out
	}
	return nil
}

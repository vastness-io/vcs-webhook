package service

import (
	"github.com/vastness-io/vcs-webhook-svc/webhook"
	"github.com/vastness-io/vcs-webhook-svc/webhook/github"
	"time"
)

// MapPostWebhookToVcsPushEvent converts a Github push event message to a PushEvent
func MapGithubPushEventToVcsPushEvent(from *github.PushEvent) *vcs.VcsPushEvent {
	if from == nil {
		return nil
	}

	var (
		out            = new(vcs.VcsPushEvent)
		fromHeadCommit = from.GetHeadCommit()
		fromCommits    = from.GetCommits()
		fromPusher     = from.GetPusher()
		fromSender     = from.GetSender()
		fromRepository = from.GetRepository()
		timestamp      string
	)

	out.VcsType = vcs.VcsType_GITHUB

	for _, commit := range fromCommits {

		outCommit := new(vcs.PushCommit)

		if commit != nil {
			outCommit = &vcs.PushCommit{
				Sha:      commit.GetSha(),
				Id:       commit.GetId(),
				TreeId:   commit.GetTreeId(),
				Distinct: commit.GetDistinct(),
				Message:  commit.GetMessage(),
				Url:      commit.GetUrl(),
				Added:    commit.GetAdded(),
				Modified: commit.GetModified(),
				Removed:  commit.GetRemoved(),
			}

			ts, err := time.Parse(time.RFC3339, commit.GetTimestamp())

			if err != nil {
				log.WithError(err).WithField("timestamp", commit.GetTimestamp()).Debug("Unable to parse timestamp")
			} else {
				timestamp = ts.UTC().Format(time.RFC3339)
				outCommit.Timestamp = timestamp
			}

			outCommitter := commit.GetCommitter()

			if outCommitter != nil {
				outCommit.Committer = &vcs.CommitAuthor{
					Name:     outCommitter.GetName(),
					Email:    outCommitter.GetEmail(),
					Date:     timestamp,
					Username: outCommitter.GetUsername(),
				}
			}

			outAuthor := commit.GetAuthor()

			if outAuthor != nil {
				outCommit.Author = &vcs.CommitAuthor{
					Name:     outAuthor.GetName(),
					Email:    outAuthor.GetEmail(),
					Date:     timestamp,
					Username: outAuthor.GetUsername(),
				}
			}

			if fromHeadCommit != nil {
				if commit.GetId() == fromHeadCommit.GetId() {
					out.HeadCommit = outCommit
				}
			}

			out.Commits = append(out.GetCommits(), outCommit)

		}

	}

	out.Ref = from.GetRef()

	if fromPusher != nil {
		author := vcs.CommitAuthor{
			Name:     fromPusher.GetName(),
			Email:    fromPusher.GetEmail(),
			Date:     fromPusher.GetDate(),
			Username: fromPusher.GetUsername(),
		}
		out.Pusher = &author
	}

	if fromSender != nil {
		sender := vcs.User{
			Id:    int64(fromSender.GetId()),
			Login: fromSender.GetLogin(),
			Url:   fromSender.GetUrl(),
			Type:  fromSender.GetType(),
			Name:  fromSender.GetName(),
			Email: fromSender.GetEmail(),
		}
		out.Sender = &sender
	}

	if fromRepository != nil {
		outRepository := vcs.Repository{
			Id:          int64(fromRepository.GetId()),
			Name:        fromRepository.GetName(),
			FullName:    fromRepository.GetFullName(),
			Description: fromRepository.GetDescription(),
			Private:     fromRepository.GetPrivate(),
			Fork:        fromRepository.GetFork(),
			Url:         fromRepository.GetUrl(),
		}

		fromOwner := fromRepository.GetOwner()

		if fromOwner != nil {
			owner := vcs.User{
				Id:    int64(fromOwner.GetId()),
				Login: fromOwner.GetLogin(),
				Url:   fromOwner.GetUrl(),
				Type:  fromOwner.GetType(),
				Name:  fromOwner.GetName(),
				Email: fromOwner.GetEmail(),
			}
			outRepository.Owner = &owner
			outRepository.Organization = &owner
			out.Organization = &owner
		}

		fromOrg := fromRepository.GetOrganization()

		if fromOrg != nil {
			org := vcs.User{
				Id:    int64(fromOrg.GetId()),
				Login: fromOrg.GetLogin(),
				Url:   fromOrg.GetUrl(),
				Type:  fromOrg.GetType(),
				Name:  fromOrg.GetName(),
				Email: fromOrg.GetEmail(),
			}
			out.Organization = &org
		}

		out.Repository = &outRepository

	}



	out.Created = from.Created
	out.Deleted = from.Deleted
	out.Forced = from.Forced

	return out
}

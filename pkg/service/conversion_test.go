package service

import (
	"encoding/json"
	"github.com/vastness-io/vcs-webhook-svc/webhook"
	"github.com/vastness-io/vcs-webhook-svc/webhook/bitbucketserver"
	"github.com/vastness-io/vcs-webhook-svc/webhook/github"
	"reflect"
	"testing"
)

type convertTestHelper = struct {
	payload        []byte
	pushEvent      *github.PushEvent
	postWebhook    *bitbucketserver.PostWebhook
	ref            string
	commits        int
	filesAdded     []int
	filesModified  []int
	filesRemoved   []int
	distinct       []bool
	authorName     []string
	authorDate     []string
	authorUsername []string
	commitMessage  []string
	commitId       []string
	treeId         []string
	commitURL      []string
	organisation   []*vcs.User
	vcsType        vcs.VcsType
}

func TestPostWebhookConversion(t *testing.T) {

	tests := make([]convertTestHelper, len(postWebhookExamplePayloads))

	for i, _ := range postWebhookExamplePayloads {

		tests[i] = convertTestHelper{
			payload:     []byte(postWebhookExamplePayloads[i]),
			postWebhook: new(bitbucketserver.PostWebhook),
			ref:         "refs/heads/master",
			commits:     1,
			filesAdded: []int{
				0,
			},
			filesModified: []int{
				2,
			},
			filesRemoved: []int{
				0,
			},
			distinct: []bool{
				false,
			},
			authorName: []string{
				"jhocman",
			},
			authorDate: []string{
				"45531-05-09T13:13:20Z",
			},
			authorUsername: []string{
				"jhocman",
			},
			commitMessage: []string{
				"Updating poms ...",
			},
			commitId: []string{
				"f259e9032cdeb1e28d073e8a79a1fd6f9587f233",
			},
			organisation: []*vcs.User{
				{
					Id:    21,
					Name:  "Iridium",
					Login: "Iridium",
					Type:  "ORG",
				},
			},
			vcsType: vcs.VcsType_BITBUCKET_SERVER,
		}
	}

	for i, _ := range tests {

		if err := json.Unmarshal(tests[i].payload, tests[i].postWebhook); err != nil {
			t.Fatal(err)
		}

		pushEvent := MapPostWebhookToVcsPushEvent(tests[i].postWebhook)

		assertNotNil(t, pushEvent)

		assertEquals(t, tests[i].ref, pushEvent.GetRef())

		assertEquals(t, tests[i].vcsType, pushEvent.GetVcsType())

		assertNotNil(t, pushEvent.GetCommits())

		assertEquals(t, tests[i].commits, len(pushEvent.GetCommits()))

		commit := pushEvent.GetCommits()[i]

		assertNotNil(t, commit)

		assertEquals(t, tests[i].filesAdded[i], len(commit.GetAdded()))

		assertEquals(t, tests[i].filesModified[i], len(commit.GetModified()))

		assertEquals(t, tests[i].filesRemoved[i], len(commit.GetRemoved()))

		assertEquals(t, tests[i].distinct[i], commit.GetDistinct())

		commitAuthor := commit.GetAuthor()

		assertNotNil(t, commitAuthor)

		assertEquals(t, tests[i].authorName[i], commitAuthor.GetName())

		assertEquals(t, tests[i].authorDate[i], commitAuthor.GetDate())

		assertEquals(t, tests[i].authorUsername[i], commitAuthor.GetUsername())

		assertEquals(t, tests[i].commitId[i], commit.GetId())

		assertEquals(t, tests[i].commitMessage[i], commit.GetMessage())

		assertEquals(t, commit.GetId(), commit.GetSha())

		assertEquals(t, "", commit.GetTreeId())

		assertEquals(t, commitAuthor, commit.GetCommitter())

		assertEquals(t, "", commit.GetUrl())

		assertEquals(t, commit, pushEvent.GetHeadCommit())

		org := pushEvent.GetOrganization()

		assertNotNil(t, org)

		assertEquals(t, tests[i].organisation[i], org)

	}
}

func TestGithubPushEventConversion(t *testing.T) {

	tests := make([]convertTestHelper, len(githubPushPayloads))

	for i, _ := range githubPushPayloads {

		tests[i] = convertTestHelper{
			payload:   []byte(githubPushPayloads[i]),
			pushEvent: new(github.PushEvent),
			ref:       "refs/heads/changes",
			commits:   1,
			filesAdded: []int{
				0,
			},
			filesModified: []int{
				1,
			},
			filesRemoved: []int{
				0,
			},
			distinct: []bool{
				true,
			},
			authorName: []string{
				"baxterthehacker",
			},
			authorDate: []string{
				"2015-05-05T23:40:15Z",
			},
			authorUsername: []string{
				"baxterthehacker",
			},
			commitMessage: []string{
				"Update README.md",
			},
			commitId: []string{
				"0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
			},
			treeId: []string{
				"f9d2a07e9488b91af2641b26b9407fe22a451433",
			},
			commitURL: []string{
				"https://github.com/baxterthehacker/public-repo/commit/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
			},
			organisation: []*vcs.User{
				{
					Name:  "baxterthehacker",
					Email: "baxterthehacker@users.noreply.github.com",
				},
			},
			vcsType: vcs.VcsType_GITHUB,
		}
	}

	for i, _ := range tests {

		if err := json.Unmarshal(tests[i].payload, tests[i].pushEvent); err != nil {
			t.Fatal(err)
		}

		pushEvent := MapGithubPushEventToVcsPushEvent(tests[i].pushEvent)

		assertNotNil(t, pushEvent)

		assertEquals(t, tests[i].ref, pushEvent.GetRef())

		assertEquals(t, tests[i].vcsType, pushEvent.GetVcsType())

		assertNotNil(t, pushEvent.GetCommits())

		assertEquals(t, tests[i].commits, len(pushEvent.GetCommits()))

		commit := pushEvent.GetCommits()[i]

		assertNotNil(t, commit)

		assertEquals(t, tests[i].filesAdded[i], len(commit.GetAdded()))

		assertEquals(t, tests[i].filesModified[i], len(commit.GetModified()))

		assertEquals(t, tests[i].filesRemoved[i], len(commit.GetRemoved()))

		assertEquals(t, tests[i].distinct[i], commit.GetDistinct())

		commitAuthor := commit.GetAuthor()

		assertNotNil(t, commitAuthor)

		assertEquals(t, tests[i].authorName[i], commitAuthor.GetName())

		assertEquals(t, tests[i].authorDate[i], commitAuthor.GetDate())

		assertEquals(t, tests[i].authorUsername[i], commitAuthor.GetUsername())

		assertEquals(t, tests[i].commitId[i], commit.GetId())

		assertEquals(t, tests[i].commitMessage[i], commit.GetMessage())

		assertEquals(t, commit.GetId(), commit.GetSha())

		assertEquals(t, tests[i].treeId[i], commit.GetTreeId())

		assertEquals(t, commitAuthor, commit.GetCommitter())

		assertEquals(t, tests[i].commitURL[i], commit.GetUrl())

		assertEquals(t, commit, pushEvent.GetHeadCommit())

		org := pushEvent.GetOrganization()

		assertNotNil(t, org)

		assertEquals(t, tests[i].organisation[i], org)
	}
}

var (
	postWebhookExamplePayloads = []string{
		`{
		   "repository":{
			  "slug":"iridium-parent",
			  "id":11,
			  "name":"iridium-parent",
			  "scmId":"git",
			  "state":"AVAILABLE",
			  "statusMessage":"Available",
			  "forkable":true,
			  "project":{
				 "key":"IR",
				 "id":21,
				 "name":"Iridium",
				 "public":false,
				 "type":"NORMAL",
				 "isPersonal":false
			  },
			  "public":false
		   },
		   "refChanges":[
			  {
				 "refId":"refs/heads/master",
				 "fromHash":"2c847c4e9c2421d038fff26ba82bc859ae6ebe20",
				 "toHash":"f259e9032cdeb1e28d073e8a79a1fd6f9587f233",
				 "type":"UPDATE"
			  }
		   ],
		   "changesets":{
			  "size":1,
			  "limit":100,
			  "isLastPage":true,
			  "values":[
				 {
					"fromCommit":{
					   "id":"2c847c4e9c2421d038fff26ba82bc859ae6ebe20",
					   "displayId":"2c847c4"
					},
					"toCommit":{
					   "id":"f259e9032cdeb1e28d073e8a79a1fd6f9587f233",
					   "displayId":"f259e90",
					   "author":{
						  "name":"jhocman",
						  "emailAddress":"jhocman@atlassian.com"
					   },
					   "authorTimestamp":1374663446000,
					   "message":"Updating poms ...",
					   "parents":[
						  {
							 "id":"2c847c4e9c2421d038fff26ba82bc859ae6ebe20",
							 "displayId":"2c847c4"
						  }
					   ]
					},
					"changes":{
					   "size":2,
					   "limit":500,
					   "isLastPage":true,
					   "values":[
						  {
							 "contentId":"2f259b79aa7e263f5829bb6e98096e7ec976d998",
							 "path":{
								"components":[
								   "iridium-common",
								   "pom.xml"
								],
								"parent":"iridium-common",
								"name":"pom.xml",
								"extension":"xml",
								"toString":"iridium-common/pom.xml"
							 },
							 "executable":false,
							 "percentUnchanged":-1,
							 "type":"MODIFY",
							 "nodeType":"FILE",
							 "srcExecutable":false,
							 "link":{
								"url":"/projects/IR/repos/iridium-parent/commits/f259e9032cdeb1e28d073e8a79a1fd6f9587f233#iridium-common/pom.xml",
								"rel":"self"
							 }
						  },
						  {
							 "contentId":"2f259b79aa7e263f5829bb6e98096e7ec976d998",
							 "path":{
								"components":[
								   "iridium-magma",
								   "pom.xml"
								],
								"parent":"iridium-magma",
								"name":"pom.xml",
								"extension":"xml",
								"toString":"iridium-magma/pom.xml"
							 },
							 "executable":false,
							 "percentUnchanged":-1,
							 "type":"MODIFY",
							 "nodeType":"FILE",
							 "srcExecutable":false,
							 "link":{
								"url":"/projects/IR/repos/iridium-parent/commits/f259e9032cdeb1e28d073e8a79a1fd6f9587f233#iridium-magma/pom.xml",
								"rel":"self"
							 }
						  }
					   ],
					   "start":0,
					   "filter":null
					},
					"link":{
					   "url":"/projects/IR/repos/iridium-parent/commits/f259e9032cdeb1e28d073e8a79a1fd6f9587f233#iridium-magma/pom.xml",
					   "rel":"self"
					}
				 }
			  ],
			  "start":0,
			  "filter":null
		   }
		}
		`,
	}

	githubPushPayloads = []string{
		`{
  "ref": "refs/heads/changes",
  "before": "9049f1265b7d61be4a8904a9a27120d2064dab3b",
  "after": "0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
  "created": false,
  "deleted": false,
  "forced": false,
  "base_ref": null,
  "compare": "https://github.com/baxterthehacker/public-repo/compare/9049f1265b7d...0d1a26e67d8f",
  "commits": [
    {
      "id": "0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
      "tree_id": "f9d2a07e9488b91af2641b26b9407fe22a451433",
      "distinct": true,
      "message": "Update README.md",
      "timestamp": "2015-05-05T19:40:15-04:00",
      "url": "https://github.com/baxterthehacker/public-repo/commit/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
      "author": {
        "name": "baxterthehacker",
        "email": "baxterthehacker@users.noreply.github.com",
        "username": "baxterthehacker"
      },
      "committer": {
        "name": "baxterthehacker",
        "email": "baxterthehacker@users.noreply.github.com",
        "username": "baxterthehacker"
      },
      "added": [

      ],
      "removed": [

      ],
      "modified": [
        "README.md"
      ]
    }
  ],
  "head_commit": {
    "id": "0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
    "tree_id": "f9d2a07e9488b91af2641b26b9407fe22a451433",
    "distinct": true,
    "message": "Update README.md",
    "timestamp": "2015-05-05T19:40:15-04:00",
    "url": "https://github.com/baxterthehacker/public-repo/commit/0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
    "author": {
      "name": "baxterthehacker",
      "email": "baxterthehacker@users.noreply.github.com",
      "username": "baxterthehacker"
    },
    "committer": {
      "name": "baxterthehacker",
      "email": "baxterthehacker@users.noreply.github.com",
      "username": "baxterthehacker"
    },
    "added": [

    ],
    "removed": [

    ],
    "modified": [
      "README.md"
    ]
  },
  "repository": {
    "id": 35129377,
    "name": "public-repo",
    "full_name": "baxterthehacker/public-repo",
    "owner": {
      "name": "baxterthehacker",
      "email": "baxterthehacker@users.noreply.github.com"
    },
    "private": false,
    "html_url": "https://github.com/baxterthehacker/public-repo",
    "description": "",
    "fork": false,
    "url": "https://github.com/baxterthehacker/public-repo",
    "forks_url": "https://api.github.com/repos/baxterthehacker/public-repo/forks",
    "keys_url": "https://api.github.com/repos/baxterthehacker/public-repo/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/baxterthehacker/public-repo/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/baxterthehacker/public-repo/teams",
    "hooks_url": "https://api.github.com/repos/baxterthehacker/public-repo/hooks",
    "issue_events_url": "https://api.github.com/repos/baxterthehacker/public-repo/issues/events{/number}",
    "events_url": "https://api.github.com/repos/baxterthehacker/public-repo/events",
    "assignees_url": "https://api.github.com/repos/baxterthehacker/public-repo/assignees{/user}",
    "branches_url": "https://api.github.com/repos/baxterthehacker/public-repo/branches{/branch}",
    "tags_url": "https://api.github.com/repos/baxterthehacker/public-repo/tags",
    "blobs_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/baxterthehacker/public-repo/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/baxterthehacker/public-repo/languages",
    "stargazers_url": "https://api.github.com/repos/baxterthehacker/public-repo/stargazers",
    "contributors_url": "https://api.github.com/repos/baxterthehacker/public-repo/contributors",
    "subscribers_url": "https://api.github.com/repos/baxterthehacker/public-repo/subscribers",
    "subscription_url": "https://api.github.com/repos/baxterthehacker/public-repo/subscription",
    "commits_url": "https://api.github.com/repos/baxterthehacker/public-repo/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/baxterthehacker/public-repo/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/baxterthehacker/public-repo/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/baxterthehacker/public-repo/contents/{+path}",
    "compare_url": "https://api.github.com/repos/baxterthehacker/public-repo/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/baxterthehacker/public-repo/merges",
    "archive_url": "https://api.github.com/repos/baxterthehacker/public-repo/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/baxterthehacker/public-repo/downloads",
    "issues_url": "https://api.github.com/repos/baxterthehacker/public-repo/issues{/number}",
    "pulls_url": "https://api.github.com/repos/baxterthehacker/public-repo/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/baxterthehacker/public-repo/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/baxterthehacker/public-repo/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/baxterthehacker/public-repo/labels{/name}",
    "releases_url": "https://api.github.com/repos/baxterthehacker/public-repo/releases{/id}",
    "created_at": 1430869212,
    "updated_at": "2015-05-05T23:40:12Z",
    "pushed_at": 1430869217,
    "git_url": "git://github.com/baxterthehacker/public-repo.git",
    "ssh_url": "git@github.com:baxterthehacker/public-repo.git",
    "clone_url": "https://github.com/baxterthehacker/public-repo.git",
    "svn_url": "https://github.com/baxterthehacker/public-repo",
    "homepage": null,
    "size": 0,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": null,
    "has_issues": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": true,
    "forks_count": 0,
    "mirror_url": null,
    "open_issues_count": 0,
    "forks": 0,
    "open_issues": 0,
    "watchers": 0,
    "default_branch": "master",
    "stargazers": 0,
    "master_branch": "master"
  },
  "pusher": {
    "name": "baxterthehacker",
    "email": "baxterthehacker@users.noreply.github.com"
  },
  "sender": {
    "login": "baxterthehacker",
    "id": 6752317,
    "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
    "gravatar_id": "",
    "url": "https://api.github.com/users/baxterthehacker",
    "html_url": "https://github.com/baxterthehacker",
    "followers_url": "https://api.github.com/users/baxterthehacker/followers",
    "following_url": "https://api.github.com/users/baxterthehacker/following{/other_user}",
    "gists_url": "https://api.github.com/users/baxterthehacker/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/baxterthehacker/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/baxterthehacker/subscriptions",
    "organizations_url": "https://api.github.com/users/baxterthehacker/orgs",
    "repos_url": "https://api.github.com/users/baxterthehacker/repos",
    "events_url": "https://api.github.com/users/baxterthehacker/events{/privacy}",
    "received_events_url": "https://api.github.com/users/baxterthehacker/received_events",
    "type": "User",
    "site_admin": false
  }
}`,
	}
)

func assertNotNil(t *testing.T, actual interface{}) {
	if actual == nil {
		t.Fatal("Should not equal nil")
	}
}

func assertEquals(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

}

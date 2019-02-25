

# Lists out the pull requests
listPulls = codecommit.list_pull_requests(
    repositoryName=repository,
    authorArn=auth,
    maxResults=5
)

# Gets the current pull request
currentPullRequest = codecommit.get_pull_request(
    pullRequestId= listPulls['pullRequestIds'][0]
)

comments = codecommit.get_comments_for_pull_request(
    pullRequestId=pullRequest
)
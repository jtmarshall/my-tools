import json
import boto3
from botocore.vendored import requests


# Codecommit client
codecommit = boto3.client('codecommit')

def lambda_handler(event, context):
    # Log the updated references from the event
    references = { reference['ref'] for reference in event['Records'][0]['codecommit']['references'] }
    print("References: "  + str(references))
    
    # Get the commit ID
    commits = [ reference['commit'] for reference in event['Records'][0]['codecommit']['references'] ]
    print("CommitID: "  + commits[0])
    
    #Get the repository from the event and show its git clone URL
    repository = event['Records'][0]['eventSourceARN'].split(':')[5]
    
    # Get commit info
    commitData = codecommit.get_commit(
        repositoryName = repository,
        commitId = commits[0]
    )
    auth = commitData['commit']['author']['name']
    # msg = commitData['commit']['message']
    
    # Lists out the pull requests
    listPulls = codecommit.list_pull_requests(
        repositoryName = repository,
        maxResults = 5
    )
    
    # Gets the current pull request, 
    currentPullRequest = codecommit.get_pull_request(
        pullRequestId = listPulls['pullRequestIds'][0]
    )
    
    print('Current pull request: ', currentPullRequest)
    
    # Get title and description of pull request
    pullRequestTitle = currentPullRequest['pullRequest']['title']
    
    # Handle empty pull descriptions
    try:
        pullRequestDescription = currentPullRequest['pullRequest']['description']
    except:
        pullRequestDescription = ""
    
    # Posting to slack
    try:
        # Use slack's incomming webhook url
        webhook_url = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
        # Format the message to send in slack
        slack_data = {'text': "*`Deploying " + repository + "`*: " + pullRequestTitle + " " + pullRequestDescription + "\n- " + auth}
        
        # Make the post to the webhook
        response = requests.post(
            webhook_url, data=json.dumps(slack_data),
            headers={'Content-Type': 'application/json'}
        )
        if response.status_code != 200:
            raise ValueError(
                'Request to slack returned an error %s, the response is:\n%s'
                % (response.status_code, response.text)
            )
        
        response = codecommit.get_repository(repositoryName=repository)
        print("Clone URL: " +response['repositoryMetadata']['cloneUrlHttp'])
        return response['repositoryMetadata']['cloneUrlHttp']
    except Exception as e:
        print(e)
        print('Error getting repository {}. Make sure it exists and that your repository is in the same region as this function.'.format(repository))
        raise e

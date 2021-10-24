# Basecamp Totoro Bot (Golang)

This is a Basecamp chatbot that posts the latest Mailchimp member stats into the
chat room every day.

It is implemented as a [Google Cloud Run](https://cloud.google.com/run/docs/quickstarts/build-and-deploy/go)
containerized application that is mostly generic i.e. could be run on any
container platform. However, it does make use of [Google Cloud Secret Manager](https://cloud.google.com/run/docs/configuring/secrets) for storing the Mailchimp API Key and Basecamp Chatbot URL.

## Setup

1. Install the [Cloud Code for VS Code extension](https://cloud.google.com/code/docs/vscode/install). This includes the dependencies including the local test environment.
2. [Set up the local Cloud Run service](https://cloud.google.com/code/docs/vscode/developing-a-cloud-run-service) using the Docker builder.
3. Ensure you have logged into the Google Cloud project to be able to access the relevant secrets.

## Configuration notes

* [Two secrets are set up](https://console.cloud.google.com/security/secret-manager?project=bc-totorobot-go) with the Cloud Run project service account given the appropriate read permissions. The [Cloud Code extension for VS Code](https://cloud.google.com/code/docs/vscode/secret-manager) can handle accessing them in local development.
* [Google Cloud Scheduler is set up](https://cloud.google.com/run/docs/triggering/using-scheduler) to trigger the service each day. It also has its own service account.

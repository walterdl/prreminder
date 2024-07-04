## Description

PRReminder monitors a Slack channel for pull request (PR) messages, initiating a series of steps using AWS Step Functions to ensure timely PR reviews. It decouples event processing using AWS Lambda and SQS, abstracts PR links, and manages state machines for each PR to periodically check their status. The system sends reminders via Slack if PRs remain unapproved within a specified timeframe, respecting business hours and optimizing the review process by integrating seamlessly with Slack and GitLab APIs.

![Architecture](.doc/design.jpg)

## Local Development

### Requirements

1. Install [Task](https://taskfile.dev/) Runner.
2. Install [CDK CLI](https://docs.aws.amazon.com/cdk/v2/guide/cli.html).
3. Install [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
4. Download Go modules by running `go get` in the root directory.
5. Install CompileDaemon locally for compilation in watch mode. Run `task install_compiledaemon` in the root. This should create the program file in the `local` directory.

### Compilation

Before running lambda Functions locally, the Go functions must be compiled. There are two compilation modes: `prod` and `watch-mode`.

- `prod`: From each Lambda function directory, run `task build`. This compilation disables debugging utilities and removes the RPC layer from the `aws-lambda-go` for a smaller deployment unit.
- `watch-mode`: From each Lambda function directory, run `task watch`. This re-compiles the function every time a Go file is modified.

### Local API

Powered by the SAM and CDK integration. Run `task start-api` in the root. It should expose a local web server in port `3000`. The output provides the different endpoints for the Lambda functions. This requires the Lambda functions to be built.


### Local Worker

- Run `task start-workers`
- `cd local/invoke-worker`
- Create a JS file from the template file and edit the Lambda invocation payload.
- Run the file with `node FILE.js`

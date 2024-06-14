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


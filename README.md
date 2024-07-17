# Coinbase Prime Activity Listener

[![GoDoc](https://godoc.org/github.com/coinbase-samples/prime-activity-listener-go?status.svg)](https://godoc.org/github.com/coinbase-samples/prime-activity-listener-go)
[![Go Report Card](https://goreportcard.com/badge/coinbase-samples/prime-activity-listener-go)](https://goreportcard.com/report/coinbase-samples/prime-activity-listener-go)


## Overview

The *Coinbase Prime Activity Listener* is a sample application that demonstrates how to poll for new portfolio [activities](https://docs.cdp.coinbase.com/prime/reference/primerestapi_getportfolioactivities) and broadcast them to external services for consumption. This sample application writes the activity data to the [Amazon Simple Notification Service](https://docs.aws.amazon.com/sns/latest/dg/welcome.html) (SNS) service and there are two subscriptions to the topic. The first subscription is an [Amazon Simple Queue Service](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/welcome.html) (SQS) queue and the second is an [Amazon Data Firehose](https://docs.aws.amazon.com/firehose/latest/dev/what-is-this-service.html) (Firehose) stream. The Firehose stream persists the activities to an [Amazon S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html) (S3) bucket, but the service also supports persisting data to Amazon OpenSearch, Amazon Kinesis Data Streams, Amazon Managed Streaming for Apache Kafka, Snowflake, Splunk, Redshift, and custom HTTP endpoints.

## License

The *Coinbase Prime Activity Listener* sample application is free and open source and released under the [Apache License, Version 2.0](LICENSE).

The application and code are only available for demonstration purposes.

## Warning

If this application is deployed using the sample AWS CloudFormation template, there will be new charges to your AWS account. For high throughput Prime portfolios, these charges may be significant. As always, continiously review your AWS bill to understand more.

## Usage

### Create Stack

The *Coinbase Prime Activity Listener* has a [sample AWS CloudFormation](infra/aws.cfn.yml) (CFN) template that can be deployed to run the application in an Amazon Elastic Container Service (Amazon ECS) cluster. This template creates all of the required resources and can be customized to suit the deployers needs. To deploy the CFN stack, [initialize your AWS credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) and then run:

 ```bash
make create-aws-stack ENV_NAME=dev PROFILE=default REGION=us-east-1
```

Customize the values of the *ENV_NAME*, *PROFILE*, and *REGION* to match the needs of your environment. The *PROFILE* argument is the name of the AWS CLI profile configured and the *REGION* argument is the [AWS Region](https://aws.amazon.com/about-aws/global-infrastructure/regions_az/) you would like to deploy to.

### Prime Credentials

Once the CFN stack is deployed, configure your credentials in [AWS Secrets Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html). The name of the empty secret uses the following format:

```
prime-activity-listener-ENV_NAME-prime-api-credentials
```

The *ENV_NAME* will the same as what was passed to the *create-aws-stack* command (e.g., dev). The credentials use the following format:

```
{
  "accessKey": "",
  "passphrase": "",
  "signingKey": "",
  "portfolioId": "",
  "svcAccountId": ""
}
```

Prime API credentials can be created in the [Prime web application](https://prime.coinbase.com), once an account is opened.

### Build/Deploy Container

The inital stack deploys the *public.ecr.aws/nginx/nginx:stable-perl-arm64v8* container as a placeholder. The ECS task definition/service requires a container and at this point, the ECR repository has not been created. Once the CFN stack is deployed, build the container image and deploy to the Amazon ECR (ECR) repository. To do this, execute:

 ```bash
make build-image ENV_NAME=dev PROFILE=default REGION=us-east-1 BUILD_ID=1
```

Again, customize the values of the *ENV_NAME*, *PROFILE*, *BUILD_ID*, and *REGION* to match the needs of your environment. The *BUILD_ID* can be set to a value that matches your container tagging practices.

After the build is deployed, update the *DockerImageUri* attribute in the [aws-dev.json](infra/aws-dev.json) file to the URI of your deployed build. The newly created ECR respository uses the following format:

```
AWS_ACCOUNT_ID.dkr.ecr.AWS_REGION.amazonaws.com/prime-activity-listener-ENV_NAME:BUILD_ID
```

Note: If your environment is not named, *dev*, you will need to create a new environment configuration file. The environment name is defined by the *ENV_NAME* CLI argument. The format for the environment configuration file name is:

```
infra/aws-ENV_NAME.json
```

The CFN stack is configured to run ARM64 containers, so you must build on an ARM64 compatible computer.

### Update Stack

Once the *DockerImageUri* attribute is updated, run the following command to update the CFN stack and run the *Coinbase Prime Activity Listener*:

 ```bash
make update-aws-stack ENV_NAME=dev PROFILE=default REGION=us-east-1
```

This command deploys the container image specified in the *DockerImageUri* and starts listening for new Prime activities.

To validate, perform some actions in the Coinbase Prime application, and then review the SQS queue and/or the S3 bucket where Firehose persists.

Note: The buffering hints for Firehose are set to the following, but these can be customized in the CFN template by adjusting the *FirehoseBufferingHintIntervalInSeconds* and/or *FirehoseBufferingHintSizeInMBs* parameters.

```
FirehoseBufferingHintSizeInMBs: 128
FirehoseBufferingHintIntervalInSeconds: 900
```

## Build

To build the sample application locally, ensure that [Go](https://go.dev/) 1.21+ is installed and then run:

```bash
go build cmd/main.go
```


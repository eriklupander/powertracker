# Powertracker
Keeping track of electricity usage using AWS services, Tibber's API and the Watty energy monitor/load balancer.

## AWS services
Uses the following services:
* AWS API Gateway
* AWS Lambda
* AWS EventBridge for scheduling
* AWS Timestream for time-series data
* AWS Secret Manager for storing Tibber API key

All lambdas are written in Golang.

AWS CDK is used for all provisioning, building the Go lambda's etc.

## CDK basics

This is a project for TypeScript development with CDK.

The `cdk.json` file tells the CDK Toolkit how to execute your app.

## Useful commands

 * `npm run build`   compile typescript to js
 * `npm run watch`   watch for changes and compile
 * `npm run test`    perform the jest unit tests
 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template

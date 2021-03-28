import * as cdk from '@aws-cdk/core';
import * as path from "path";
import * as lambda from '@aws-cdk/aws-lambda';
import {Schedule} from "@aws-cdk/aws-events"
import {HttpApi, HttpMethod} from "@aws-cdk/aws-apigatewayv2"
import {LambdaProxyIntegration} from "@aws-cdk/aws-apigatewayv2-integrations"
import iam = require("@aws-cdk/aws-iam");

const ruleCdk = require('@aws-cdk/aws-events')
const targets = require('@aws-cdk/aws-events-targets')

export class PowertrackerStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // set up Timestream DB
        // const timeStreamDB = new PowerTrackerTimestreamConstruct(this, "powertracker_timestream", {
        //     databaseName: "powertracker",
        //     tableName: "power_record"
        // })


        // Build powerecorder lambda
        const powerRecorderFunction = this.buildGolangLambda('powerRecorder', path.join(__dirname, '../functions/powerRecorder'), 'main');

        // Build EventBridge rule with cron expression and bind to lambda
        const rule = new ruleCdk.Rule(this, "collect_power_rule", {
            description: "Invoked every minute to collect current power state",
            schedule: Schedule.expression("cron(0/5 * * * ? *)")
        });
        rule.addTarget(new targets.LambdaFunction(powerRecorderFunction))

        // Build API lambda
        const policyStmt = new iam.PolicyStatement()
        policyStmt.addActions(
            "secretsmanager:GetSecretValue",
            "timestream:*"
        )
        policyStmt.addResources(
            "*"
            // "arn:aws:secretsmanager:*:secret:prod/tibber_config-*",
            // "arn:aws:timestream:eu-west-1:378539896247:database/powertracker"
        )

        powerRecorderFunction.addToRolePolicy(policyStmt)

        const exporterLambdaFn = this.buildGolangLambda('exporter-api', path.join(__dirname, '../functions/exporter'), 'main');
        const policyStmt2 = new iam.PolicyStatement()
        policyStmt2.addActions(
            "secretsmanager:GetSecretValue",
            "timestream:*",
        )

        policyStmt2.addResources(
            "*",
            "arn:aws:secretsmanager:*:secret:prod/tibber_config-*",
        )
        exporterLambdaFn.addToRolePolicy(policyStmt2)

        // Create Rest API Gateway in front of the lambda
        const apiGtw = this.createApiGatewayForLambda("exporter-api-endpoint", exporterLambdaFn, 'Exposed endpoint')

        // Output the DNS of your API gateway deployment
        new cdk.CfnOutput(this, 'lambda-url', {value: apiGtw.url!})
    }

    // buildGolangLambda builds a docker image with the code and creates the lambda function.
    buildGolangLambda(id: string, lambdaPath: string, handler: string): lambda.Function {
        const environment = {
            CGO_ENABLED: '0',
            GOOS: 'linux',
            GOARCH: 'amd64',
        };
        return new lambda.Function(this, id, {
            code: lambda.Code.fromAsset(lambdaPath, {
                bundling: {
                    image: lambda.Runtime.GO_1_X.bundlingDockerImage,
                    user: "root",
                    environment,
                    command: [
                        'bash', '-c', [
                            'make lambda-build',
                        ].join(' && ')
                    ]
                }
            }),
            handler,
            runtime: lambda.Runtime.GO_1_X,
            timeout: cdk.Duration.seconds(10),
        });
    }

    // createApiGatewayForLambda creates a HTTP API Gateway for the supplied lambda function
    createApiGatewayForLambda(id: string, handler: lambda.Function, desc: string): HttpApi {

        const httpApi = new HttpApi(this, id, {
            description: desc,
        })
        const lambdaProxyIntegration = new LambdaProxyIntegration({
            handler: handler,
        })
        httpApi.addRoutes({
            integration: lambdaProxyIntegration,
            methods: [HttpMethod.GET],
            path: '/',
        })
        return httpApi
    }
}
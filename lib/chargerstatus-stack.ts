import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as path from "path";
import {Schedule} from "aws-cdk-lib/aws-events"
import {GolangBuilder} from "./golanglambda-builder";
import {APIGatewayBuilder} from "./apigw-builder";
import iam = require("aws-cdk-lib/aws-iam");

const ruleCdk = require('aws-cdk-lib/aws-events')
const targets = require('aws-cdk-lib/aws-events-targets')

    export class ChargerStatusStack extends cdk.Stack {
        constructor(scope: Construct, id: string, props?: cdk.StackProps) {
            super(scope, id, props);

        // IAM policies
        const secretsPolicy = new iam.PolicyStatement({
            actions: ["secretsmanager:GetSecretValue"],
            resources: ["arn:aws:secretsmanager:*:secret:prod/tibber_config-*"]
        })

        // Build lambda that reads data from Chargefinder and stores in InfluxDB cloud
        const golangBuilder = new GolangBuilder(this, "golang builder");
        const chargerStatusLambdaFn = golangBuilder
            .buildGolangLambda('chargerStatus', path.join(__dirname, '../functions/statusrecorder'), 60);

        // Build EventBridge rule with cron expression and bind to lambda to trigger chargerStatus lambda
        const rule = new ruleCdk.Rule(this, "collect_charger_status_rule", {
            description: "Invoked every 15 minutes to collect current charger state",
            schedule: Schedule.expression("cron(0/15 * * * ? *)")
        });
        rule.addTarget(new targets.LambdaFunction(chargerStatusLambdaFn))

        // Add IAM for chargerStatus recorder
        chargerStatusLambdaFn.addToRolePolicy(secretsPolicy)

        // Build Status API lambda
        const statusApiLambdaFn = golangBuilder.buildGolangLambda('status-api', path.join(__dirname, '../functions/statusapi'), 30);

        // Add IAM for statusApiLambdaFn recorder
        statusApiLambdaFn.addToRolePolicy(secretsPolicy)

        // Create HTTP API Gateway in front of the lambda
        const apiGtw = new APIGatewayBuilder(scope, 'api-gw-status')
            .createApiGatewayForLambda("status-api-endpoint", statusApiLambdaFn, this);

        // Output the hostname of your the API gateway
        new cdk.CfnOutput(this, 'lambda-url', {value: apiGtw.url!})
    }
}
import * as cdk from '@aws-cdk/core';
import * as path from "path";
import {Schedule} from "@aws-cdk/aws-events"
import {GolangBuilder} from "./golanglambda-builder";
import iam = require("@aws-cdk/aws-iam");

const ruleCdk = require('@aws-cdk/aws-events')
const targets = require('@aws-cdk/aws-events-targets')

export class ChargerStatusStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // IAM policies
        const secretsPolicy = new iam.PolicyStatement({
            actions: ["secretsmanager:GetSecretValue"],
            resources: ["arn:aws:secretsmanager:*:secret:prod/tibber_config-*"]
        })

        // Build lambda that reads data from Chargefinder and stores in InfluxDB cloud
        const chargerStatusFunction = new GolangBuilder(this, "golang builder")
            .buildGolangLambda('chargerStatus', path.join(__dirname, '../functions/statusrecorder'), 60);

        // Build EventBridge rule with cron expression and bind to lambda to trigger chargerStatus lambda
        const rule = new ruleCdk.Rule(this, "collect_charger_status_rule", {
            description: "Invoked every 15 minutes to collect current charger state",
            schedule: Schedule.expression("cron(0/15 * * * ? *)")
        });
        rule.addTarget(new targets.LambdaFunction(chargerStatusFunction))

        // Add IAM for powerrecorder
        chargerStatusFunction.addToRolePolicy(secretsPolicy)
     }
}
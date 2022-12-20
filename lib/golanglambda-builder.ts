import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from 'constructs';
import {Duration} from "aws-cdk-lib";

export class GolangBuilder extends Construct {

    buildGolangLambda(id: string, lambdaPath: string, timeout: number): lambda.Function {
        const environment = {
            CGO_ENABLED: '0',
            GOOS: 'linux',
            GOARCH: 'amd64',
        };
        return new lambda.Function(this, id, {
            code: lambda.Code.fromAsset(lambdaPath, {
                bundling: {
                    image: lambda.Runtime.GO_1_X.bundlingImage,
                    user: "root",
                    environment,
                    command: [
                        'bash', '-c', [
                            'make lambda-build',
                        ].join(' && ')
                    ]
                }
            }),
            handler: 'main',
            runtime: lambda.Runtime.GO_1_X,
            timeout: Duration.seconds(timeout),
           // timeout: cdk.Duration.seconds(timeout),
        });
    }
}

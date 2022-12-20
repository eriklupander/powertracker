import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from 'constructs';
import {HttpApi, HttpMethod, PayloadFormatVersion} from "@aws-cdk/aws-apigatewayv2-alpha";
import {HttpLambdaIntegration} from "@aws-cdk/aws-apigatewayv2-integrations-alpha";
import {Stack} from "aws-cdk-lib";

export class APIGatewayBuilder extends Construct{
    // createApiGatewayForLambda creates a HTTP API Gateway for the supplied lambda function.
    createApiGatewayForLambda(id: string, handler: lambda.Function, scope: Stack): HttpApi {

        const httpIntegration = new HttpLambdaIntegration(id, handler);

        const httpApi = new HttpApi(scope, 'HttpApi');

        httpApi.addRoutes({
            path: '/books',
            methods: [ HttpMethod.GET ],
            integration: httpIntegration,
        });
        return httpApi

        // const httpApi = new HttpApi(this, id, {
        //     description: desc,
        // })
        // const httpIntegration = new HttpLambdaIntegration(id, handler)
        // /*const lambdaProxyIntegration = new HttpLambdaIntegration({
        //     handler: handler,
        //     payloadFormatVersion: PayloadFormatVersion.VERSION_1_0,
        // })*/
        // httpApi.addRoutes({
        //     integration: httpIntegration,
        //     methods: [HttpMethod.GET],
        //     path: '/',
        // })
        // return httpApi
    }
}
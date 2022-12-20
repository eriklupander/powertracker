#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { ChargerStatusStack } from '../lib/chargerstatus-stack';

const app = new cdk.App();
new ChargerStatusStack(app, 'ChargerStatusStack');
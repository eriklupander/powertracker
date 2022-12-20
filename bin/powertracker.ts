#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { PowertrackerStack } from '../lib/powertracker-stack';

const app = new cdk.App();
new PowertrackerStack(app, 'PowertrackerStack');
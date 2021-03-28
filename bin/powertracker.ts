#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { PowertrackerStack } from '../lib/powertracker-stack';

const app = new cdk.App();
new PowertrackerStack(app, 'PowertrackerStack');
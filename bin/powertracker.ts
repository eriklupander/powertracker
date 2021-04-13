#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { PowertrackerStack } from '../lib/powertracker-stack';

const app = new cdk.App();
new PowertrackerStack(app, 'PowertrackerStack');
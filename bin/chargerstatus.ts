#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { ChargerStatusStack } from '../lib/chargerstatus-stack';

const app = new cdk.App();
new ChargerStatusStack(app, 'ChargerStatusStack');
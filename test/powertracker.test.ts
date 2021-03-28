import { expect as expectCDK, matchTemplate, MatchStyle } from '@aws-cdk/assert';
import * as cdk from '@aws-cdk/core';
import * as Powertracker from '../lib/powertracker-stack';

test('Empty Stack', () => {
    const app = new cdk.App();
    // WHEN
    const stack = new Powertracker.PowertrackerStack(app, 'PowertrackerStack');
    // THEN
    expectCDK(stack).to(matchTemplate({
      "Resources": {}
    }, MatchStyle.EXACT))
});


import { configureGenkit } from '@genkit-ai/core';
import { startFlowsServer } from '@genkit-ai/flow';
import { vertexAI } from '@genkit-ai/vertexai';
import { firebase } from '@genkit-ai/firebase';
import { dotprompt, promptRef } from '@genkit-ai/dotprompt';

const LOCATION = process.env.LOCATION|| 'us-central1';
const PROJECT_ID = process.env.PROJECT_ID;

configureGenkit({
  plugins: [
  
    vertexAI({ projectId: PROJECT_ID, location: LOCATION }),
    firebase(),
    dotprompt()
  ],
  logLevel: 'debug',
  enableTracingAndMetrics: true,
    telemetry: {
    instrumentation: 'firebase',
    logger: 'firebase',
    }
});


export {UserProfileFlowPrompt, UserProfileFlow} from './userProfileFlow'
export {QueryTransformPrompt, QueryTransformFlow} from './queryTransformFlow'
export {MovieFlowPrompt, MovieFlow} from './movieFlow'
export {movieDocFlow} from './docRetriever'


// Start a flow server, which exposes your flows as HTTP endpoints. This call
// must come last, after all of your plug-in configuration and flow definitions.
// You can optionally specify a subset of flows to serve, and configure some
// HTTP server options, but by default, the flow server serves all defined flows.
startFlowsServer();
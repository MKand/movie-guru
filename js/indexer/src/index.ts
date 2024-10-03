
import { configureGenkit } from '@genkit-ai/core';
import { startFlowsServer } from '@genkit-ai/flow';
import { vertexAI } from '@genkit-ai/vertexai';
import { firebase } from '@genkit-ai/firebase';
import { ProcessMovies } from './addData';

let location = process.env.LOCATION;
if (!location) {
  location = 'europe-west4';
}
configureGenkit({
  plugins: [
  
    vertexAI({ projectId: process.env.PROJECT_ID, location: location }),
    firebase()
  ],
  // Log debug output to tbe console.
  logLevel: 'debug',
  // Perform OpenTelemetry instrumentation and enable trace collection.
  enableTracingAndMetrics: true,
    telemetry: {
    instrumentation: 'firebase',
    logger: 'firebase',
    }
});

ProcessMovies()


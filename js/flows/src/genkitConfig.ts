import { gemini15Flash, vertexAI,  } from '@genkit-ai/vertexai';
import { enableFirebaseTelemetry} from '@genkit-ai/firebase';
import { initializeApp } from 'firebase-admin/app';

import { genkit } from 'genkit';

const LOCATION = process.env.LOCATION|| 'us-central1';
const PROJECT_ID = process.env.PROJECT_ID;


initializeApp({
  projectId: PROJECT_ID,
});

enableFirebaseTelemetry();


export const ai = genkit({
    plugins: [vertexAI({location: LOCATION, projectId: PROJECT_ID})],
    model: gemini15Flash, // set default model
  });
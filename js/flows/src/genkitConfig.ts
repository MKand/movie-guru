import { gemini15Flash, vertexAI, textEmbedding004  } from '@genkit-ai/vertexai';
import { enableFirebaseTelemetry} from '@genkit-ai/firebase';
import { initializeApp } from 'firebase-admin/app';
import { genkitEval, GenkitMetric } from '@genkit-ai/evaluator';

import { genkit } from 'genkit';

const LOCATION = process.env.LOCATION|| 'us-central1';
const PROJECT_ID = process.env.PROJECT_ID;


initializeApp({
  projectId: PROJECT_ID,
});

enableFirebaseTelemetry();


export const ai = genkit({
    plugins: [
      vertexAI({location: LOCATION, projectId: PROJECT_ID}),
      genkitEval({
        judge: gemini15Flash,
        metrics: [GenkitMetric.FAITHFULNESS, GenkitMetric.ANSWER_RELEVANCY],
        embedder: textEmbedding004, // GenkitMetric.ANSWER_RELEVANCY requires an embedder
      }),
    ],
    model: gemini15Flash, // set default model
    
  });
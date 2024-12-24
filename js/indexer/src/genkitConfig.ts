import {  vertexAI } from '@genkit-ai/vertexai';
import { genkit } from 'genkit';

const LOCATION = process.env.LOCATION|| 'us-central1';
const PROJECT_ID = process.env.PROJECT_ID ;

export const ai = genkit({
    plugins: [vertexAI({location: LOCATION, projectId: PROJECT_ID})]
  });
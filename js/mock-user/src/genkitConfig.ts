/**
 * Copyright 2025 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
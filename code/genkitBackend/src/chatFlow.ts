// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { ai } from './genkitConfig';
import { z } from 'genkit';
import {userPreferenceAgent} from './userPreferenceAgent'

// Define a prompt that represents a specialist agent
const chatAgent = ai.chat(
  {
    input: {
      schema: z.object({
      userProfile: z.any().describe('A list of user preferences with likes and dislikes categorized by actor, director, genre, or a catch all category called other.'),
      movieContext: z.array(z.any()).describe('A list of movies that have been retrieved from the database and are relevant to the user\'s query.'),
      userMessage: z.string().describe('The original message sent my the user'),
    })},
    system: 'say hi and try your best to respond',
    tools: [userPreferenceAgent],
  },
);


export const chatFlow = ai.defineFlow(
    {
        name: 'chatFlow',
        inputSchema: z.string(),
        outputSchema: z.string(),
    },
    async (message) => {
        const { text } = await chatAgent.send(message);
        return text
    }
);
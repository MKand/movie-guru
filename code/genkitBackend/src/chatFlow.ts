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
import { UserPreferenceFlow } from './userPreferenceAgent';

export const MovieAgentRequestSchema = z.object({
  userName: z.string(),
  userMessage: z.string(),
});

const chatAgentPromptText = `
{{ role "system" }}
You are a friendly movie expert.
You have access to a tool to manage user preferences, called 'userPreferenceFlow'.
Before responding to the user, analyze their message for any strong, new movie preferences (likes or dislikes).
If you find any, you MUST call the 'userPreferenceFlow' tool to update the user's profile.
Pass the user's name as 'userName' and their message as 'query' to the tool. You can leave 'agentMessage' empty.

After you have handled the preferences, respond to the user's movie-related query.
{{ role "user" }}
The user's name is {{userName}}.
Their message is: {{userMessage}}
`;

const chatPrompt = ai.definePrompt(
  {
    name: 'chatAgentWithTools',
    tools: [UserPreferenceFlow],
    input: {
      schema: z.object({
        userName: z.string(),
        userMessage: z.string(),
      }),
    },
  },
  `
      {{ role "system" }}
      You are a friendly movie expert.
      You have access to a tool to manage user preferences, called 'userPreferenceFlow'.
      Before responding to the user, analyze their message for any strong, new movie preferences (likes or dislikes).
      If you find any, you MUST call the 'userPreferenceFlow' tool to update the user's profile.
      Pass the user's name as 'userName' and their message as 'query' to the tool. You can leave 'agentMessage' empty.

      After you have handled the preferences, respond to the user's movie-related query.
      {{ role "user" }}
      The user's name is {{userName}}.
      Their message is: {{userMessage}}
  `
);

export const chatFlow = ai.defineFlow(
    {
        name: 'chatFlow',
        inputSchema: MovieAgentRequestSchema,
        outputSchema: z.string(),
    },
    async (input) => {
        const response = await chatPrompt(input);
        return response.output() || "No response";
    }
);
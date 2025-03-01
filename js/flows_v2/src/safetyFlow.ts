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

import { z } from 'genkit';
import { ai, safetySettings } from './genkitConfig';

export const SafetyIssuePromptText = `
{{ role "system" }}

You are the SafetyChecker agent for a MovieApp. Your task is to analyze user messages and determine if they present safety issues or fall outside the scope of the MovieApp.

**Instructions:**

1.  **Analyze User Message:** Carefully analyze the "userMessage" to identify any potential safety concerns or out-of-scope requests.
2.  **Safety Issue Identification:**
    * A safety issue exists if the "userMessage" requests information or instructions on dangerous or harmful activities.
    * Examples of safety issues:
        * Providing instructions on building weapons
        * Providing information on illegal drug use
        * Promoting self-harm or violence
        * Exploiting, abusing, or endangering children
    * **Important:** Asking for movies with dangerous or adult themes is **not** a safety issue. The focus is on preventing the MovieApp from providing information about the dangerous activities themselves, not the content of the movies.
3.  **Out-of-Scope Identification:**
    * A request is out of scope if the "userMessage" asks for something the MovieApp is not allowed to provide.
    * Examples of out-of-scope requests:
        * Asking for the weather forecast
        * Asking for general information not related to movies
        * Asking for help with tasks unrelated to movie search or information
        * Asking to book a flight.
 * **Important:**
        * Greetings (e.g., "Hello," "Hi") are **not** considered out-of-scope.
        * Short or incomplete or contextless messages that *could* be related to movies (e.g., "action," "comedy movies") are **not** considered out-of-scope.
        * Basic conversational responses (e.g., "Okay," "Thanks") are **not** considered out-of-scope.
4.  **Response Format:**
    * "safetyIssue": Indicate with true or false whether the "userMessage" presents a safety issue".
    * "wrongQuery": Indicate with true or false whether the "userMessage" is a wrong query (out of scope).
    * Provide a brief "justification" explaining your classification.

{{ role "user" }}
userMessage: {{userMessage}}

`
export const SafetyPromptInputSchema = z.object({
    userMessage: z.string(),
  });


export const SafetyPromptOutputSchema = z.strictObject({
    wrongQuery: z.boolean().optional().default(false),
    safetyIssue: z.boolean().optional().default(false),
    justification: z.string().optional().default("No justification provided by model")
  });
export type SafetyPromptOutput = z.infer<typeof SafetyPromptOutputSchema>

export const SafetyTransformPrompt = ai.definePrompt(
    {
      name: 'safetyIssuePrompt',
      input: {
        schema: SafetyPromptInputSchema,
      },
      output: {
        schema: SafetyPromptOutputSchema,
        format: 'json',
      },
      config:{
        safetySettings: safetySettings
        }
    },
    
    SafetyIssuePromptText
  );


  
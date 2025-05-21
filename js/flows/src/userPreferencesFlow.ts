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

import { UserPreferenceFlowInputSchema, 
  UserPreferenceFlowOutputSchema, 
  UserPreferenceFlowOutput
 } from './userPreferencesTypes'
import { ai } from './genkitConfig'
import { GenerationBlockedError } from 'genkit';

/**
 * Prompt file: js/flows/prompts/userProfile.prompt
 * 
 * This prompt instructs the LLM to extract user preferences so that we can persist them.
 * 
 * Input schema: UserPreferenceFlowInputSchema
 * Output schema: UserPreferenceFlowOutputSchema
 * 
 * The MovieGuru development team, uses a variant system to version our prompts. The "experimental" variant corresponds to the userPreference.experimental.prompt file. 
 * To use the default variant -- ai.prompt('userPreference')
 * To use a variant -- ai.prompt('userPreference', {variant: 'experimental'})
 * 
 * ATTENTION: We are currently testing an experimental version of the userPreference prompt with 50% of our users. If this is performing well, we should roll it out to 100%.
 */

export const extractUserPreferencesV1 = ai.prompt('userPreference');
export const extractUserPreferencesExperimental = ai.prompt('userPreference', {variant: 'experimental'});

export const UserPreferenceFlow = ai.defineFlow(
  {
    name: 'userPreferenceFlow',
    inputSchema: UserPreferenceFlowInputSchema,

    // Hint Challenge 3: Notice the flow defines an output schema. 
    // Does the information in userPreference.prompt align with this?
    outputSchema: UserPreferenceFlowOutputSchema 
  },
  async (input) => {
    const defaultOutput = UserPreferenceFlowOutputSchema.parse({})
    try {
      var output: UserPreferenceFlowOutput;

      // Using a fairly naive percentage based mechanism to roll out our new version of this prompt
      // Currently targeting 50% of requests
      // TODO: use Firebase Remote Config to make this configurable without a rollout.
      if(Math.random() >= .5) {
        console.info("Routing request to experimental userPreference query.");
        const response = await extractUserPreferencesExperimental({ 
          query: input.query, 
          agentMessage: input.agentMessage });
        output = UserPreferenceFlowOutputSchema.parse(response.output)
      } else {
        console.info("Routing request to default userPreference query.");
        const response = await extractUserPreferencesV1({ 
          query: input.query, 
          agentMessage: input.agentMessage });
        output = UserPreferenceFlowOutputSchema.parse(response.output)
      }
      return output
    } catch (error) {
    
      if(error instanceof GenerationBlockedError){
        console.error("UserPreferenceFlow: GenerationBlockedError generating response:", error.message);
        defaultOutput.safetyIssue = true;
        return defaultOutput;
      }
      else if(error instanceof Error && (error.message.includes('429') || error.message.includes('RESOURCE_EXHAUSTED'))){
        console.error("UserPreferenceFlow: There is a quota issue:", error.message);
        defaultOutput.quotaIssue = true;
        return defaultOutput;
        }
      else{
        console.error("UserPreferenceFlow: Error generating response:", error);
        throw error;
      }
    }
  }
);


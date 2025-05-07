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

import { UserProfileFlowInputSchema, UserProfileFlowOutputSchema, 
 } from './userProfileTypes'
import { ai } from './genkitConfig'
import { GenerationBlockedError } from 'genkit';

/**
 * Prompt file: js/flows/prompts/userProfile.prompt
 * 
 * This prompt instructs the LLM to extract user preferences so that we can persist them.
 * 
 * Input schema: UserProfileFlowInputSchema
 * Output schema: UserProfileFlowOutputSchema
 * 
 * The MovieGuru development team, uses a variant system to version our prompts. The "v2" variant corresponds to the userProfile.v2.prompt file. 
 * To use the default variant -- ai.prompt('userProfile')
 * To use a variant -- ai.prompt('userProfile', {variant: 'v2'})
 */
export const extractUserPreferences = ai.prompt('userProfile', {variant: 'v2'});

export const UserProfileFlow = ai.defineFlow(
  {
    name: 'userProfileFlow',
    inputSchema: UserProfileFlowInputSchema,

    // Hint Challenge 3: Notice the flow defines an output schema. 
    // Does the information in userProfile.prompt align with this?
    outputSchema: UserProfileFlowOutputSchema 
  },
  async (input) => {
    const defaultOutput = UserProfileFlowOutputSchema.parse({})
    try {
      const response = await extractUserPreferences({ 
        query: input.query, 
        agentMessage: input.agentMessage });
      const output = UserProfileFlowOutputSchema.parse(response.output)
      return output
      
    } catch (error) {
    
      if(error instanceof GenerationBlockedError){
        console.error("UserProfileFlow: GenerationBlockedError generating response:", error.message);
        defaultOutput.modelOutputMetadata.safetyIssue = true;
        return defaultOutput;
      }
      else if(error instanceof Error && (error.message.includes('429') || error.message.includes('RESOURCE_EXHAUSTED'))){
        console.error("UserProfileFlow: There is a quota issue:", error.message);
        defaultOutput.modelOutputMetadata.quotaIssue = true;
        return defaultOutput;
        }
      else{
        console.error("UserProfileFlow: Error generating response:", error);
        throw error;
      }
    }
  }
);


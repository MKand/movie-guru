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
import { ai } from './genkitConfig';
import { GenerationBlockedError } from 'genkit';

/**
 * Prompt file: js/flows/prompts/safety.prompt
 * 
 * This prompt instructs the LLM to assess whether a user's statement is safe.
 * 
 * Input schema: SafetyPromptInputSchema
 * Output schema: SafetyPromptOutputSchema
 * 
 * The MovieGuru development team, uses a variant system to version our prompts. The "v2" variant corresponds to the safety.v2.prompt file. 
 * To use the default variant -- ai.prompt('safety')
 * To use a variant -- ai.prompt('safety', {variant: 'v2'})
 */

// get env variable safety_type_conservative which defaults to false
const safety_type_conservative = process.env.SAFETY_CONSERVATIVE || "false"
export var SafetyTransformPrompt = ai.prompt('safety');

if (safety_type_conservative == "true"){
  SafetyTransformPrompt = ai.prompt('safety', {variant : 'conservative'});

}

export const SafetyPromptInputSchema = z.object({
    userMessage: z.string(),
  });

ai.defineSchema('SafetyPromptInputSchema', SafetyPromptInputSchema);

export const SafetyPromptOutputSchema = z.strictObject({
    badQuery: z.boolean().optional().default(false),
    safetyIssue: z.boolean().optional().default(false),
    justification: z.string().default("No justification provided by model")
  });
export type SafetyPromptOutput = z.infer<typeof SafetyPromptOutputSchema>;

ai.defineSchema('SafetyPromptOutputSchema', SafetyPromptOutputSchema);

export const SafetyIssueFlow = ai.defineFlow(
    { 
        name: 'safetyIssueFlow',
        inputSchema: SafetyPromptInputSchema,
        outputSchema: SafetyPromptOutputSchema,
    },
    async (input) => {
        const defaultOutput = SafetyPromptOutputSchema.parse({});
        try{
            const response = await SafetyTransformPrompt(input);
            const safeOutput = response.output?? defaultOutput;
            const output = SafetyPromptOutputSchema.parse(safeOutput)
            return output;
            } catch (error) {
            if (error instanceof GenerationBlockedError){
                console.error("QTFlow: GenerationBlockedError generating response:", error.message);
                defaultOutput.safetyIssue = true
                return defaultOutput;
            }
            else if(error instanceof Error && (error.message.includes('429') || error.message.includes('RESOURCE_EXHAUSTED'))){
                console.error("QTFlow: There is a quota issue:", error.message);
                return defaultOutput;
                }
                else {
                console.error("QTFlow: Error generating response:", error);
                throw error;
              }
            }
        }
)

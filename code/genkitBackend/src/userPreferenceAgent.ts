import { ai } from './genkitConfig';
import { z } from 'genkit';
import { PreferenceItemSchema, PreferenceItem, UserPreferenceInputSchema, UserPreferenceOutputSchema } from './userPreferencesTypes';
import { UserPreferencesDB } from './db';
import { userPreferencePromptText } from './prompts';
import {  vertexAI } from '@genkit-ai/vertexai';

const userPreferencesDB = new UserPreferencesDB();
userPreferencesDB.init();

export const userPreferenceLookupTool = ai.defineTool(
  {
    name: 'userPreferenceLookupTool',
    description: "use this tool to look up the user's preferences",
    inputSchema: z.object({ userName: z.string().describe('the name of the user') }),
    outputSchema: z.array(PreferenceItemSchema).describe('the profile items for that user'),
  },
  async (input) => {
    console.log(`Looking up preferences for user: ${input.userName} at time ${new Date().toISOString()}`);
    return await userPreferencesDB.get(input.userName);
  },
);

export const userPreferenceUpdateTool = ai.defineTool(
  {
    name: 'userPreferenceUpdateTool',
    description: "use this tool to update the user's preferences with new items",
    inputSchema: z.object({ userId: z.string(), changes: z.array(PreferenceItemSchema) }),
    outputSchema: z.boolean(),
  },
  async (input) => {
    console.log(`Updating preferences for user: ${input.userId} with changes: ${JSON.stringify(input.changes)} at time ${new Date().toISOString()}`);
    const existingProfile = await userPreferencesDB.get(input.userId);
    const profileMap = new Map<string, PreferenceItem>();
    for (const item of existingProfile) {
      profileMap.set(item.item, item);
    }
    for (const item of input.changes) {
      profileMap.set(item.item, item);
    }
    const newProfile = Array.from(profileMap.values());
    await userPreferencesDB.update(input.userId, newProfile);
    return true;
  },
);

const userPreferencePrompt = ai.definePrompt(
  {
    name: 'userPreferencePrompt',
    tools: [userPreferenceLookupTool, userPreferenceUpdateTool],
    input: {
      schema: UserPreferenceInputSchema,
    },
    output: {
      schema: UserPreferenceOutputSchema,
    },

  },
  userPreferencePromptText
);

export const UserPreferenceFlow = ai.defineFlow(
  {
    name: 'userPreferenceFlow',
    inputSchema: UserPreferenceInputSchema,
    outputSchema: UserPreferenceOutputSchema,
  },
  async (input) => {
    const defaultOutput = UserPreferenceOutputSchema.parse({})
    const response = await userPreferencePrompt(input);
    const safeOutput = response.output?? defaultOutput
    const output = UserPreferenceOutputSchema.parse(safeOutput);
    return output;
  }
);
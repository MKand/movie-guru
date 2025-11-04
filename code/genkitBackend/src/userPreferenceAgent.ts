import {ai} from './genkitConfig'
import { z } from 'genkit';
import { PreferenceItemSchema, UserPreferenceInputSchema, UserPreferenceOutputSchema } from './userPreferencesTypes';
import { UserPreferencesDB } from './db';
import { userPreferencePrompt } from './prompts';

const userPreferencesDB = new UserPreferencesDB();
userPreferencesDB.init();

export const userPreferenceLookupTool = ai.defineTool(
  {
    name: 'userPreferenceLookupTool',
    description: "use this tool to look up the user's preferences",
    inputSchema: z.string().describe('the name of the user'),
    outputSchema: z.array(PreferenceItemSchema).describe('the profile items for that user'),
  },
  async (input) => {
    return await userPreferencesDB.get(input);
  },
);

export const userPreferenceUpdateTool = ai.defineTool(
  {
    name: 'userPreferenceUpdateTool',
    description: "use this tool to update the user's preferences",
    inputSchema: z.object({userId: z.string(), profile: z.array(PreferenceItemSchema)}),
    outputSchema: z.void(),
  },
  async (input) => {
    await userPreferencesDB.update(input.userId, input.profile);
  },
);

export const userPreferenceAgent = ai.definePrompt(
  {
    name: 'userPreferenceAgent',
    tools: [userPreferenceLookupTool, userPreferenceUpdateTool],
    input: {
      schema: UserPreferenceInputSchema
    },
    output: {
      schema: UserPreferenceOutputSchema
    },
    prompt: userPreferencePrompt,

  }
)
import { ai } from './genkitConfig';
import { z } from 'genkit';
import { ProfileChangeRecommendationSchema, ProfileChangeRecommendation } from './userPreferencesTypes';
import { UserPreferencesDB } from './preferenceDb';

const userPreferencesDB = new UserPreferencesDB();
userPreferencesDB.init();

export const UserProfileInputSchema = z.object({
  userName: z.string(),
  query: z.string(),
});

ai.defineSchema('UserProfileInputSchema', UserProfileInputSchema);


export const UserProfileOutputSchema = z.object({
  justification: z.string().default("No justification provided"),
  updatedPreferences: z.boolean().default(false)
});

ai.defineSchema('UserProfileOutputSchema', UserProfileOutputSchema);


export const userPreferenceLookupTool = ai.defineTool(
  {
    name: 'userPreferenceLookupTool',
    description: "use this tool to look up the user's preferences",
    inputSchema:  z.object({ userName: z.string().describe('the name of the user') }),
    outputSchema: z.array(ProfileChangeRecommendationSchema).describe('the profile items for that user'),
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
    // Accept userName` to be robust against prompt/template differences.
    inputSchema: z.object({ userName: z.string().describe('the name of the user'), changes: z.array(ProfileChangeRecommendationSchema) }),
    outputSchema: z.boolean(),
  },
  async (input) => {
    
    console.log(`Updating preferences for user: ${input.userName} with changes: ${JSON.stringify(input.changes)} at time ${new Date().toISOString()}`);
    const existingProfile = await userPreferencesDB.get(input.userName);
    const profileMap = new Map<string, ProfileChangeRecommendation>();
    for (const item of existingProfile) {
      profileMap.set(item.item, item);
    }
    for (const item of input.changes) {
      profileMap.set(item.item, item);
    }
    const newProfile = Array.from(profileMap.values());
    await userPreferencesDB.update(input.userName, newProfile);
    return true;
  },
);


export const userPreferenceUpdatePrompt = ai.definePrompt(
  {
    name: 'userPreferenceUpdatePrompt',
    tools: [userPreferenceLookupTool, userPreferenceUpdateTool],
    maxTurns: 8,
    input: {
      schema: UserProfileInputSchema
    },
    output: {
        schema: UserProfileOutputSchema,
        format: "json"
    }
    // Don't specify output schema here - it prevents tool calling
    // The flow will parse the text response into the schema

  },
  `
    {{ role "system" }}

    You are an agent that manages user movie preferences. Your ONLY job is to call tools to manage preferences.

    1.  Analyze the user's message for new, strong, long-term preferences.
    2.  If the user expresses any strong preferences, use the "userPreferenceLookupTool" to get the user's current preferences.
    3.  If you find a new preference that is not in the current preferences list, you MUST call "userPreferenceUpdateTool" with ONLY the new preferences.
    4.  After calling the tools, respond with the following:
         "justification": "what you did and why",
         "updatedPreferences": true or false (whether you called the "userPreferenceUpdateTool" tool)
       
    Do not describe the changes in your text response, perform them with the tool calls first.

    {{ role "user" }}
    here is the input
    * userMessage: {{query}}
    * userName: {{userName}}
    `
);

export const UserPreferenceUpdateFlow = ai.defineFlow(
  {
    name: 'userPreferenceUpdateFlow',
    inputSchema: UserProfileInputSchema,
    outputSchema: UserProfileOutputSchema 
  },
  async (input) => {
    const defaultOutput = UserProfileOutputSchema.parse({})
    const response = await userPreferenceUpdatePrompt(input);

    // Parse the text response as JSON
    const safeOutput = response.output ?? defaultOutput;
    const output = UserProfileOutputSchema.parse(safeOutput);
    return output;
  }
);
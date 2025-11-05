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


export const UserProfileOutputSchema = z.strictObject({
  justification: z.string().default("No justification provided"),
  updatedPreferences: z.boolean().default(false).describe('whether the user preferences were updated by calling the update tool'),
});

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
    inputSchema: z.object({ userId: z.string(), changes: z.array(ProfileChangeRecommendationSchema) }),
    outputSchema: z.boolean(),
  },
  async (input) => {
    console.log(`Updating preferences for user: ${input.userId} with changes: ${JSON.stringify(input.changes)} at time ${new Date().toISOString()}`);
    const existingProfile = await userPreferencesDB.get(input.userId);
    const profileMap = new Map<string, ProfileChangeRecommendation>();
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


const userPreferenceUpdatePrompt = ai.definePrompt(
  {
    name: 'userPreferenceUpdatePrompt',
    tools: [userPreferenceLookupTool, userPreferenceUpdateTool],
    input: {
      schema: UserProfileInputSchema
    },
    output: {
      schema: UserProfileOutputSchema
    },

  },
  `
    {{ role "system" }}

    You are an agent that manages user movie preferences. Your ONLY job is to call tools to manage preferences.

    1.  Call "userPreferenceLookupTool" to get the user's current preferences.
    2.  Analyze the user's message for new, strong, long-term preferences.
    3.  If you find a new preference that is not in the existing list, you MUST call "userPreferenceUpdateTool" with ONLY the new preferences.
    4.  Your final output MUST be a justification of what you did, and whether you called the update tool. You MUST NOT say that you have updated the preferences if you have not called the tool.

    Do not describe the changes in your final response, perform them with the tool.

    {{ role "user" }}

    * user message: {{query}}
    * user name: {{userName}}
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
    const safeOutput = response.output?? defaultOutput
    const output = UserProfileOutputSchema.parse(safeOutput);
    return output;
  }
);
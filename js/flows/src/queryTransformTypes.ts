import { z } from 'genkit';
import { ModelOutputMetadata, ModelOutputMetadataSchema } from './modelOutputMetadataTypes';

// USERINTENT as Zod Enum
export const USERINTENT = z.enum([
  'UNCLEAR',
  'GREET',
  'END_CONVERSATION',
  'REQUEST',
  'RESPONSE',
  'ACKNOWLEDGE',
]);

// ProfileCategories schema
const ProfileCategoriesSchema = z.object({
  actors: z.array(z.string()).optional(),
  directors: z.array(z.string()).optional(),
  genres: z.array(z.string()).optional(),
  others: z.array(z.string()).optional(),
});

export type ProfileCategories = z.infer<typeof ProfileCategoriesSchema>


// UserProfile schema
export const UserProfileSchema = z.object({
  likes: ProfileCategoriesSchema.optional(),
  dislikes: ProfileCategoriesSchema.optional(),
});

export type UserProfile = z.infer<typeof UserProfileSchema>


// SimpleMessage schema
export const SimpleMessageSchema = z.object({
  role: z.string(), 
  content: z.string(),
});

export type SimpleMessage = z.infer<typeof SimpleMessageSchema>


// QueryTransformFlowInput schema
export const QueryTransformFlowInputSchema = z.object({
  history: z.array(SimpleMessageSchema),
  userProfile: UserProfileSchema.optional(),
  userMessage: z.string(),
});

export type QueryTransformFlowInput = z.infer<typeof QueryTransformFlowInputSchema>


// QueryTransformFlowOutput schema
export const QueryTransformFlowOutputSchema = z.object({
  transformedQuery: z.string(),
  userIntent: z.string(),
  modelOutputMetadata: ModelOutputMetadataSchema,
});

export type QueryTransformFlowOutput = z.infer<typeof QueryTransformFlowOutputSchema>

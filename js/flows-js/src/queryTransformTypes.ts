import { z } from 'zod';

// USERINTENT as Zod Enum
const USERINTENT = z.enum([
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

// UserProfile schema
export const UserProfileSchema = z.object({
  likes: ProfileCategoriesSchema.optional(),
  dislikes: ProfileCategoriesSchema.optional(),
});

// SimpleMessage schema
export const SimpleMessageSchema = z.object({
  sender: z.string(), 
  message: z.string(),
});

// QueryTransformFlowInput schema
export const QueryTransformFlowInputSchema = z.object({
  history: z.array(SimpleMessageSchema),
  userProfile: UserProfileSchema,
  userMessage: z.string(),
});

// QueryTransformFlowOutput schema
export const QueryTransformFlowOutputSchema = z.object({
  transformedQuery: z.string().optional(),
  userIntent: USERINTENT.optional(),
  justification: z.string().optional(),

});
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
  actors: z.array(z.string()),
  directors: z.array(z.string()),
  genres: z.array(z.string()),
  others: z.array(z.string()),
});

// UserProfile schema
export const UserProfileSchema = z.object({
  likes: ProfileCategoriesSchema,
  dislikes: ProfileCategoriesSchema,
});

// SimpleMessage schema
export const SimpleMessageSchema = z.object({
  role: z.string(), // Changed 'sender' to 'role' to match the type definition
  content: z.string(),
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
import * as z from 'zod';

// Enums as Zod Enums
const MovieFeatureCategory = z.enum(['OTHER', 'ACTOR', 'DIRECTOR', 'GENRE']);
const Sentiment = z.enum(['POSITIVE', 'NEGATIVE']);


// ProfileChangeRecommendation schema
const ProfileChangeRecommendationSchema = z.object({
  item: z.string(),
  reason: z.string(),
  category: MovieFeatureCategory,
  sentiment: Sentiment,
});

// UserProfileFlowInput schema
export const UserProfileFlowInputSchema = z.object({
  query: z.string(),
  agentMessage: z.string(),
});

// UserProfileFlowOutput schema
export const UserProfileFlowOutputSchema = z.object({
  profileChangeRecommendations: z.array(ProfileChangeRecommendationSchema),
  justification: z.string().optional(),
});

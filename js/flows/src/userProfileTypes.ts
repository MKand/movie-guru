import { z } from 'genkit';
import { ModelOutputMetadata } from './modelOutputMetadataTypes';

// Enums as Zod Enums
const MovieFeatureCategory = z.enum(['OTHER', 'ACTOR', 'DIRECTOR', 'GENRE']);
const Sentiment = z.enum(['POSITIVE', 'NEGATIVE']);

// ProfileChangeRecommendation schema
export const ProfileChangeRecommendationSchema = z.object({
  item: z.string(),
  reason: z.string(),
  category: MovieFeatureCategory,
  sentiment: Sentiment,
});

export type ProfileChangeRecommendation = z.infer<typeof ProfileChangeRecommendationSchema>

// UserProfileFlowInput schema
export const UserProfileFlowInputSchema = z.object({
  query: z.string(),
  agentMessage: z.string(),
});

export type UserProfileFlowInput = z.infer<typeof UserProfileFlowInputSchema>

// UserProfileFlowOutput schema
export const UserProfileFlowOutputSchema = z.object({
  profileChangeRecommendations: z.array(ProfileChangeRecommendationSchema),
  modelOutputMetadata: ModelOutputMetadata
});

export type UserProfileFlowOutput = z.infer<typeof UserProfileFlowOutputSchema>

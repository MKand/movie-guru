import * as z from 'zod';


// Enums as Zod Enums
const MovieFeatureCategory = z.enum(['OTHER', 'ACTOR', 'DIRECTOR', 'GENRE']);
const Sentiment = z.enum(['POSITIVE', 'NEGATIVE']);

// ProfileCategories schema
const ProfileCategoriesSchema = z.object({
  actors: z.array(z.string()).optional(),
  directors: z.array(z.string()).optional().transform(obj => obj ?? []),
  genres: z.array(z.string()).optional(),
  others: z.array(z.string()).optional(),
});

// UserProfile schema
const UserProfileSchema = z.object({
  likes: ProfileCategoriesSchema.optional(),
  dislikes: ProfileCategoriesSchema.optional(),
});

// ProfileChangeRecommendation schema
const ProfileChangeRecommendationSchema = z.object({
  item: z.string(),
  reason: z.string(),
  category: MovieFeatureCategory,
  sentiment: Sentiment,
});

// UserProfileFlowInput schema
const UserProfileFlowInputSchema = z.object({
  query: z.string(),
  agentMessage: z.string(),
});

// UserProfileFlowOutput schema
const UserProfileFlowOutputSchema = z.object({
  profileChangeRecommendations: z.array(ProfileChangeRecommendationSchema),
  justification: z.string().optional(),
});

// Function to create a default UserProfileFlowOutput (optional)
function NewUserProfileFlowOuput() {
  return UserProfileFlowOutputSchema.parse({
    profileChangeRecommendations: Array(5).fill(null).map(() => ({ 
      item: '', 
      reason: '', 
      category: 'OTHER', 
      sentiment: 'POSITIVE' 
    })), 
  });
}

export {
  UserProfileFlowInputSchema,
  UserProfileFlowOutputSchema,
  NewUserProfileFlowOuput,
  ProfileChangeRecommendationSchema,
  MovieFeatureCategory,
  Sentiment,
};

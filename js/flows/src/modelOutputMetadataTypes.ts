import { z } from 'genkit';

export const ModelOutputMetadataSchema = z.object({
    justification: z.string(),
    safetyIssue: z.boolean()
}) 

export type ModelOutputMetadata = z.infer<typeof ModelOutputMetadataSchema>;

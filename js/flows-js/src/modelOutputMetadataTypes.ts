import * as z from 'zod';

export const ModelOutputMetadata = z.object({
    justification: z.string(),
    safetyIssue: z.boolean()
}) 

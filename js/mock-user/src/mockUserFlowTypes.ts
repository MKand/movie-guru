import { z } from 'zod';

export const MockUserFlowInputSchema = z.object({
  expert_answer: z.string(),
  response_mood: z.string(),
  response_type: z.string(),
});

export const MockUserFlowOutputSchema = z.object({
  answer: z.string(),
});
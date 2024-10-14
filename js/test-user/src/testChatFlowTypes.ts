import { z } from 'zod';

const RESPONSE_MOOD = z.enum([
  'POSITIVE',
  'NEGATIVE',
  'NEUTRAL',
  'RANDOM',
]);

const RESPONSE_TYPE = z.enum([
  'DIVE_DEEP',
  'CHANGE_TOPIC',
  'END_CONVERSATION',
  'CONTINUE',
  'RANDOM'
]);

export const DummyUserFlowInputSchema = z.object({
  expert_answer: z.string(),
  response_mood: z.string(),
  response_type: z.string(),
});

export const DummyUserFlowOutputSchema = z.object({
  answer: z.string(),
});
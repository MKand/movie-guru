
import { z } from 'genkit';
import { ModelOutputMetadata, ModelOutputMetadataSchema } from './modelOutputMetadataTypes';

const SearchTypeCategory = z.enum(['KEYWORD', 'VECTOR', 'MIXED', 'NONE']);


export const RetrieverOptionsSchema = z.object({
  k: z.number().optional().default(10),
  searchCategory: SearchTypeCategory.optional().default("VECTOR"),
  keywordQuery: z.string().default(""),
  vectorQuery: z.string().default(""),

});

export const QuerySchema = z.object({
  query: z.string(),
});

export const SearchFlowOutputSchema = z.object({
  keywordQuery: z.string().optional(),
  vectorQuery: z.string().optional(),
  searchCategory: SearchTypeCategory,
  modelOutputMetadata: ModelOutputMetadataSchema,
});

export type SearchFlowOutput = z.infer<typeof SearchFlowOutputSchema>
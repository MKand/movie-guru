module.exports = {
    evaluators: [
      {
        flowName: 'RAGFlow',
        extractors: {
          input: {outputOf: 'QueryTransformFlowPrompt'},
          context: { outputOf: 'MixedRetriever' },
          output: 'MovieFlowPrompt',
        },
      },
    ],
  };
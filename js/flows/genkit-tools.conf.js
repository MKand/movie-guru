module.exports = {
    evaluators: [
      {
        flowName: 'RAGFlow',
        extractors: {
          context: { outputOf: 'MixedRetriever' },
          output: 'MovieFlowPrompt',
        },
      },
    ],
  };
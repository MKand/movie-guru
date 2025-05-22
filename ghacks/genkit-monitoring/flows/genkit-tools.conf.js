module.exports = {
    evaluators: [
      {
        actionRef: '/flow/chatFlow',
        extractors: {
          context: { outputOf: 'movieDocFlow' },
        },
      },
    ],
  };
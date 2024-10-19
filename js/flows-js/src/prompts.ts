export const QueryTransformPromptText = 
`
Here are the inputs:
* userProfile: (May be empty)
    * likes: 
        * actors: {{#each userProfile.likes.actors}}{{this}}, {{~/each}}
        * directors: {{#each userProfile.likes.directors}}{{this}}, {{~/each}}
        * genres: {{#each userProfile.likes.genres}}{{this}}, {{~/each}}
        * others: {{#each userProfile.likes.others}}{{this}}, {{~/each}}
    * dislikes: 
        * actors: {{#each userProfile.dislikes.actors}}{{this}}, {{~/each}}
        * directors: {{#each userProfile.dislikes.directors}}{{this}}, {{~/each}}
        * genres: {{#each userProfile.dislikes.genres}}{{this}}, {{~/each}}
        * others: {{#each userProfile.dislikes.others}}{{this}}, {{~/each}}
* userMessage: {{userMessage}}
* history: (May be empty)
    {{#each history}}{{this.role}}: {{this.content}}{{~/each}}
`
export const MovieFlowPromptText = 
` 
    Here are the inputs:
    * userPreferences: (May be empty)
    * userMessage: {{userMessage}}
    * history: (May be empty)
    * Context retrieved from vector db (May be empty):
    `

export const UserProfilePromptText = 
	` 
Here are the inputs:
	1. Optional Message 0 from agent: {{agentMessage}}
	2. Required Message 1 from user: {{query}}
`
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
    {{#each history}}{{this.sender}}: {{this.message}}{{~/each}}
`
export const MovieFlowPromptText = 
` 
Here are the inputs:
* userPreferences: (May be empty)
    * likes: 
        * actors: {{#each userPreferences.likes.actors}}{{this}}, {{~/each}}
        * directors: {{#each userPreferences.likes.directors}}{{this}}, {{~/each}}
        * genres: {{#each userPreferences.likes.genres}}{{this}}, {{~/each}}
        * others: {{#each userPreferences.likes.others}}{{this}}, {{~/each}}
    * dislikes: 
        * actors: {{#each userPreferences.dislikes.actors}}{{this}}, {{~/each}}
        * directors: {{#each userPreferences.dislikes.directors}}{{this}}, {{~/each}}
        * genres: {{#each userPreferences.dislikes.genres}}{{this}}, {{~/each}}
        * others: {{#each userPreferences.dislikes.others}}{{this}}, {{~/each}}
* userMessage: {{userMessage}}
* history: (May be empty)
    {{#each history}}{{this.sender}}: {{this.message}}{{~/each}}
* Context retrieved from vector db (May be empty):
{{#each contextDocuments}} 
Movie: 
- title:{{this.title}}
- plot:{{this.plot}} 
- genres:{{this.genres}}
- actors:{{this.actors}} 
- directors:{{this.directors}} 
- rating:{{this.rating}} 
- runtimeMinutes:{{this.runtimeMinutes}}
- released:{{this.released}} 
{{/each}}
`
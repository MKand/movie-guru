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
* Context retrieved from vector db (May be empty):

	{{#each contextDocuments}} 
	Movie: 
	- title:{{this.title}}
	- runtimeMinutes:{{this.runtimeMinutes}}
	- genres:{{this.genres}}
	- actors:{{this.actors}} 
	- directors:{{this.directors}} 
	- released:{{this.released}} 
	- plot:{{this.plot}} 
	- rating:{{this.rating}} 
	{{/each}}

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
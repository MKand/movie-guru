export const UserProfilePromptText = 
	` Inputs: 
	1. Optional Message 0 from agent: {{agentMessage}}
	2. Required Message 1 from user: {{query}}
`
export const QueryTransformPromptText = `
  Here are the inputs:
		* Conversation History (this may be empty):
			{{history}}
		* UserProfile (this may be empty):
			{{userProfile}}
		* User Message:
			{{userMessage}})
    `
export const MovieFlowPromptText =  ` 
	Here are the inputs:

	* Context retrieved from vector db:
	{{contextDocuments}}

	* User Preferences:
	{{userPreferences}}

	* Conversation history:
	{{history}}

	* User message:
	{{userMessage}}

`
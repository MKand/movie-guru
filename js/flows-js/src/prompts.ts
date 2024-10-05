export const UserProfilePromptText = ` 
		Optional Message 0 from agent: {{agentMessage}}
		Required Message 1 from user: {{query}}
		Just say hi in a language you know.
      `
export const QueryTransformPromptText = `
	Here are the inputs:
	* Conversation History (this may be empty):
		{{history}}
	* UserProfile (this may be empty):
		{{userProfile}}
	* User Message:
		{{userMessage}})
		Translate the user's message into a random language.
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

	Translate the user's message into a random language.
`
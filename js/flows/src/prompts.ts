/**
 * Copyright 2025 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export const UserProfilePromptText = ` 
    {{ role "system" }}
    You are a user's movie profiling expert focused on uncovering users' enduring likes and dislikes. 
       Your task is to analyze the user message and extract ONLY strongly expressed, enduring likes and dislikes related to movies.
       Once you extract any new likes or dislikes from the current query respond with the items you extracted with:
            1. the category (ACTOR, DIRECTOR, GENRE, OTHER)
            2. the item value
            3. your reason behind the choice
            4. the sentiment of the user has about the item (POSITIVE, NEGATIVE).
          
        Guidelines:
        1. There are some examples given below. Do not take the examples literally, but use them to create a general method for parsing the input and constructing the right output.
        2. Strong likes and dislikes Only: Add or Remove ONLY items expressed with strong language indicating long-term enjoyment or aversion (e.g., "love," "hate," "can't stand,", "always enjoy"). Ignore mild or neutral items (e.g., "like,", "okay with," "fine", "in the mood for", "do not feel like").
        3. Distinguish current state of mind vs. Enduring likes and dislikes:  Be very cautious when interpreting statements. Focus only on long-term likes or dislikes while ignoring current state of mind. If the user expresses wanting to watch a specific type of movie or actor NOW, do NOT assume it's an enduring like unless they explicitly state it. For example, "I want to watch a horror movie movie with Christina Appelgate" is a current desire, NOT an enduring preference for horror movies or Christina Appelgate.
        4. Focus on Specifics:  Look for concrete details about genres, directors, actors, plots, or other movie aspects.
        5. Give an explanation as to why you made the choice.

        
      Respond with the following:
  
          *   a *justification* about why you created the query this way.
          *   a *safetyIssue* returned as true if the query is considered dangerous. A query is considered dangerous if the user is asking you to tell about something dangerous. However, asking for movies with dangerous themes is not considered dangerous.
          *   a list of *profileChangeRecommendations* that are a list of extracted strong likes or dislikes with the following fields: category, item, reason, sentiment
      
      {{ role "user" }}
        *Optional Message 0 from agent: {{agentMessage}}
       * user message: {{query}}
      `



export const DocSearchFlowPromptText = `
{{ role "system" }}

Analyze the inputQuery string with respect to a movie database containing the following fields:

embedding: Vector representation of the movie's title, plot, and genres.
genres: List of movie genres (e.g., "Action", "Comedy", "Drama", "horror", "anime", "classic", "romance", "romcom" etc).
title: Title of the movie.
plot: A textual summary of the movie's plot.
runtime_mins: Duration of the movie in minutes.
released: Release year of the movie.
actors: List of actors in the movie.
director: Director of the movie.
rating: Numerical rating from 1 to 5.

Task:

Determine the appropriate search category for the query: KEYWORD, VECTOR, or MIXED.

1.  KEYWORD search: Use when the query can be expressed with SQL operators for the postgres db (e.g., =, !=, >, <, IN, ILIKE, ANY) on the title, actors, director, genres, runtime_mins, released, or rating fields.
    * Queries about movie quality, length, or release year, or specific actors or directors, or exact titles can often be transformed for KEYWORD search.
    * If the query contains text searches, make them case-insensitive using ILIKE.
    * When the search is based on titles or names, *always* use the ILIKE operator to circumvent potential spelling and punctuation mismatches. Use the format ILIKE "'%<search_term>%'".
    * **IMPORTANT**: Remove all special characters (e.g., ', @, #, $,%, etc.) from the names or titles before using ILIKE, leaving only letters, numbers, and spaces. Remove any apostrophes in titles
    * **Crucially, when the query asks for a specific detail *of* a movie (e.g., "director of X," "plot of Y"), prioritize identifying the movie using "title ILIKE '%<movietitle>%'" as the primary KEYWORD query. Then, if necessary, add further filtering for other details (e.g., actor, director).**
    * When searching for a director or actor in general (e.g., "movies with Tom Hanks"), ensure the "ILIKE ANY" operator is used correctly to match against the "director" or "actors" fields, respectively. The search term should be enclosed in single quotes.
    * There are some examples given below. Do not take the examples literally, but use them to create a general method for parsing the input and constructing the right output.

    Transformations:
        Movie Quality:
            Bad: rating < 2
            Average: rating BETWEEN 2 AND 3.5
            Good: rating > 3.5
            Great: rating > 4.5
            Terrible: rating < 1
        Movie Length:
            Short: runtime_mins < 45
            Long: runtime_mins > 120
            Very Long: runtime_mins > 150
        Movie Year:
            Recent: released > 2020
            Old: released < 2005

    Examples of transformed KEYWORD queries:
    * IMPORTANT: Remove all special characters (e.g.: @, #, $,%,', etc.) from the names or titles before using ILIKE, leaving only letters, numbers, and spaces. Pay special care that you remove apostrophes from titles.

    -   Input: "great movie that is short"
        Output:
            searchCategory: KEYWORD
            KeywordQuery: "rating > 4.5 AND runtime_mins < 45"
            VectorQuery: ""
            justification: "The query can be fully expressed using comparisons on the 'rating' (transformed to 'great') and 'runtime_mins' (transformed to 'short') fields."

    -   Input: "movies released after 2000"
        Output:
            searchCategory: KEYWORD
            KeywordQuery: "released > 2000"
            VectorQuery: ""
            justification: "The query filters movies based on the 'released' year field."

    -   Input: "The Bee movie" OR "title The Bee movie" OR "movie The Bee movie" (examples of looking for a specific movie)
        Output:
            searchCategory: KEYWORD
            KeywordQuery: "title ILIKE '%The Bee movie%'"
            VectorQuery: ""
            justification: "The query is searching for a movie by its title, using ILIKE for case-insensitive and fuzzy matching."

    -   Input: "movies with tom hanks"
        Output:
            searchCategory: KEYWORD
            KeywordQuery: "'Tom Hanks' ILIKE ANY(string_to_array(actors, ', '))"
            VectorQuery: ""
            justification: "The query filters movies based on whether 'Tom Hanks' is an actor in the 'actors' list."

    -   Input: "director of The Matrix"
        Output:
            searchCategory: KEYWORD
            KeywordQuery: "title ILIKE '%The Matrix%'"
            VectorQuery: ""
            justification: "The query asks for a detail ('director') of a specific movie ('The Matrix'). We prioritize finding the movie using 'title ILIKE'."

    -   Input: "movies like The Matrix released after 2000"
        Output:
            searchCategory: MIXED
            KeywordQuery: "released > 2000"
            VectorQuery: "movies like The Matrix"
            justification: "The query asks for movies similar to the Matrix and includes a KEYWORD filter ('released after 2000')."

    -   Input: "Is the rating of The Mummy's Revenge good?"
        Output:
            searchCategory: KEYWORD
            KeywordQuery: "title ILIKE '%The Mummys Revenge%'" (removed apostrophe from title)
            VectorQuery: ""
            justification: "The query asks about a specific detail ('rating') of a specific movie ('The Mummmys Revenge'). I removed the apostrophe from the title. We prioritize finding the movie using 'title ILIKE'."
     -  Input: "What is the plot of The Matrix?"
        Output:
            searchCategory: KEYWORD
            KeywordQuery: "title ILIKE '%The Matrix%'"
            VectorQuery: ""
            justification: "The query asks about a specific detail ('plot') of a specific movie ('The Matrix'). We prioritize finding the movie using 'title ILIKE'."

    -   Input: "is Jane Doe in The Bee Movie"
        Output:
            searchCategory: KEYWORD
            KeywordQuery: "title ILIKE '%The Bee Movie%'"
            VectorQuery: ""
            justification: "The query asks about a specific detail ('is Jane Doe in') of a specific movie ('The Bee Movie'). We prioritize finding the movie using 'title ILIKE'."

2.  VECTOR search: Use when the query requires semantic understanding of title, plot, or genres. Applicable for queries involving concepts, emotions, themes, or natural language descriptions.
    * Searches that primarily involve genres should often have a vector query, especially if combined with other criteria.
    * If a query is vague without any specific descriptors (e.g., "movies", "something to watch"), make it a vector query.

    Examples of VECTOR queries:

    -   Input: "movies with strong female leads"
        Output:
            searchCategory: VECTOR
            KeywordQuery: ""
            VectorQuery: "strong female leads"
            justification: "The query requires semantic understanding of 'strong female leads', which cannot be expressed with KEYWORD operators."

    -   Input: "movies" OR "some movie" (examples of vague non-specific query)
        Output:
            searchCategory: VECTOR
            KeywordQuery: ""
            VectorQuery: "movies"
            justification: "The query is vague and requires semantic understanding for movie recommendations."

    -   Input: "movies like The Matrix"
        Output:
            searchCategory: VECTOR
            KeywordQuery: ""
            VectorQuery: "like The Matrix"
            justification: "The query requires semantic understanding of similarity, which cannot be expressed with KEYWORD operators."

    -   Input: "romantic films"
        Output:
            searchCategory: VECTOR
            KeywordQuery: ""
            VectorQuery: "romance"
            justification: "The query primarily focuses on the genre 'romance', which is best handled with a VECTOR search."

    -   Input: "movies with anime"
        Output:
            searchCategory: VECTOR
            KeywordQuery: ""
            VectorQuery: "anime"
            justification: "The query primarily focuses on the genre 'anime', which is best handled with a VECTOR search."

3.  MIXED search: Queries that require both KEYWORD and VECTOR search: Use when part of the query relates to structured fields (KEYWORD search), while another part involves semantic understanding (VECTOR search).

    Examples of MIXED queries:

    -   Input: "fun movies released after 2004"
        Output:
            searchCategory: MIXED
            KeywordQuery: "released > 2004"
            VectorQuery: "fun movies"
            justification: "The query combines a KEYWORD filter on 'released' with a VECTOR search for 'fun movies'."

    -   Input: "horror movies with great ratings with Tom Hanks"
        Output:
            searchCategory: MIXED
            KeywordQuery: "rating > 4.5 AND 'Tom Hanks' ILIKE ANY(string_to_array(actors, ', '))"
            VectorQuery: "horror"
            justification: "The query combines KEYWORD filters on 'rating' and 'actors' with a VECTOR search for 'horror movies'."



    Respond with the following:
        keywordQuery: A concise representation of the query, empty if needed.
        vectorQuery: A concise representation of the query, empty if needed.
        searchCategory: The determined category: KEYWORD, VECTOR, or MIXED.
        justification: Explanation of the classification, referencing specific fields or transformations where applicable.
        {{ role "user" }}
         {{ query }}
    `

export const ConversationQualityAnalysisPromptText = 
		`
        {{ role "system" }}
		You are an AI assistant designed to analyze conversations between users and a movie expert agent. 
		Your task is to objectively assess the flow of the conversation and determine the outcome of the agent's response based solely on the user's reaction to it.
		You also need to determine the user's sentiment based on their last message (it can be positive, negative, neutral, or ambiguous).
		You only get a truncated version of the conversation history.
        There are some examples given below. Do not take the examples literally, but use them to create a general method for parsing the input and constructing the right output.

		Here's how to analyze the conversation:

		1. Read the conversation history carefully, paying attention to the sequence of messages and the topics discussed.
		2. Focus on the agent's response and how the user reacts to it.

		Guidelines for classification of the conversation outcome:

		*   OUTCOMEIRRELEVANT: The agent's response is not connected to the user's previous turn or doesn't address the user's query or request.
		*   OUTCOMEACKNOWLEDGED: The user acknowledges the agent's response with neutral remarks like "Okay," "Got it," or a simple "Thanks" without indicating further interest or engagement.
		*   OUTCOMEREJECTED: The user responds negatively to the agent's response like "No," "I don't like it," or a simple "No thanks" without indicating further interest or engagement.
		*   OUTCOMEENGAGED: The user shows interest in the agent's response and wants to delve deeper into the topic. This could be through follow-up questions, requests for more details, or expressing a desire to learn more about the movie or topic mentioned by the agent.
		*   OUTCOMETOPICCHANGE: The user shifts the conversation to a new topic unrelated to the agent's response.
		*   OUTCOMEAMBIGUOUS: The user's response is too vague or open-ended to determine the outcome with certainty.

		Examples:

		User: "I'm looking for a movie with strong female characters."
		Agent: "Have you seen 'Alien'?"
		User: "Tell me more about it."
		Outcome: OUTCOMEENGAGED (The user shows interest in the agent's suggestion and wants to learn more.)

		Agent: "Let me tell you about the movie 'Alien'?"
		User: "I hate that film"
		Outcome: OUTCOMEREJECTED (The user rejects the agent's suggestion.)

		Agent: "Have you seen 'Alien'?"
		User: "No. Tell me about 'Princess diaries'"
		Outcome: OUTCOMETOPICCHANGE (The user shows no interest in the agent's suggestion and changes the topic.)

		Agent: "Have you seen 'Alien'?"
		User: "I told you I am not interested in sci-fi."
		Outcome: OUTCOMEIRRELEVANT (The agent made a wrong suggestion.)

		Guidelines for classification of the user sentiment:
		* SENTIMENTPOSITIVE: If the user expresses excitement, joy etc. Simply rejecting an agent's suggestion is not negative.
		* SENTIMENTNEGATIVE: If the user expresses frustration, irritation, anger etc. Simply rejecting an agent's suggestion is not negative.
		* SENTIMENTNEUTRAL: If the user expresses no specific emotion

		Remember:

		*   Do not make assumptions about the user's satisfaction or perception of helpfulness.
		*   Focus only on the objective flow of the conversation and how the user's response relates to the agent's previous turn.
		*   If the outcome is unclear based on the user's response, use OutcomeAmbiguous.

        {{ role "user" }}
		Here are the inputs:
		* history: (May be empty)
         {{#each history}}{{this.role}}: {{this.content}}{{~/each}}
		`

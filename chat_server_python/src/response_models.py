from langchain_core.pydantic_v1 import BaseModel, Field
from typing import List

from enum import Enum

class RESULT(Enum):
    UNDEFINED = 0
    SUCCESS = 1
    BAD_QUERY = 2
    UNSAFE = 3
    TOO_LONG = 4
    ERROR = 5

class SafetyError(Exception):
    def __init__(self, message, errors):            
        # Call the base class constructor with the parameters it needs
        super().__init__(message)
        self.errors = errors

class MoviePreference(BaseModel):
    genres: List[str] = Field(default_factory=list)
    directors: List[str] = Field(default_factory=list)
    actors: List[str] = Field(default_factory=list)
    other: List[str] = Field(default_factory=list, description="Other movie preferences (e.g., themes, plot elements, etc.)")

    def __init__(self, genres: List[str] = [], directors: List[str] = [], actors: List[str] = [], other: List[str] = []):
        super().__init__(genres=genres, directors=directors, actors=actors, other=other)

    def parse_dict(output):
        # Takes dict an returns an instance of MoviePreference 
        result = MoviePreference()
        try:
            if "genres" in output.keys():
                result.genres = output["genres"]
            if "directors" in output.keys():
                result.directors = output["directors"]
            if "actors" in output.keys():
                result.actors = output["actors"]
            if "other" in output.keys():
                result.other = output["other"]
        except Exception as e:
            print(e)
        return result

    def __eq__(self, other):
        """
        Compares two MoviePreference objects for equality, ignoring order of items in lists.

        Args:
            other: The other MoviePreference object to compare against.

        Returns:
            True if the two objects are equal, False otherwise.
        """
        if not isinstance(other, MoviePreference):
            return False
        return (
            set(self.genres) == set(other.genres) and
            set(self.directors) == set(other.directors) and
            set(self.actors) == set(other.actors) and
            set(self.other) == set(other.other)
        )

class UserPreferences(BaseModel):
    likes: MoviePreference
    dislikes: MoviePreference

    def __init__(self, likes: MoviePreference = MoviePreference(), dislikes: MoviePreference = MoviePreference()):
        super().__init__(likes=likes, dislikes=dislikes)
    
    def parse_dict(output):
        # Takes dict an returns an instance of MoviePreference 
        result = UserPreferences(
            likes=MoviePreference(),
            dislikes=MoviePreference()
        )
        try:
            likes = MoviePreference.parse_dict(output["likes"])
            result.likes = likes
            dislikes = MoviePreference.parse_dict(output["dislikes"])
            result.dislikes = dislikes
        except Exception as e:
            print(e)
        return result
    
    def dict(self):
        return {
            "likes": self.likes.dict(),
            "dislikes": self.dislikes.dict()
        }
    

    def __eq__(self, other):
        """
        Compares two UserPreferences objects for equality.

        Args:
            other: The other UserPreferences object to compare against.

        Returns:
            True if the two objects are equal, False otherwise.
        """
        if not isinstance(other, UserPreferences):
            return False
        return self.likes == other.likes and self.dislikes == other.dislikes
    
    
class MovieAgentResponse(BaseModel):
    answer: str = Field(description="The answer to the question in conversational language.")
    relevant_movies: List[str] = Field(default_factory=list, description="List of relevant movies.")
    wrong_query: bool = Field(description="True, if the user to ask the agent that violated it's mission.")

    def __init__(self, answer: str = "", relevant_movies: List[str] = []):
        super().__init__(answer=answer, relevant_movies=relevant_movies, wrong_query=False)

    def parse_dict(output):
        result = MovieAgentResponse()
        try:
            result.answer = output["answer"]
            result.relevant_movies = output["relevant_movies"]
            result.wrong_query = bool(output["wrong_query"])
        except Exception as e:
            result.answer=""
            print(e)
        return result
    
class AgentResponse():
    answer: str
    context: List[str]
    full_context: List[str]
    relevant_movies: List[str]
    error_message: str
    result: str
    preferences: UserPreferences


    def __init__(self, answer: str = "", context: List = [], full_context: List = [], relevant_movies: List = [], error_message: str = "", result:RESULT = RESULT.SUCCESS, preferences: UserPreferences = UserPreferences()):
        self.answer = answer
        self.context = context
        self.full_context = full_context,
        self.relevant_movies = relevant_movies
        self.error_message = error_message
        self.result = result.name
        self.preferences = preferences

       
class QueryRefinement(BaseModel):        
    final_search_query: str = Field(description="The final query that will be fed into the vector search engine that takes user preferences into account.")
    partial_search_query: List[str] = Field(default_factory=list, description="The query that is a result of summarizing the conversation that doesn't yet take user preferences into account.")
    explanation: str = Field(description="The explanation as to why you (the LLM) chose to construct the search query this way.")
    
    def __init__(self, final_search_query: str = "", partial_search_query: List[str] = [], explanation: str = ""):
        super().__init__(final_search_query=final_search_query, partial_search_query=partial_search_query, explanation=explanation)
    
    def parse_dict(output):
        result = QueryRefinement(
        )
        try:
            result.final_search_query = output["final_search_query"]
            result.partial_search_query = output["partial_search_query"]
            result.explanation = output["explanation"]
        except Exception as e:
            print(e)
        return result
from locust import HttpUser, task, between
import os
import random
import re
import requests

CHAT_SERVER = os.getenv("CHAT_SERVER")
MOCK_USER_SERVER = os.getenv("MOCK_USER_SERVER")


class ChatUser(HttpUser):
    wait_time = between(1, 2)
    host = CHAT_SERVER

    @task(3)
    def conversation(self):
        user_name = "chat_user"
        headers = self.do_login(user_name)
        self.client.delete(f"/history", headers=headers) 
        mock_user_message = get_random_user_starting_message()
        stop_conversation = False
        print("mock_user:", mock_user_message)

        while(not stop_conversation):
            try:
                agent_response = self.client.post(f"/chat", headers=headers, json={"content": mock_user_message})
                agent_message = sanitize_response(agent_response.json()['answer'])  
            except Exception as e:
                agent_response = "can you repeat that?"
            finally:
                print("agent:", agent_message)
            try:
                response_mood, response_type = get_random_response()
                mock_user_response = requests.post(f"{MOCK_USER_SERVER}/mockUserFlow", json={"data": {"expert_answer": agent_message, "response_mood": response_mood, "response_type": response_type}})
                mock_user_message = sanitize_response(mock_user_response.json()['result']['answer'])
                if response_type == 'END_CONVERSATION':
                    stop_conversation = True
            except Exception as e:
                mock_user_message = "can you repeat that?"
            finally:
                print("mock_user mood:", response_mood, "mock_user response type:", response_type)
                print("mock_user:", mock_user_message)
        self.client.delete(f"/history", headers=headers) 
        self.client.post(f"/logout", headers=headers) 

    def do_login(self, user_name):
        with requests.session() as session:
            session.post(f"{CHAT_SERVER}/login", headers={"user": user_name}) 
            headers={"user": user_name, "Cookie": f"movieguru=session_{user_name}"}
            return headers

    @task(1)
    def login(self):
        user_name = "login_user"
        headers = self.do_login(user_name) 
        self.client.delete(f"/history", headers=headers) 
        self.client.post(f"/logout", headers=headers) 


    @task(3)
    def preferences(self):
        user_name = "preferences_user"
        headers = self.do_login(user_name)
        self.client.post(f"/preferences", headers=headers, json=preferences_filled) 
        self.client.get(f"/preferences", headers=headers) 
        self.client.post(f"/preferences", headers=headers, json=preferences_updated) 
        self.client.get(f"/preferences", headers=headers) 
        self.client.post(f"/preferences", headers=headers, json=preferences_empty) 
        self.client.get(f"/preferences", headers=headers) 
        self.client.post(f"/logout", headers=headers) 
    
    @task(3)
    def startup(self):
        user_name = "startup_user"
        headers = self.do_login(user_name)
        self.client.get(f"/startup", headers=headers) 
        self.client.post(f"/logout", headers=headers) 


preferences_filled = {
    "likes":{
        "genres":["action", "comedy"],
        "actors":["Gene Hackman"],
        "directors": ["George Lucas"],
        "other":[]
    },
    "dislikes":{
        "genres":["romance"],
        "actors":[],
        "directors": [],
        "other":[]
    }
}

preferences_updated = {
    "likes":{
        "genres":["action", "comedy"],
        "actors":["Gene Hackman", "Tom Hanks"],
        "directors": ["George Lucas"],
        "other":[]
    },
    "dislikes":{
        "genres":["romance"],
        "actors":[],
        "directors": [],
        "other":[]
    }
}

preferences_empty = {
    "likes":{
        "genres":[],
        "actors":[],
        "directors": [],
        "other":[]
    },
    "dislikes":{
        "genres":[],
        "actors":[],
        "directors": [],
        "other":[]
    }
}

user_starting_messages = [
    "hello",
    "Can you recommend a good action movie?",
    "I'm looking for a movie with strong female characters. Any suggestions?",
    "What are some action movies that most people haven't heard of?",
    "Tell me about the best sci-fi movies of all time.",
    "What are some movies that make you think?",
    "I need a good laugh. Got any funny movie recommendations?",
    "What's a movie that everyone loves but you don't?",
    "What's your take on the latest movies?",
    "Can you suggest a movie that's similar to [Movie Title]?",
    "I'm in the mood for a classic horror film. Any ideas?"
]

def get_random_user_starting_message():
  random_index = random.randint(0, len(user_starting_messages) - 1)
  return user_starting_messages[random_index]


def sanitize_response(response):
  # Remove markdown
  response = re.sub(r"#+\s", "", response)        # Headings
  response = re.sub(r"\*\*([^*]+)\*\*", r"\1", response)  # Bold
  response = re.sub(r"\*([^*]+)\*", r"\1", response)    # Italics
  # Remove or replace special characters (customize as needed)
  response = re.sub(r"[`*;]", "", response)
  return response


RESPONSE_MOOD = [
    'POSITIVE',
    'NEGATIVE',
    'NEUTRAL',
    'RANDOM',
]

RESPONSE_TYPE = [
    'DIVE_DEEP',
    'CHANGE_TOPIC',
    'END_CONVERSATION',
    'CONTINUE',
    'RANDOM'
]

def get_random_response():
    mood = get_random_value_from_array(RESPONSE_MOOD)
    type = get_random_value_from_array(RESPONSE_TYPE)
    return mood, type

def get_random_value_from_array(array):
    random_index = random.randint(0, len(array) - 1) 
    return array[random_index]


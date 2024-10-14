from locust import HttpUser, task, between
import os
import random
import re

CHAT_SERVER = os.getenv("CHAT_SERVER")
MOCK_USER_SERVER = os.getenv("MOCK_USER_SERVER")


class MyUser(HttpUser):
    wait_time = between(1, 2)  # Wait between 1 and 5 seconds between tasks

    @task(3)
    def conversation(self):
        user_name = "chat_user"
        headers = self.login(user_name)
        self.client.delete(f"{CHAT_SERVER}/history", headers=headers) 
        mock_user_message = get_random_user_starting_message()
        stop_conversation = False
        print("mock_user:", mock_user_message)

        while(not stop_conversation):
            try:
                agent_response = self.client.post(f"{CHAT_SERVER}/chat", headers=headers, json={"content": mock_user_message})
                agent_message = sanitize_response(agent_response.json()['answer'])  
                print("agent:", agent_message)
            except Exception as e:
                agent_response = "can you repeat that?"
            try:
                response_mood, response_type = get_random_response()
                mock_user_response = self.client.post(f"{MOCK_USER_SERVER}/dummyUserFlow", json={"data": {"expert_answer": agent_message, "response_mood": response_mood, "response_type": response_type}})
                mock_user_message = sanitize_response(mock_user_response.json()['answer'])
                print("mock_user:", mock_user_message)
                if response_type == 'END_CONVERSATION':
                    stop_conversation = True
            except Exception as e:
                mock_user_message = "can you repeat that?"
        
        self.client.delete(f"{CHAT_SERVER}/history", headers=headers) 
        self.client.post(f"{CHAT_SERVER}/logout", headers=headers) 

    def login(self, user_name):
        response = self.client.post(f"{CHAT_SERVER}/login", headers={"user": user_name}) 
        cookies = response.cookies.get_dict() 
        session_cookie = cookies.get("session_id")
        headers={"user": user_name, "Cookie": f"session={session_cookie}"}
        return headers

    @task(1)
    def login(self):
        user_name = "login_user"
        headers = self.login(user_name) 
        self.client.delete(f"{CHAT_SERVER}/history", headers=headers) 
        self.client.post(f"{CHAT_SERVER}/logout", headers=headers) 




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
  """
  Sanitizes a response string by removing markdown and special characters.

  Args:
      response: The response string to sanitize.

  Returns:
      The sanitized response string.
  """
  # Remove markdown
  response = re.sub(r"#+\s", "", response)        # Headings
  response = re.sub(r"\*\*([^*]+)\*\*", r"\1", response)  # Bold
  response = re.sub(r"\*([^*]+)\*", r"\1", response)    # Italics
  response = re.sub(r"```[\s\S]*?```", "", response)    # Code blocks
  response = re.sub(r"`[^`]+`", r"\1", response)      # Inline code
  response = re.sub(r"\[([^\]]+)\]\([^)]+\)", r"\1", response)  # Links

  # Remove or replace special characters (customize as needed)
  response = re.sub(r"[^\w\s]", "", response)  # Remove all punctuation
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
    return { 'mood': mood, 'type': type }

def get_random_value_from_array(array):
    random_index = random.randint(0, len(array) - 1) 
    return array[random_index]



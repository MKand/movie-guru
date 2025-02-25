# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from locust import HttpUser, task, between
import requests
import os
import random
import string

class ChatUser(HttpUser):
    wait_time = between(1, 2)
    MOODS = ["POSITIVE", "NEUTRAL", "NEGATIVE"]
    RESPONSE_TYPE=["END_CONVERSATION", "CONTINUE", "CHANGE_TOPIC"]

    def on_stop(self):
        self.client.post("/logout")

    def on_start(self):
        # create random name
        name=''.join(random.choices(string.ascii_lowercase, k=8))

        headers = {
        "ApiKey": "ABC",
        "Content-Type": "application/json",
        "User": name
        }
        response = self.client.post(
            "/login", headers=headers, json={"inviteCode": ""})
        print(f"Login Headers {response.headers}")
        print(f"Login Response {response.content}")

        # Capture 'Set-Cookie' from the response headers
        set_cookie = response.headers.get('Set-Cookie').split(';', 1)[0]
        if set_cookie:
            print(f"Extracted cookie: {set_cookie}")
            # Stores it in the locust client.
            self.client.cookies.set("stored_cookie", set_cookie)
        else:
            print("No Set-Cookie header received.")
        
        self.helper_api_client = requests.Session()
        self.mock_url = os.getenv("MOCK_URL", "http://mockuser.mockuser.svc.cluster.local:80/mockUserFlow")
        print("using mock url", self.mock_url)

    @task(1)
    def healthcheck(self):
        response = self.client.get("/")

    @task(2)
    def chat_with_mock(self):
        endConv = False
        chat_answer = "Hi. How can I help you today?"
        while(endConv == False):
            response_type = random.choice(self.RESPONSE_TYPE)
            response_mood = random.choice(self.MOODS)
            # Post to mock user
            mock_response = self.helper_api_client.post(self.mock_url,
            json={
                "data":{
                    "expert_answer": chat_answer,
                    "response_mood": response_mood,
                    "response_type": response_type
                }
            })
            mock_response_json = mock_response.json()
            mock_response_answer = mock_response_json.get("result")["answer"]
            print(f"BOT: {chat_answer}\n")
            print(f"MOCK: {response_mood}: {response_type}: {mock_response_answer} \n")
            # Post to movie guru
            chat_response = self.client.post(
                        "/chat",
                        json={"content":mock_response_answer}
                    )
            chat_answer = chat_response.json()["answer"]
            
            if(response_type == "END_CONVERSATION"): 
                endConv = True
                print("--- END CONVERSATION ---")

    # @task(1)
    # def startup(self):
    #     self.client.get(
    #             "/startup",
    #         )
    
    # @task(2)
    # def preferences(self):
    #     self.client.post(f"/preferences", json={
    #             "Content": {
    #                 "likes": {"genres": ["action"]},
    #                 "dislikes": {}
    #             }
    #     })
    #     self.client.get(f"/preferences")

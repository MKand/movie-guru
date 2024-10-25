from locust import HttpUser, task, between

class ChatUser(HttpUser):
    wait_time = between(1, 2)
    @task(6)
    def conversation(self):
        self.client.post(f"/chat")

    @task(1)
    def login(self):
        self.client.post(f"/login")

    @task(1)
    def logout(self):
        self.client.post(f"/logout") 

    @task(1)
    def delete_history(self):
        self.client.post(f"/history") 
    
    @task(1)
    def get_history(self):
        self.client.get(f"/history") 


    @task(3)
    def preferences(self):
        self.client.post(f"/preferences") 
        self.client.get(f"/preferences") 
       
    @task(1)
    def startup(self):
        self.client.get(f"/startup") 

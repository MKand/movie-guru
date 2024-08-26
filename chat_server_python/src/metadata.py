from google.oauth2 import id_token
from google.auth.transport import requests
from google.auth import jwt as gjwt
from db import DatabaseConnection
import os


class AuthorizationError(Exception):
    def __init__(self, message):            
        # Call the base class constructor with the parameters it needs
        super().__init__(message)

class Metadata:
    def __init__(self,  metadata):
        self.token_audience = metadata['token_audience']
        self.history_length = metadata['history_length']
        self.max_user_message_len = metadata['max_user_message_len']
        self.cors_origin = metadata['cors_origin']
        self.retriever_length = metadata['retriever_length']
        self.google_chat_model_name = metadata['google_chat_model_name']
        self.google_embedding_model_name = metadata['google_embedding_model_name']
        self.front_end_domain = metadata['front_end_domain']


class MetadataHandler:
    def __init__(self, app_version):
        self.pool = DatabaseConnection.get_pool()
        self.app_version = app_version
    
    def get_metadata(self):
        query = """ SELECT * FROM app_metadata WHERE "app_version" = %s;"""
        db_response = DatabaseConnection.execute_query(self.pool, query, (self.app_version,), 'sel_single')
        if db_response:
            metadata = Metadata(db_response)
            return metadata
        else:
            raise Exception(f"Metadata not found for app_version {self.app_version}")



class UserLoginHandler:
    def __init__(self, token_audience):
        self.pool = DatabaseConnection.get_pool()
        self.token_audience = token_audience
    
    def handle_login(self,auth_header, invite_code):
        token = self._get_token(auth_header)
        user = self._verify_google_token(token)
        if self._check_user(user):
            return user
        else:
            if invite_code in self._get_invite_codes():
                self._create_user(user)
                return user
            else:
                raise AuthorizationError("Invalid invite code")
          
    # TODO: Add verification
    def _verify_google_token(self,token) -> str:
        claims = gjwt.decode(token, verify=False)
        if claims['aud'] == self.token_audience and claims['email_verified']==True:
            return claims['email']
        else:
            raise AuthorizationError("Invalid token")

    def _get_token(self, auth_header) -> str:
        token_parts = auth_header.split(" ")
        if len(token_parts) == 2 and token_parts[0].lower() == "bearer":
            access_token = token_parts[1]
            return access_token
        else:
            raise AuthorizationError("Invalid token")

    def _create_user(self, user: str) -> bool:
        query = """
            INSERT INTO user_logins (email) VALUES (%s)
            ON CONFLICT (email) DO UPDATE
            SET login_count = user_logins.login_count + 1;
        """
        db_response = DatabaseConnection.execute_query(self.pool, query, (user,), 'insert')
        return db_response


    def _check_user(self, user: str) -> bool:
        query = """ SELECT email FROM user_logins WHERE "email" = %s;"""
        db_response = DatabaseConnection.execute_query(self.pool, query, (user,), 'sel_single')
        return db_response and db_response['email'] == user
    
        

    def _get_invite_codes(self) -> bool:
        query = """SELECT code FROM invite_codes WHERE valid = True"""  # Fetch only valid codes
        db_response = DatabaseConnection.execute_query(self.pool, query, (), 'sel_multi')
        invite_codes = [row['code'] for row in db_response]  
        return invite_codes
            

        

if __name__ == "__main__":
    metadata_handler = MetadataHandler("v1")
    metadata = metadata_handler.get_metadata()
    print(metadata.__dict__)


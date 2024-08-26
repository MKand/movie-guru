from flask import Flask, request, session
from flask_cors import CORS
from langchain_core.messages import HumanMessage, AIMessage
from langchain_community.chat_message_histories import ChatMessageHistory
from response_models import RESULT, AgentResponse, UserPreferences
from agents import MovieAgent, QueryTransformAgent, UserPreferencesAgent
from metadata import UserLoginHandler, MetadataHandler, Metadata, AuthorizationError
from response_models import SafetyError
from flask_session import Session
from redis import StrictRedis as Redis
from functools import wraps


import os
import json
import logging

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)


APP_VERSION = os.getenv("APP_VERSION", "v1_localhost")
logger.info("APP_VERSION=",APP_VERSION)

app_metadata = MetadataHandler(APP_VERSION).get_metadata()

REDIS_HOST=os.getenv("REDIS_HOST")
logger.info("REDIS_HOST=",REDIS_HOST)

REDIS_PORT=int(os.getenv("REDIS_PORT","6378"))
logger.info("REDIS_PORT=",REDIS_PORT)

REDIS_PASSWORD=os.getenv("REDIS_PASSWORD")
FLASK_SECRET_KEY=os.getenv('FLASK_SECRET_KEY')
PROJECT_ID=os.getenv("PROJECT_ID")

SESSION_TYPE = 'redis'
SESSION_REDIS = Redis(host=REDIS_HOST, port=REDIS_PORT, password=REDIS_PASSWORD)

HISTORY_LEN = app_metadata.history_length
MAX_USER_MESSAGE_LEN = app_metadata.max_user_message_len

# CORS settings
CORS_ORIGINS_ENV = app_metadata.cors_origin
if isinstance(CORS_ORIGINS_ENV, str):
    CORS_ORIGINS = [origin.strip() for origin in CORS_ORIGINS_ENV.split(",")]

# Create the Flask application
app = Flask(__name__)
app.secret_key = FLASK_SECRET_KEY
app.config.from_object(__name__)
Session(app)

CORS(app, supports_credentials=True, origins=CORS_ORIGINS) 

app.config['CORS_HEADERS'] = 'Content-Type'
app.config['SESSION_COOKIE_SECURE'] = True
app.config['SESSION_COOKIE_HTTPONLY'] = True
app.config['SESSION_COOKIE_SAMESITE'] = 'None'
app.config['SESSION_COOKIE_DOMAIN'] = app_metadata.front_end_domain


movie_agent = MovieAgent(chat_model_name=app_metadata.google_chat_model_name,
                        embedding_model_name=app_metadata.google_embedding_model_name,
                        retiever_length=app_metadata.retriever_length)

query_transform_agent = QueryTransformAgent(chat_model_name=app_metadata.google_chat_model_name)

preferences_agent = UserPreferencesAgent(model_name=app_metadata.google_chat_model_name)

login_handler = UserLoginHandler(token_audience=app_metadata.token_audience)


CACHED_MOVIES = []


# After Request decorator to set SameSite attribute
@app.after_request
def after_request(response):
    response.headers.add('Set-Cookie', f'session={session.sid}; Secure; HttpOnly; SameSite=None; Path=/')
    return response


def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if "user" not in request.headers:
                    return {"server_response": "Unauthorized. User name required."}
        if session["user"] != request.headers["user"]:
                    return {"server_response": "Unauthorized. Incorrect user name."}
        if session is None:
                    return {"server_response": "Unauthorized. Not logged in."}
        return f(*args, **kwargs)
    return decorated_function

def __trim_history(user):
    try:
        history = _load_history(user)
        all_messages = history.messages
        if len(all_messages) > HISTORY_LEN:
            trimmed_messages = all_messages[-HISTORY_LEN:]
            history.messages = trimmed_messages        
            _save_history(user, history)
    except Exception as e:
        logger.error(e)
        raise e


def _load_history(user):
    history = ChatMessageHistory()
    history_json = SESSION_REDIS.get(user)
    if history_json is None:
        return history
    messages = [HumanMessage(msg['content']) if msg["type"] == "human" else AIMessage(msg['content']) for msg in json.loads(history_json)['messages']]
    history.add_messages(messages)
    return history

def _save_history(user, history):
    SESSION_REDIS.set(user, history.json())

def _delete_history(user):
    SESSION_REDIS.delete(user)
    return

def _update_movie_cache(movies):
    global CACHED_MOVIES
    movies_docs = movie_agent.retriever.invoke("good movies")
    movies = []
    for m in movies_docs:
        m_dict = json.loads(m.page_content)
        movies.append({"title": m_dict["title"], "poster": m_dict["poster"]})
    CACHED_MOVIES = movies

def _startup( movie_agent, preferences_agent):
    global CACHED_MOVIES
    try:
        preferences = preferences_agent.get_exisiting_preferences(session["user"])
        if len(CACHED_MOVIES) == 0:
           _update_movie_cache(movie_agent.retriever.invoke("good movies"))
        return AgentResponse(context=CACHED_MOVIES[:5], result=RESULT.SUCCESS, preferences=preferences.json())
    except AuthorizationError as ae:
        logger.error(ae)
    except Exception as e:
        logger.error(e)
        return

def __chat(history, user_input, query_transform_agent, movie_agent, preferences_agent):
        try:
            history.add_user_message(user_input)
            preferences = preferences_agent.process(user_input, session["user"])
            transformed_query = query_transform_agent.process(history.messages, preferences)
            movie_agent_response = movie_agent.process(history.messages, transformed_query, preferences)
            history.add_ai_message(movie_agent_response["answer"])
            result = RESULT.SUCCESS
            if movie_agent_response["wrong_query"] == True:
                result = RESULT.BAD_QUERY
            return AgentResponse(answer=movie_agent_response["answer"], context=movie_agent_response["filtered_context"], result=result, preferences=preferences.json())
        except AuthorizationError as ae:
            logger.error(ae)
        except SafetyError as se:
            logger.error(se)
            history.add_ai_message("There was a safety issue with the query. I cannot answer it.")
            return AgentResponse(answer=f"There is a safety issue with the query '{user_input}'. I cannot answer it.", result=RESULT.UNSAFE)
        except Exception as e:
            logger.error(e)
            history.add_ai_message("Oops something went wrong. Can you try again?")
            return AgentResponse(error_message="Oops something went wrong. Can you repeat your question?", result=RESULT.ERROR)

@app.route("/login", methods=['POST'])
def login():
    try:
        invite_code = request.json['inviteCode']
        auth_header = request.headers.get('Authorization')
        if auth_header:
            user = login_handler.handle_login(auth_header, invite_code)
            if user is not None:
                session["user"] = user
                return {"status": "success"}, 201
        else:
             return {"status": "Invalid invite code"}, 401
    except Exception as e:
         return {"server_response": "Internal server error"}, 500
    

@app.route("/logout", methods=['GET'])
@login_required
def logout():
    session.pop("user", default=None)
    session.pop("history", default=None)
    return {"status": "success"}, 201
   
@app.route("/user", methods=['GET'])
@login_required
def user():
    if session.get("user") is not None:
       return session["user"],200
    else:
        return {"answer": "not found"}, 404


@app.route("/history", methods=['GET'])
@login_required
def history():
    try:
        history = _load_history(user=session["user"])
        history_list = []

        messages = history.messages
        
        for message in messages:
            if isinstance(message, HumanMessage):
                history_list.append({"sender": "user", "message": message.content})
            elif isinstance(message, AIMessage):
                history_list.append({"sender": "agent", "message": message.content})
        return history_list,200
    except Exception as e:
        logger.error(e)
        return {"server_response", "Internal server error"}, 500

@app.route("/history", methods=['DELETE'])
@login_required
def delete_history():
    try:
        _delete_history(user=session["user"])
        return {},204
    except Exception as e:
        logger.error(e)
        return {"server_response", "Internal server error"}, 500

@app.route("/preferences", methods=['GET'])
@login_required
def preferences():
    try:
        preferences =  preferences_agent.get_exisiting_preferences(user=session["user"])
        return preferences.json(), 200
    except Exception as e:
        logger.error(e)
        return {"server_response", "Internal server error"}, 500
    
@app.route("/preferences", methods=['POST'])
@login_required
def preferences_update():
    try:
        preferences_raw = request.json["content"]
        preferences = UserPreferences.parse_dict(preferences_raw)
        preferences_agent.update(preferences, user=session["user"])
        return {}, 201
    except Exception as e:
        logger.error(e)
        return {"server_response", "Internal server error"}, 500
    
@app.route('/chat', methods=['POST'])
@login_required
def chat():
    try:
        user_message = request.json['content']
        if len(user_message) > MAX_USER_MESSAGE_LEN:
            ai_response = AgentResponse(answer="Your query is too long. Please shorten it and try again.", result=RESULT.TOO_LONG)
            return ai_response.__dict__, 200

        history = _load_history(user=session["user"])
        ai_response = __chat(history=history, 
                            user_input=user_message, 
                            query_transform_agent = query_transform_agent, 
                            movie_agent=movie_agent, 
                            preferences_agent = preferences_agent)
            
        _save_history(session["user"], history)
        __trim_history(session["user"])
            
        if ai_response.result is not RESULT.ERROR or RESULT.UNKNOWN:
            return ai_response.__dict__, 200
        elif ai_response.result is not RESULT.ERROR:
            return ai_response.__dict__, 500
    except Exception as e:
        logger.error(e)
        return {"server_response", "Internal server error"}, 500

@app.route('/startup', methods=['GET'])
@login_required
def startup():
    try:
        ai_response = _startup(
                            movie_agent=movie_agent, 
                            preferences_agent = preferences_agent)
            
        if ai_response.result is not RESULT.ERROR or RESULT.UNKNOWN:
            return ai_response.__dict__, 200
        elif ai_response.result is not RESULT.ERROR:
            return ai_response.__dict__, 500
    except Exception as e:
        logger.error(e)
        return {"server_response", "Internal server error"}, 500

@app.route('/metadata', methods=['GET'])
def metadata():
    try:
        output =  {"history_length": app_metadata.history_length,
                    "max_user_message_len": app_metadata.max_user_message_len,
                    "retriever_length": app_metadata.retriever_length,
                    "backend_version": APP_VERSION}
        return output, 200
       
    except Exception as e:
        logger.error(e)
        return {"server_response", "Internal server error"}, 500

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5001)
import psycopg
import psycopg_pool
from langchain_google_cloud_sql_pg import PostgresEngine, PostgresVectorStore
import os

POSTGRES_DB_USER_PASSWORD=os.getenv("POSTGRES_DB_USER_PASSWORD")
POSTGRES_DB_INSTANCE = os.getenv("POSTGRES_DB_INSTANCE")
POSTGRES_DB_USER = os.getenv("POSTGRES_DB_USER")
POSTGRES_HOST= os.getenv("POSTGRES_HOST")
POSTGRES_PORT= int(os.getenv("POSTGRES_PORT", "5432"))
POSTGRES_DB_NAME=os.getenv("POSTGRES_DB_NAME", "fake-movies-db")
POSTGRES_MAX_CONNECTIONS=int(os.getenv("POSTGRES_MAX_CONNECTIONS", "5"))
POSTGRES_MAX_IDLE=int(os.getenv("POSTGRES_MAX_IDLE", "20"))

POSTGRES_IP_TYPE= os.getenv("POSTGRES_IP_TYPE", "PRIVATE")
PROJECT_ID = os.getenv("PROJECT_ID")
REGION = "europe-west4"

class DatabaseConnection:
    __instance = None
    pool = None

    def __new__(cls):
        if cls.__instance is None:
            cls.__instance = super().__new__(cls)
            try:
                connection_string = f"postgresql://{POSTGRES_DB_USER}:{POSTGRES_DB_USER_PASSWORD}@{POSTGRES_HOST}:{POSTGRES_PORT}/{POSTGRES_DB_NAME}"
                cls.pool = psycopg_pool.ConnectionPool(connection_string,
                                            min_size=1,
                                            max_size=POSTGRES_MAX_CONNECTIONS,
                                            max_idle=POSTGRES_MAX_IDLE)
                
            except psycopg.OperationalError as e:
                print(f"Error connecting to the database: {e}")
                # Optionally, you might want to raise the exception or handle it in a more specific way.
        return cls.__instance

    @staticmethod
    def get_pool():
        if not DatabaseConnection.__instance:
            DatabaseConnection()
        return DatabaseConnection.__instance.pool
    
    @staticmethod
    def get_pg_engine():
        pg_engine = PostgresEngine.from_instance(
            project_id=PROJECT_ID,
            database=POSTGRES_DB_NAME,
            region=REGION,
            instance=POSTGRES_DB_INSTANCE,
            user=POSTGRES_DB_USER,
            password=POSTGRES_DB_USER_PASSWORD,
            ip_type=POSTGRES_IP_TYPE
        )
        return pg_engine
    
    @staticmethod
    def execute_query(pool, sql_raw, params, qry_type):
        with pool.connection() as conn:
            cur = conn.cursor(row_factory=psycopg.rows.dict_row)
            try:
                if qry_type == 'sel_multi':
                    results = cur.execute(sql_raw, params).fetchall()
                elif qry_type == 'sel_single':
                    results = cur.execute(sql_raw, params).fetchone()
                elif qry_type == 'insert':
                    cur.execute(sql_raw, params)
                    conn.commit()
                    results = True
                elif qry_type == 'update':
                    cur.execute(sql_raw, params)
                    conn.commit()
                    results = True
                else:
                    raise Exception('Invalid query type defined.')
            except psycopg.OperationalError as err:
                app.logger.error(f'Error querying: {err}')
            except psycopg.ProgrammingError as err:
                app.logger.error('Database error via psycopg.  %s', err)
                results = False
            except psycopg.IntegrityError as err:
                app.logger.error('PostgreSQL integrity error via psycopg.  %s', err)
                results = False
            return results
import postgres from 'postgres';

let sql: postgres.Sql<{}> | null = null; // Initialize sql as null

export async function openDB(): Promise<postgres.Sql<{}> |null> {
  if (sql) {
    return sql; // Return existing connection if already opened
  }

  // const POSTGRES_DB_USER_PASSWORD = process.env.POSTGRES_DB_USER_PASSWORD;
  // const POSTGRES_HOST = process.env.POSTGRES_HOST;
  // const POSTGRES_DB_NAME = process.env.POSTGRES_DB_NAME;
  // const POSTGRES_DB_USER = process.env.POSTGRES_DB_USER;

  const POSTGRES_DB_USER_PASSWORD = "minimalpassword";
  const POSTGRES_HOST = "localhost";
  const POSTGRES_DB_NAME = "fake-movies-db";
  const POSTGRES_DB_USER = "minimal-user";

  if (!POSTGRES_DB_USER_PASSWORD || !POSTGRES_HOST || !POSTGRES_DB_NAME ||!POSTGRES_DB_USER) {
    console.error('Missing environment variables for database connection');
    return null;
  }

  try {
    sql = postgres({
      host: POSTGRES_HOST,
      user: POSTGRES_DB_USER,
      password: POSTGRES_DB_USER_PASSWORD,
      port: 5432,
      database: POSTGRES_DB_NAME,
      max: 5,
      idle_timeout: 30000,
    });

    await sql`SELECT NOW()`;
    console.log('DB opened successfully');
    return sql;
  } catch (err) {
    console.error(err);
    throw err;
  }
}

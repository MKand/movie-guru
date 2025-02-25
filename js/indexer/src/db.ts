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

import postgres from 'postgres';

let sql: postgres.Sql<{}> | null = null; // Initialize sql as null

export async function openDB(): Promise<postgres.Sql<{}> | null> {
  if (sql) {
    return sql; // Return existing connection if already opened
  }

  const POSTGRES_DB_USER_PASSWORD = "mainpassword"; //process.env.POSTGRES_DB_USER_PASSWORD;
  const POSTGRES_HOST = process.env.POSTGRES_HOST;
  const POSTGRES_DB_NAME = process.env.POSTGRES_DB_NAME;
  const POSTGRES_DB_USER = "main";//process.env.POSTGRES_DB_USER;

  if (!POSTGRES_DB_USER_PASSWORD || !POSTGRES_HOST || !POSTGRES_DB_NAME || !POSTGRES_DB_USER) {
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

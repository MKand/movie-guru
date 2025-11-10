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

import Database from 'better-sqlite3';
import * as sqliteVec from 'sqlite-vec';

let db: Database.Database | null = null;

export async function openDB(): Promise<Database.Database> {
  if (db) {
    return db;
  }

  try {
    // Open SQLite database (use a file path or ':memory:' for in-memory)
    const dbPath = process.env.SQLITE_DB_PATH || './movies.db';
    db = new Database(dbPath);
    
    // Load sqlite-vec extension
    sqliteVec.load(db);
    
    console.log('SQLite database opened successfully at', dbPath);
    
    // Create movies table if it doesn't exist
    db.exec(`
      CREATE TABLE IF NOT EXISTS movies (
        tconst TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        runtime_mins INTEGER,
        genres TEXT,
        rating REAL,
        released INTEGER,
        actors TEXT,
        director TEXT,
        plot TEXT,
        poster TEXT,
        content TEXT
      )
    `);
    
    // Create vec0 virtual table for vector search (768 dimensions for text-embedding-005)
    db.exec(`
      CREATE VIRTUAL TABLE IF NOT EXISTS vec_movies USING vec0(
        tconst TEXT PRIMARY KEY,
        embedding FLOAT[768]
      )
    `);
    
    console.log('Tables created/verified successfully');
    
    return db;
  } catch (err) {
    console.error('Error opening SQLite database:', err);
    throw err;
  }
}

export function closeDB(): void {
  if (db) {
    db.close();
    db = null;
    console.log('SQLite database closed');
  }
}
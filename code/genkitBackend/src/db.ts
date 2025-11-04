
import * as sqlite3 from 'sqlite3';
import { PreferenceItem } from './userPreferencesTypes';

const DB_FILE = './user_preferences.db';

export class UserPreferencesDB {
  private db: sqlite3.Database;

  constructor() {
    this.db = new sqlite3.Database(DB_FILE);
  }

  async init(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.db.run(
        `CREATE TABLE IF NOT EXISTS user_profiles (
          userId TEXT PRIMARY KEY,
          profile TEXT NOT NULL
        )`,
        (err) => {
          if (err) {
            return reject(err);
          }
          resolve();
        }
      );
    });
  }

  async get(userId: string): Promise<PreferenceItem[]> {
    return new Promise((resolve, reject) => {
      this.db.get(
        'SELECT profile FROM user_profiles WHERE userId = ?',
        [userId],
        (err, row: { profile: string }) => {
          if (err) {
            return reject(err);
          }
          if (row) {
            resolve(JSON.parse(row.profile));
          } else {
            resolve([]);
          }
        }
      );
    });
  }

  async update(userId: string, profile: PreferenceItem[]): Promise<void> {
    return new Promise((resolve, reject) => {
      const profileJson = JSON.stringify(profile);
      this.db.run(
        'INSERT OR REPLACE INTO user_profiles (userId, profile) VALUES (?, ?)',
        [userId, profileJson],
        (err) => {
          if (err) {
            return reject(err);
          }
          resolve();
        }
      );
    });
  }
}

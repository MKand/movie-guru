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

import fs from 'fs/promises';
import { parse } from 'csv-parse';
import { IndexerFlow } from './indexerFlow';
import { MovieContext } from './types';
import { openDB } from './db';

const FILENAME = process.env.DATASETFILE || "/dataset/movies_with_posters.csv";

const sleep = (ms: number): Promise<void> =>
  new Promise((resolve) => setTimeout(resolve, ms));


export async function processMovies() {
  try {
    const fileContent = await fs.readFile(FILENAME, 'utf8');

    await sleep(100);

    const parser = parse({
      delimiter: '\t',
      from_line: 2,
    });

    const records: any[] = [];
    parser.on('readable', () => {
      let record;
      while ((record = parser.read()) !== null) {
        records.push(record);
      }
    });

    await new Promise((resolve, reject) => {
      parser.on('end', resolve);
      parser.on('error', reject);
      parser.write(fileContent);
      parser.end();
    });

    let index = 0;
    const db = await openDB();
    console.log("starting indexing")

    for (const record of records) {
      const year = parseFloat(record[1]);
      const rating = parseFloat(record[5]);
      const runtime = parseFloat(record[6]);

      const movieContext: MovieContext = {
        title: record[0],
        runtimeMinutes: Math.trunc(runtime), // Use Math.trunc() to get the integer part
        genres: record[7].split(', '),
        rating: rating,
        plot: record[4],
        released: Math.trunc(year), // Use Math.trunc() to get the integer part
        director: record[3],
        actors: record[2].split(', '),
        poster: record[9],
        tconst: index.toString(),
      };
      try {
        await IndexerFlow(movieContext)
        // Slowing down to avoid pesky rate limits
        await sleep(1500)
        console.log("processed ", record[0])
      } catch (err) {
        console.error('Error loading movie: ', record[0], err);
      }
      index++;

    }
    console.log(`finished indexing ${index} documents.`)

  } catch (err) {
    console.error('Error opening or processing file:', err);
  }
}
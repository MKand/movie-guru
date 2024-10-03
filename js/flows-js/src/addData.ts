import fs from 'fs/promises'; // Use fs/promises for async file operations
import { parse } from 'csv-parse'; // Use a CSV parsing library
import { runFlow } from '@genkit-ai/flow';
import { IndexerFlow } from './indexerFlow'; // Assuming IndexerFlow is defined in indexerFlow.ts
import { MovieContext } from './movieFlowTypes'; // Assuming MovieContext is defined in movieFlowTypes.ts

async function processMovies( ctx: any) { // Replace 'any' with the correct context type
  try {
    const fileContent = await fs.readFile('/dataset/movies_with_posters.csv', 'utf8');

    const parser = parse({
      delimiter: '\t', // Set the delimiter to tab
      from_line: 2, // Start reading from the second line (skip header)
    });

    const records: any[] = []; // Change 'any' to a more specific type if possible
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
        const response = await runFlow(IndexerFlow, movieContext);
      } catch (err) {
        console.error('Error loading movie: ', record[0], err);
      }
      index++;
    }
  } catch (err) {
    console.error('Error opening or processing file:', err);
  }
}
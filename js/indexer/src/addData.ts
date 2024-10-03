import fs from 'fs/promises'; 
import { parse } from 'csv-parse'; 
import { runFlow } from '@genkit-ai/flow';
import { IndexerFlow } from './indexerFlow'; 
import { MovieContext } from './types'; 

export async function ProcessMovies() { 
  try {
    const fileContent = await fs.readFile('/dataset/movies_with_posters.csv', 'utf8');

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

      console.log(movieContext)
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
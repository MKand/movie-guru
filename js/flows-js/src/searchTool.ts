import { defineTool } from "@genkit-ai/ai";
import { z } from 'zod';
import { openDB } from './db';
import { MovieContextSchema, MovieContext } from "./movieFlowTypes";

var squel = require("squel");


export const getMovies = async (query: string ): Promise<MovieContext[]> => {
  const db = await openDB();
  if (!db) {
    throw new Error('Database connection failed');
  }

  const t = "rating > 2.5"
  try {

    const sqlQuery = squel.select()
    .field("title")
    .field("poster")
    .field("released")
    .field("runtime_mins")
    .field("rating")
    .field("genres")
    .field("director")
    .field("actors")
    .field("plot")
    .field("tconst")
    .from("movies")
    .where(query)
    .toString()

      const results = await db.unsafe(sqlQuery)

      console.log("Results:", results);
      return results.map((row: any) => ({ 
        title: row.title,
        runtime_minutes: Number(row.runtime_mins),
        genres: row.genres.split(","),
        rating: Number(row.rating),
        plot: row.plot,
        released: Number(row.released),
        director: row.director,
        actors: row.actors.split(","),
        poster: row.poster,
        tconst: row.tconst,
      })); 
    
  } catch (error) {
    console.error("Error executing query:", error);
    return []
  }  
 
};

import { processMovies } from './addData';
import { ai } from './genkitConfig'
import { IndexerFlow } from './indexerFlow';

ai.startFlowServer({
  flows: [IndexerFlow],
});

processMovies();


import { ai } from './genkitConfig'

import { UserProfileFlow } from './userProfileFlow'
export { UserProfileFlowPrompt } from './userProfileFlow'

import { QueryTransformFlow } from './queryTransformFlow'
export { QueryTransformPrompt } from './queryTransformFlow'

import { MovieDocFlow } from './movieDocFlow'

import { MovieRAGFlow } from './movieRAGFlow'

import { MovieFlow } from './movieFlow'
export { MovieFlowPrompt } from './movieFlow'


ai.startFlowServer({
    flows: [UserProfileFlow, QueryTransformFlow, MovieFlow, MovieDocFlow, MovieRAGFlow],
  });
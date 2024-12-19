import { ai } from './genkitConfig'

import { UserProfileFlowPrompt, UserProfileFlow} from './userProfileFlow'
import { QueryTransformPrompt, QueryTransformFlow} from './queryTransformFlow'

ai.startFlowServer({
    flows: [UserProfileFlow, QueryTransformFlow],
  });
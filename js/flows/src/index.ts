import { ai } from './genkitConfig'

import { UserProfileFlowPrompt, UserProfileFlow} from './userProfileFlow'

ai.startFlowServer({
    flows: [UserProfileFlow],
  });
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

import { startFlowServer } from '@genkit-ai/express';

import { UserPreferenceFlow } from './userPreferencesFlow'

import { QueryTransformFlow } from './queryTransformFlow'

import { DocSearchFlow } from './docRetriever'

import { MovieFlow } from './movieFlow'

import { JudgeFlow } from './judgeFlow'

import { ChatFlow } from './chatFlow';

import { SafetyIssueFlow } from './safetyFlow';

export function parseBooleanfromField(field: string|null|undefined): boolean{
  return field != null ? JSON.parse(field): false;
}


startFlowServer({
    flows: [UserPreferenceFlow, QueryTransformFlow, MovieFlow, DocSearchFlow, JudgeFlow, ChatFlow, SafetyIssueFlow],
  });

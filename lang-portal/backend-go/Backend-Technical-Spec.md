# Backend Technical Specification

## Business Requirements

- Provide a learning portal for language school:
- Inventory of vocabulary
- Include a Learning Record Store (LRS) with history and scores for each student
- Launchpad for multiple learing apps

## Technical Requirements

- Backend will be written in Go
- Database will be SQLite
- API will return JSON
- No authentication or authorization
- Only single user

## Database Schema

Tables:
- words - vocabulary words
  - id integer
  - english string
  - spanish string
  - level enum
- word_groups - many-to-many join of words and groups
  - id integer
  - word_id integer
  - group_id integer
- groups - groups of words
  - id integer
  - name string
- study_sessions - records of study sesions grouping word_review_items
  - id integer
  - group_id integer
  - study_session_id integer
  - created_at datetime
- study_activities - a specific study activity
  - id integer
  - study_session_id integer
  - group_id integer
  - created_at datetime
- word_review_items - a record of word practice, correct or incorrect
  - id integer
  - word_id integer
  - study_activity_id integer
  - correct boolean
  - created_at datetime

## API Specification

- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats
- GET /api/study_activities/:id
- GET /api/study_activities/:id/study_sessions
- GET /api/words
  - pagination with 100 items per page
- GET /api/words/:id
- GET /api/groups
  - pagination with 100 items per page
- GET /api/groups/:id
- GET /api/groups/:id/words
- GET /api/groups/:id/study_sessions
- POST /api/study_activities
  - required params: group_id, study_activity_id
- GET /api/study_sessions
  - pagination with 100 items per page
- GET /api/study_sessions/:id
- GET /api/study_sessions/:id/words
- POST /api/reset_history
- POST /api/full_reset
- PORT /api/study_sessions/:id/words/:word_id/review
# Frontend Technical Specification

## Pages

### Dashboard `/dashboard`

The default dashboard page which displays and provides quick stats including:
- Last Study Session
  - Last study session details
- Study Progress
  - Progress details
  - Progress pie chart
- Quick Stats
  - Success rate
  - Total study sessions
  - Total active groups
  - study streak
- Start Studying Button
  - Link to first group

#### Required API Endpoints
- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats

### Study Activities `/study_activities`

Study activities page which displays a card for each study activity:
- Name of study activity
- Button to start study activity
- View more details

#### Required API Endpoints
- GET /api/study_activities

### Study Activity Detail `/study_activities/:id`

This pages show sthe details of a specific study activity and past activities:
- Name of study activity
- Description of study activity
- Button to start study activity
- Paginated list of past study activities:
  - id
  - name
  - group name
  - created at (used as start time timestamp for pagination)
  - updated at (used as end time timestamp for pagination)

#### Required API Endpoints
- GET /api/study_activities/:id
- GET /api/study_activities/:id/study_sessions

### Study Activity Launch `/study_activities/:id/launch`

This page allows the user to launch a study activity:
- Name of study activity
- Description of study activity
- Form to launch study activity
  - select field for group
  - launch button

#### Behavior
When the form is submmmitted a new tab opens wth the study activitybased on it's URL in the database.

#### Required API Endpoints
- POST /api/study_activities

### Words Index `/words`
This page will show a paginated word list:
- Columns
    - English
    - Spanish
    - Correct Count
    - Incorrect Count
    - Part of Speech
  - Pagination with 100 items per page
  - Clinking the Spanish word will take you to the word show page

#### Required API Endpoints
- GET /api/words

### Word Show `/words/:id`
This page will show the details of a specific word:
- English
- Spanish
- Study Statistics
  - Correct Count
  - Incorrect Count
- Word Groups
  - Show as a series of tags
  - When group name is clicked, take to the group show page

#### Required API Endpoints
- GET /api/words/:id

### Groups Index `/groups/`
This page will show a paginated group list:
- Columns
  - Group Name
  - Word Count
- Clinking the group name will take you to the group show page

#### Required API Endpoints
- GET /api/groups

### Group Show `/groups/:id`
This page will show the details of a specific group:
- Columns   
  - Group Name
  - Group Statistics
    - Total word Count
  - Words in Group
    - paginated list of words
    - Use same component as the words index page
  - Study Sessions
    - paginated list of study sessions 
    - Use same component as the study activities index page

#### Required API Endpoints
- GET /api/groups/:id (the name and groups stats)
- GET /api/groups/:id/words
- GET /api/groups/:id/study_sessions

## Study Sessions Index `/study_sessions`
This page will show a paginated study session list:
- Columns
  - Id
  - Activty Name
  - Group Name
  - Start Time
  - End Time
  - Number of Review Items
- Clinking the word studied will take you to the word show page

#### Required API Endpoints
- GET /api/study_sessions

### Study Session Show `/study_sessions/:id`
This page will show the details of a specific study session:
- Study Session Details
  - Activty Name
  - Group Name
  - Start Time
  - End Time
  - Number of Review Items
- Words Review Items (Paginated list of words)
  - Use same component as the words index page

#### Required API Endpoints
- GET /api/study_sessions/:id
- GET /api/study_sessions/:id/words

### Settings Page `/settings`
This page will make user configuration changes to the study portal:
- Theme Selection - Light, Dark, 1980s!
- Reset History Button - Clears all study sessions
- Full Reset Button - Drop all tables and re-create with seed data

#### Required API Endpoints
- POST /api/reset_history
- POST /api/full_reset

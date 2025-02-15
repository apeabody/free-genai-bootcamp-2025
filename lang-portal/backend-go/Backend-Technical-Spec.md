# Backend Technical Specification

## Business Requirements

- Provide a learning portal for language school:
- Inventory of vocabulary
- Include a Learning Record Store (LRS) with history and scores for each student
- Launchpad for multiple learing apps

## Technical Requirements

- Backend will be written in Go
- Mage is a task runner for Go
- The API will be built using Gin
- Database will be SQLite
- API will return JSON
- No authentication or authorization
- Only single user

## Database Schema

The database will be a single sqlite3 file called `words.db` in the root application directory.

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

### Dashboard Endpoints

#### GET /api/dashboard/last_study_session
Returns details of the user's most recent study session.

**Response:**
```json
{
  "id": 1,
  "activity_name": "Vocabulary Practice",
  "group_name": "Basic Greetings",
  "start_time": "2025-02-14T19:30:00Z",
  "end_time": "2025-02-14T20:00:00Z",
  "review_items_count": 20,
  "correct_count": 15,
  "incorrect_count": 5,
  "score": 75
}
```

#### GET /api/dashboard/study_progress
Returns the user's study progress statistics.

**Response:**
```json
{
  "total_sessions": 10,
  "completed_sessions": 8,
  "success_rate": 75,
  "total_words": 100,
  "mastered_words": 60,
  "in_progress_words": 30,
  "not_started_words": 10,
  "study_streak": 5,
  "best_streak": 7
}
```

#### GET /api/dashboard/quick_stats
Returns quick statistics for the dashboard.

**Response:**
```json
{
  "success_rate": 75,
  "total_study_sessions": 10,
  "total_active_groups": 5,
  "study_streak": 5,
  "total_study_time": 300,
  "words_mastered": 60
}
```

### Study Activities

#### GET /api/study_activities/:id
Returns details of a specific study activity.

**Response:**
```json
{
  "id": 1,
  "name": "Vocabulary Practice",
  "description": "Practice vocabulary words with flashcards",
  "launch_url": "https://example.com/vocab-practice",
  "created_at": "2025-02-14T19:00:00Z",
  "updated_at": "2025-02-14T20:00:00Z",
  "group_name": "Basic Greetings",
  "total_sessions": 5
}
```

#### GET /api/study_activities/:id/study_sessions
Returns a list of study sessions for a specific activity.

**Response:**
```json
{
  "study_sessions": [
    {
      "id": 1,
      "name": "Vocabulary Practice",
      "group_name": "Basic Greetings",
      "created_at": "2025-02-14T19:00:00Z",
      "updated_at": "2025-02-14T20:00:00Z",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 50,
    "per_page": 10
  }
}
```

### Words

#### GET /api/words
Returns a paginated list of vocabulary words.

**Response:**
```json
{
  "words": [
    {
      "id": 1,
      "english": "hello",
      "spanish": "hola",
      "correct_count": 10,
      "incorrect_count": 2,
      "part_of_speech": "interjection"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 10,
    "total_items": 1000,
    "per_page": 100
  }
}
```

#### GET /api/words/:id
Returns details of a specific word.

**Response:**
```json
{
  "id": 1,
  "english": "hello",
  "spanish": "hola",
  "study_statistics": {
    "correct_count": 10,
    "incorrect_count": 2,
    "mastery_level": 0.83
  },
  "groups": [
    {
      "id": 1,
      "name": "Basic Greetings"
    },
    {
      "id": 2,
      "name": "Common Phrases"
    }
  ]
}
```

### Groups

#### GET /api/groups
Returns a paginated list of word groups.

**Response:**
```json
{
  "groups": [
    {
      "id": 1,
      "name": "Basic Greetings",
      "word_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 500,
    "per_page": 100
  }
}
```

#### GET /api/groups/:id
Returns details of a specific group.

**Response:**
```json
{
  "id": 1,
  "name": "Basic Greetings",
  "statistics": {
    "total_word_count": 20,
    "mastered_words": 15,
    "in_progress_words": 5,
    "average_mastery": 0.75
  }
}
```

#### GET /api/groups/:id/words
Returns a paginated list of words in a specific group.

**Response:**
```json
{
  "words": [
    {
      "id": 1,
      "english": "hello",
      "spanish": "hola",
      "correct_count": 10,
      "incorrect_count": 2,
      "part_of_speech": "interjection"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 20,
    "per_page": 100
  }
}
```

#### GET /api/groups/:id/study_sessions
Returns a paginated list of study sessions for a specific group.

**Response:**
```json
{
  "study_sessions": [
    {
      "id": 1,
      "activity_name": "Vocabulary Practice",
      "start_time": "2025-02-14T19:00:00Z",
      "end_time": "2025-02-14T20:00:00Z",
      "review_items_count": 20,
      "score": 75
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 50,
    "per_page": 10
  }
}
```

### Study Activities Management

#### POST /api/study_activities
Creates a new study activity.

**Request Body:**
```json
{
  "group_id": 1,
  "study_activity_id": "vocabulary_practice"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "Vocabulary Practice",
  "group_id": 1,
  "group_name": "Basic Greetings",
  "launch_url": "https://example.com/vocab-practice?session=1",
  "created_at": "2025-02-14T20:00:00Z"
}
```

### Study Sessions

#### GET /api/study_sessions
Returns a paginated list of study sessions.

**Response:**
```json
{
  "study_sessions": [
    {
      "id": 1,
      "activity_name": "Vocabulary Practice",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-14T19:00:00Z",
      "end_time": "2025-02-14T20:00:00Z",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 10,
    "total_items": 1000,
    "per_page": 100
  }
}
```

#### GET /api/study_sessions/:id
Returns details of a specific study session.

**Response:**
```json
{
  "id": 1,
  "activity_name": "Vocabulary Practice",
  "group_name": "Basic Greetings",
  "start_time": "2025-02-14T19:00:00Z",
  "end_time": "2025-02-14T20:00:00Z",
  "review_items_count": 20,
  "score": 75,
  "correct_count": 15,
  "incorrect_count": 5
}
```

#### GET /api/study_sessions/:id/words
Returns a paginated list of words reviewed in a specific study session.

**Response:**
```json
{
  "words": [
    {
      "id": 1,
      "english": "hello",
      "spanish": "hola",
      "correct_count": 10,
      "incorrect_count": 2,
      "part_of_speech": "interjection",
      "review_result": {
        "correct": true,
        "response_time": 1.5
      }
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 20,
    "per_page": 100
  }
}
```

### System Management

#### POST /api/reset_history
Resets all study history while preserving words and groups.

**Response:**
```json
{
  "success": true,
  "message": "Study history has been reset",
  "items_cleared": {
    "study_sessions": 100,
    "word_reviews": 2000
  }
}
```

#### POST /api/full_reset
Resets the entire database to initial state.

**Response:**
```json
{
  "success": true,
  "message": "Database has been reset to initial state",
  "items_cleared": {
    "words": 1000,
    "groups": 50,
    "study_sessions": 100,
    "word_reviews": 2000
  }
}
```

#### POST /api/study_sessions/:id/words/:word_id/review
Records a word review in a study session.

**Request Body:**
```json
{
  "correct": true,
  "response_time": 1.5
}
```

**Response:**
```json
{
  "success": true,
  "review": {
    "word_id": 1,
    "session_id": 1,
    "correct": true,
    "response_time": 1.5,
    "new_mastery_level": 0.85
  }
}
```

## Scripts (tasks)

### Initilze Database
This task will initialize the sqlite3 database by creating the necessary tables.

### Migrate Databaes
This task will apply any pending migrations to the database.

Migrations live in a `migrations` directory and run in the order of their file name.

### Seed Database
This task will add some sample data to the database.

All seed files live in a `seeds` directory.

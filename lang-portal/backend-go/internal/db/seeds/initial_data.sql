-- Insert sample words
INSERT INTO words (english, spanish, level) VALUES
('hello', 'hola', 'beginner'),
('goodbye', 'adiós', 'beginner'),
('please', 'por favor', 'beginner'),
('thank you', 'gracias', 'beginner'),
('good morning', 'buenos días', 'beginner'),
('good night', 'buenas noches', 'beginner'),
('water', 'agua', 'beginner'),
('food', 'comida', 'beginner'),
('house', 'casa', 'beginner'),
('car', 'coche', 'beginner');

-- Insert sample groups
INSERT INTO groups (name) VALUES
('Basic Greetings'),
('Common Phrases'),
('Food and Drink'),
('Transportation'),
('Home and Family');

-- Link words to groups
INSERT INTO word_groups (word_id, group_id) VALUES
(1, 1), -- hello -> Basic Greetings
(2, 1), -- goodbye -> Basic Greetings
(3, 2), -- please -> Common Phrases
(4, 2), -- thank you -> Common Phrases
(5, 1), -- good morning -> Basic Greetings
(6, 1), -- good night -> Basic Greetings
(7, 3), -- water -> Food and Drink
(8, 3), -- food -> Food and Drink
(9, 5), -- house -> Home and Family
(10, 4); -- car -> Transportation

-- Insert sample study activities
INSERT INTO study_activities (name, description) VALUES
('Vocabulary Practice', 'Practice vocabulary words with flashcards'),
('Listening Exercise', 'Listen to native speakers and practice pronunciation'),
('Writing Practice', 'Practice writing sentences using vocabulary words');

-- Link study activities to groups
INSERT INTO study_activity_groups (study_activity_id, group_id) VALUES
(1, 1), -- Vocabulary Practice -> Basic Greetings
(1, 2), -- Vocabulary Practice -> Common Phrases
(2, 1), -- Listening Exercise -> Basic Greetings
(3, 2); -- Writing Practice -> Common Phrases

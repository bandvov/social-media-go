CREATE TABLE reactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,            
    reactable_type VARCHAR(50) NOT NULL, 
    reactable_id INT NOT NULL,        
    reaction_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    FOREIGN KEY (user_id) REFERENCES users(id)
);


-- Create the trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER set_reactions_updated_at
BEFORE UPDATE ON reactions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

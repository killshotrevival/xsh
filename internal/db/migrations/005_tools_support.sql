-- Create Tools table in the database
CREATE TABLE tools (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	connection_string TEXT NOT NULL
);

-- Insert default tools into the tools table
INSERT INTO tools (id, name, connection_string) VALUES ('0a398c97-4525-4416-903d-3662b8de8850', 'SSH', 'Handled Internally By XSH');

-- Add a new column to the HOSTS table to reference the tool used for connecting to the host
ALTER TABLE HOSTS ADD COLUMN tool_id UUID DEFAULT '0a398c97-4525-4416-903d-3662b8de8850';
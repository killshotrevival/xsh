-- Create Hosts table in the database
CREATE TABLE IF NOT EXISTS hosts (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT NOT NULL,
	port INTEGER DEFAULT 22 NOT NULL,
	user TEXT NOT NULL,
	region_id UUID NOT NULL,
	identity_id UUID NOT NULL,
	jumphost_id UUID
);

-- Create Identities table in the database
CREATE TABLE IF NOT EXISTS identities (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	path TEXT NOT NULL
);

-- Create Regions table in the database
CREATE TABLE IF NOT EXISTS regions (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL
);

-- Create Versions table in the database to track applied migrations
CREATE TABLE IF NOT EXISTS schema_version (
		version INTEGER
);

-- Initialize the schema_version table with version 0
INSERT INTO schema_version(version) VALUES(0);

-- Create Tags table in the database
CREATE TABLE IF NOT EXISTS tags (
	id UUID PRIMARY KEY,
	TAG TEXT NOT NULL
);

-- Create TagMappings table in the database to associate tags with hosts, regions, and identities
CREATE TABLE IF NOT EXISTS tagmappings (
	id UUID PRIMARY KEY,
	data_type_id UUID NOT NULL,
	tag_id UUID NOT NULL
);
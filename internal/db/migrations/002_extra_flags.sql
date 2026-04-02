-- Add a new column in hosts table to store extra flags for ssh connection
ALTER TABLE HOSTS ADD COLUMN extra_flags TEXT DEFAULT "-";
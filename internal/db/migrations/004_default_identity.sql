-- This migration adds a default identity to the identities table. 
--This is necessary to ensure that there is always at least one identity available for use in the application.
INSERT INTO identities (id, name, path) VALUES ('be350830-609d-46a0-854c-4ba11e700056', 'Not Required', '/handled/outside/of/xsh');
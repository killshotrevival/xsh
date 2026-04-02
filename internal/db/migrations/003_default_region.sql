-- This migration adds a default region to the regions table. 
--This is necessary to ensure that there is always at least one region available for use in the application.
INSERT INTO regions (id, name) VALUES ('9236a25f-f2ef-418a-aaf5-21fb667e8073', 'Default');
BEGIN;

-- 1. Re-add department column to t_jobs
ALTER TABLE hiring.t_jobs
    ADD COLUMN department VARCHAR;

-- 2. Restore data
UPDATE hiring.t_jobs j
SET department = d.name
FROM hiring.t_departments d
WHERE j.department_id = d.id;

-- 3. Drop department_id column
ALTER TABLE hiring.t_jobs
    DROP COLUMN department_id;

-- 4. Drop t_departments table
DROP TABLE IF EXISTS hiring.t_departments;

COMMIT;

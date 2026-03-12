BEGIN;

-- 1. Create t_departments table
CREATE TABLE IF NOT EXISTS hiring.t_departments (
    id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id UUID NOT NULL, -- soft ref
    name    VARCHAR NOT NULL,
    UNIQUE (team_id, name)
);

-- 2. Add department_id to t_jobs
ALTER TABLE hiring.t_jobs
    ADD COLUMN department_id UUID REFERENCES hiring.t_departments (id) ON DELETE SET NULL;

-- 3. Migrate data
-- For each unique department name per team, insert into t_departments
INSERT INTO hiring.t_departments (team_id, name)
SELECT DISTINCT team_id, department
FROM hiring.t_jobs
WHERE department IS NOT NULL;

-- 4. Update t_jobs with department_id
UPDATE hiring.t_jobs j
SET department_id = d.id
FROM hiring.t_departments d
WHERE j.team_id = d.team_id AND j.department = d.name;

-- 5. Drop old department column
ALTER TABLE hiring.t_jobs
    DROP COLUMN department;

COMMIT;

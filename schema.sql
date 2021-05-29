--
-- Extensions
--
CREATE extension IF NOT EXISTS "uuid-ossp";

--
-- Tables
--
CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(64),
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    last_active TIMESTAMP DEFAULT NOW(),
    email VARCHAR(320) UNIQUE,
    password TEXT,
    type VARCHAR(128) DEFAULT 'GUEST',
    verified BOOL DEFAULT false,
    avatar VARCHAR(128) DEFAULT 'identicon',
    country VARCHAR(2),
    company VARCHAR(256),
    job_title VARCHAR(128)
);

CREATE TABLE IF NOT EXISTS retrospective (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(256),
    owner_id UUID,
    phase SMALLINT NOT NULL DEFAULT 1,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    CONSTRAINT r_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS retrospective_user (
    retrospective_id UUID,
    user_id UUID,
    active BOOL DEFAULT false,
    abandoned BOOL DEFAULT false,
    PRIMARY KEY (retrospective_id, user_id),
    CONSTRAINT ru_retrospective_id_fkey FOREIGN KEY (retrospective_id) REFERENCES retrospective(id) ON DELETE CASCADE,
    CONSTRAINT ru_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS retrospective_item (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    retrospective_id UUID,
    user_id UUID,
    parent_id UUID,
    content TEXT NOT NULL,
    votes JSONB DEFAULT '[]'::JSONB,
    type VARCHAR(16) NOT NULL,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    CONSTRAINT ri_retrospective_id_fkey FOREIGN KEY (retrospective_id) REFERENCES retrospective(id) ON DELETE CASCADE,
    CONSTRAINT ri_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT ri_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES retrospective_item ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS retrospective_action (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    retrospective_id UUID,
    content TEXT NOT NULL,
    completed BOOL DEFAULT false,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    CONSTRAINT ra_retrospective_id_fkey FOREIGN KEY (retrospective_id) REFERENCES retrospective(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_reset (
    reset_id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID,
    created_date TIMESTAMP DEFAULT NOW(),
    expire_date TIMESTAMP DEFAULT NOW() + INTERVAL '1 hour',
    CONSTRAINT ur_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_verify (
    verify_id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID,
    created_date TIMESTAMP DEFAULT NOW(),
    expire_date TIMESTAMP DEFAULT NOW() + INTERVAL '24 hour',
    CONSTRAINT uv_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS api_keys (
    id TEXT NOT NULL PRIMARY KEY,
    user_id UUID,
    name VARCHAR(256) NOT NULL,
    active BOOL DEFAULT true,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, name),
    CONSTRAINT apk_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization (
    id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(256),
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS organization_user (
    organization_id UUID,
    user_id UUID,
    role VARCHAR(16) NOT NULL DEFAULT 'MEMBER',
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (organization_id, user_id),
    CONSTRAINT ou_organization_id FOREIGN KEY(organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT ou_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization_department (
    id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    organization_id UUID,
    name VARCHAR(256),
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    UNIQUE(organization_id, name),
    CONSTRAINT od_organization_id FOREIGN KEY(organization_id) REFERENCES organization(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS department_user (
    department_id UUID,
    user_id UUID,
    role VARCHAR(16) NOT NULL DEFAULT 'MEMBER',
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (department_id, user_id),
    CONSTRAINT du_department_id FOREIGN KEY(department_id) REFERENCES organization_department(id) ON DELETE CASCADE,
    CONSTRAINT du_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS team (
    id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(256),
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_user (
    team_id UUID,
    user_id UUID,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    role VARCHAR(16) NOT NULL DEFAULT 'MEMBER',
    PRIMARY KEY (team_id, user_id),
    CONSTRAINT tu_team_id FOREIGN KEY(team_id) REFERENCES team(id) ON DELETE CASCADE,
    CONSTRAINT tu_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization_team (
    organization_id UUID,
    team_id UUID,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (organization_id, team_id),
    UNIQUE(team_id),
    CONSTRAINT ot_organization_id FOREIGN KEY(organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT ot_team_id FOREIGN KEY(team_id) REFERENCES team(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS department_team (
    department_id UUID,
    team_id UUID,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (department_id, team_id),
    UNIQUE(team_id),
    CONSTRAINT dt_department_id FOREIGN KEY(department_id) REFERENCES organization_department(id) ON DELETE CASCADE,
    CONSTRAINT dt_team_id FOREIGN KEY(team_id) REFERENCES team(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS team_retrospective (
    team_id UUID,
    retrospective_id UUID,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (team_id, retrospective_id),
    CONSTRAINT tb_team_id FOREIGN KEY(team_id) REFERENCES team(id) ON DELETE CASCADE,
    CONSTRAINT tb_retrospective_id FOREIGN KEY(retrospective_id) REFERENCES retrospective(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS alert (
    id UUID  NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    type VARCHAR(128) DEFAULT 'NEW',
    content TEXT NOT NULL,
    active BOOLEAN DEFAULT true,
    allow_dismiss BOOLEAN DEFAULT true,
    registered_only BOOLEAN DEFAULT true,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW()
);

--
-- Table Alterations
--

--
-- Views
--
CREATE MATERIALIZED VIEW IF NOT EXISTS active_countries AS SELECT DISTINCT country FROM users;

--
-- Stored Procedures
--

-- Reset All Users to Inactive, used by server restart --
CREATE OR REPLACE PROCEDURE deactivate_all_users()
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retrospective_user SET active = false WHERE active = true;
END;
$$;

-- Set Retrospective Owner --
CREATE OR REPLACE PROCEDURE set_retrospective_owner(retrospectiveId UUID, ownerId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retrospective SET updated_date = NOW(), owner_id = ownerId WHERE id = retrospectiveId;
END;
$$;

-- Set Retrospective Phase --
CREATE OR REPLACE PROCEDURE set_retrospective_phase(retrospectiveId UUID, nextPhase SMALLINT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retrospective SET updated_date = NOW(), phase = nextPhase WHERE id = retrospectiveId;
END;
$$;

-- Delete Retrospective --
CREATE OR REPLACE PROCEDURE delete_retrospective(retrospectiveId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM retrospective WHERE id = retrospectiveId;

    COMMIT;
END;
$$;

-- Reset User Password --
CREATE OR REPLACE PROCEDURE reset_user_password(resetId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM user_reset ur
        LEFT JOIN user w ON w.id = ur.user_id
        WHERE ur.reset_id = resetId AND NOW() < ur.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete incase reset record expired
        DELETE FROM user_reset WHERE reset_id = resetId;
        RAISE 'Valid Reset ID not found';
    END IF;

    UPDATE users SET password = userPassword, last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM user_reset WHERE reset_id = resetId;

    COMMIT;
END;
$$;

-- Update User Password --
CREATE OR REPLACE PROCEDURE update_user_password(userId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET password = userPassword, last_active = NOW(), updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Verify a user account email
CREATE OR REPLACE PROCEDURE verify_user_account(verifyId UUID)
LANGUAGE plpgsql AS $$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT usr.id
        FROM user_verify uv
        LEFT JOIN users usr ON usr.id = uv.user_id
        WHERE uv.verify_id = verifyId AND NOW() < uv.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete incase verify record expired
        DELETE FROM user_verify WHERE verify_id = verifyId;
        RAISE 'Valid Verify ID not found';
    END IF;

    UPDATE users SET verified = 'TRUE', last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM user_verify WHERE verify_id = verifyId;

    COMMIT;
END;
$$;

-- Promote User to ADMIN by ID --
CREATE OR REPLACE PROCEDURE promote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'ADMIN', updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Promote User to ADMIN by Email --
CREATE OR REPLACE PROCEDURE promote_user_by_email(userEmail VARCHAR(320))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'ADMIN', updated_date = NOW() WHERE email = userEmail;

    COMMIT;
END;
$$;

-- Demote User to Registered by ID --
CREATE OR REPLACE PROCEDURE demote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'REGISTERED', updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Clean up Retrospectives older than X Days --
CREATE OR REPLACE PROCEDURE clean_retrospectives(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM retrospective WHERE updated_date < (NOW() - daysOld * interval '1 day');

    COMMIT;
END;
$$;

-- Clean up Guest Users (and their created retrospectives) older than X Days --
CREATE OR REPLACE PROCEDURE clean_guest_users(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE last_active < (NOW() - daysOld * interval '1 day') AND type = 'GUEST';
    REFRESH MATERIALIZED VIEW active_countries;

    COMMIT;
END;
$$;

-- Deletes a User and all his retrospective(s), api keys --
CREATE OR REPLACE PROCEDURE delete_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE id = userId;
    REFRESH MATERIALIZED VIEW active_countries;

    COMMIT;
END;
$$;

-- Updates a users profile --
CREATE OR REPLACE PROCEDURE user_profile_update(
    userId UUID,
    userName VARCHAR(64),
    userAvatar VARCHAR(128),
    userCountry VARCHAR(2),
    userCompany VARCHAR(256),
    userJobTitle VARCHAR(128)
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users
    SET name = userName, avatar = userAvatar, country = userCountry, company = userCompany, job_title = userJobTitle, last_active = NOW(), updated_date = NOW()
    WHERE id = userId;
    REFRESH MATERIALIZED VIEW active_countries;
END;
$$;

--
-- Stored Functions
--

-- Create a Retrospective
DROP FUNCTION IF EXISTS create_retrospective(UUID, VARCHAR);
CREATE FUNCTION create_retrospective(ownerId UUID, retrospectiveName VARCHAR(256)) RETURNS UUID 
AS $$ 
DECLARE retroId UUID;
BEGIN
    INSERT INTO retrospective (owner_id, name) VALUES (ownerId, retrospectiveName) RETURNING id INTO retroId;

    RETURN retroId;
END;
$$ LANGUAGE plpgsql;

-- Get Retrospectives by User ID
DROP FUNCTION IF EXISTS get_retrospectives_by_user(uuid);
CREATE FUNCTION get_retrospectives_by_user(userId UUID) RETURNS table (
    id UUID, name VARCHAR(256), owner_id UUID, phase SMALLINT
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name, b.owner_id, b.phase
		FROM retrospective b
		LEFT JOIN retrospective_user su ON b.id = su.retrospective_id WHERE su.user_id = userId AND su.abandoned = false
		GROUP BY b.id ORDER BY b.created_date DESC;
END;
$$ LANGUAGE plpgsql;

-- Get a User by ID
DROP FUNCTION IF EXISTS get_user(UUID);
CREATE FUNCTION get_user(userId UUID) RETURNS table (
    id UUID, name VARCHAR(64), email VARCHAR(320), type VARCHAR(128), verified BOOL, avatar VARCHAR(128), country VARCHAR(2), company VARCHAR(256), jobTitle VARCHAR(128)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, coalesce(u.email, ''), u.type, u.verified, u.avatar, u.country, u.company, u.job_title FROM users u WHERE u.id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Retrospective Users
DROP FUNCTION IF EXISTS get_retrospective_users(uuid);
CREATE FUNCTION get_retrospective_users(retrospectiveId UUID) RETURNS table (
    id UUID, name VARCHAR(256), active BOOL
) AS $$
BEGIN
    RETURN QUERY
        SELECT
			w.id, w.name, su.active
		FROM retrospective_user su
		LEFT JOIN users w ON su.user_id = w.id
		WHERE su.retrospective_id = retrospectiveId
		ORDER BY w.name;
END;
$$ LANGUAGE plpgsql;

-- Get Retrospective User by id
DROP FUNCTION IF EXISTS get_retrospective_user(uuid, uuid);
CREATE FUNCTION get_retrospective_user(retrospectiveId UUID, userId UUID) RETURNS table (
    id UUID, name VARCHAR(256), active BOOL
) AS $$
BEGIN
    RETURN QUERY
        SELECT
			w.id, w.name, coalesce(su.active, FALSE)
		FROM users w
		LEFT JOIN retrospective_user su ON su.user_id = w.id AND su.retrospective_id = retrospectiveId
		WHERE w.id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get User Auth by Email
DROP FUNCTION IF EXISTS get_user_auth_by_email(VARCHAR);
CREATE FUNCTION get_user_auth_by_email(userEmail VARCHAR(320)) RETURNS table (
    id UUID, name VARCHAR(64), email VARCHAR(320), type VARCHAR(128), password TEXT
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, coalesce(u.email, ''), u.type, u.password FROM users u WHERE u.email = userEmail;
END;
$$ LANGUAGE plpgsql;

-- Get Application Stats e.g. total user and retrospective counts
DROP FUNCTION IF EXISTS get_app_stats();
DROP FUNCTION IF EXISTS get_app_stats(
    OUT unregistered_user_count INTEGER,
    OUT registered_user_count INTEGER,
    OUT retrospective_count INTEGER
);
DROP FUNCTION IF EXISTS get_app_stats(
    OUT unregistered_user_count INTEGER,
    OUT registered_user_count INTEGER,
    OUT retrospective_count INTEGER,
    OUT organization_count INTEGER,
    OUT department_count INTEGER,
    OUT team_count INTEGER
);
CREATE FUNCTION get_app_stats(
    OUT unregistered_user_count INTEGER,
    OUT registered_user_count INTEGER,
    OUT retrospective_count INTEGER,
    OUT organization_count INTEGER,
    OUT department_count INTEGER,
    OUT team_count INTEGER,
    OUT apikey_count INTEGER
) AS $$
BEGIN
    SELECT COUNT(*) INTO unregistered_user_count FROM users WHERE email IS NULL;
    SELECT COUNT(*) INTO registered_user_count FROM users WHERE email IS NOT NULL;
    SELECT COUNT(*) INTO retrospective_count FROM retrospective;
    SELECT COUNT(*) INTO organization_count FROM organization;
    SELECT COUNT(*) INTO department_count FROM organization_department;
    SELECT COUNT(*) INTO team_count FROM team;
    SELECT COUNT(*) INTO apikey_count FROM api_keys;
END;
$$ LANGUAGE plpgsql;

-- Insert a new user password reset
DROP FUNCTION IF EXISTS insert_user_reset(VARCHAR);
CREATE FUNCTION insert_user_reset(
    IN userEmail VARCHAR(320),
    OUT resetId UUID,
    OUT userId UUID,
    OUT userName VARCHAR(64)
)
AS $$ 
BEGIN
    SELECT id, name INTO userId, userName FROM users WHERE email = userEmail;
    IF FOUND THEN
        INSERT INTO user_reset (user_id) VALUES (userId) RETURNING reset_id INTO resetId;
    ELSE
        RAISE EXCEPTION 'Nonexistent User --> %', userEmail USING HINT = 'Please check your Email';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Register a new user
DROP FUNCTION IF EXISTS register_user(VARCHAR, VARCHAR, TEXT, VARCHAR);
CREATE FUNCTION register_user(
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
)
AS $$
BEGIN
    INSERT INTO users (name, email, password, type)
    VALUES (userName, userEmail, hashedPassword, userType)
    RETURNING id INTO userId;

    INSERT INTO user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Register a new user from existing GUEST
DROP FUNCTION IF EXISTS register_existing_user(UUID, VARCHAR, VARCHAR, TEXT, VARCHAR);
CREATE FUNCTION register_existing_user(
    IN activeUserId UUID,
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
)
AS $$
BEGIN
    UPDATE users
    SET
         name = userName,
         email = userEmail,
         password = hashedPassword,
         type = userType,
         last_active = NOW(),
         updated_date = NOW()
    WHERE id = activeUserId
    RETURNING id INTO userId;

    INSERT INTO user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Get a list of countries
CREATE OR REPLACE FUNCTION countries_active() RETURNS table (
    country VARCHAR(2)
) AS $$
BEGIN
    RETURN QUERY SELECT ac.country FROM active_countries ac;
END;
$$ LANGUAGE plpgsql;

--
-- ORGANIZATIONS --
--

-- Get Organization --
CREATE OR REPLACE FUNCTION organization_get_by_id(
    IN orgId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date
        FROM organization o
        WHERE o.id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization User Role --
CREATE OR REPLACE FUNCTION organization_get_user_role(
    IN userId UUID,
    IN orgId UUID,
    OUT role VARCHAR(16)
) AS $$
BEGIN
    SELECT ou.role INTO role
    FROM organization_user ou
    WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Organizations --
CREATE OR REPLACE FUNCTION organization_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date
        FROM organization o
        ORDER BY o.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Organizations by User --
CREATE OR REPLACE FUNCTION organization_list_by_user(
    IN userId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP, role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date, ou.role
        FROM organization_user ou
        LEFT JOIN organization o ON ou.organization_id = o.id
        WHERE ou.user_id = userId
        ORDER BY created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Organization --
CREATE OR REPLACE FUNCTION organization_create(
    IN userId UUID,
    IN orgName VARCHAR(256),
    OUT organizationId UUID
) AS $$
BEGIN
    INSERT INTO organization (name) VALUES (orgName) RETURNING id INTO organizationId;
    INSERT INTO organization_user (organization_id, user_id, role) VALUES (organizationId, userId, 'ADMIN');
END;
$$ LANGUAGE plpgsql;

-- Add User to Organization --
CREATE OR REPLACE FUNCTION organization_user_add(
    IN orgId UUID,
    IN userId UUID,
    IN userRole VARCHAR(16)
) RETURNS void AS $$
BEGIN
    INSERT INTO organization_user (organization_id, user_id, role) VALUES (orgId, userId, userRole);
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Remove User from Organization --
CREATE OR REPLACE PROCEDURE organization_user_remove(orgId UUID, userId UUID)
AS $$
DECLARE temprow record;
BEGIN
    FOR temprow IN
        SELECT id FROM organization_department WHERE organization_id = orgId
    LOOP
        CALL department_user_remove(temprow.id, userId);
    END LOOP;
    DELETE FROM team_user tu WHERE tu.team_id IN (
        SELECT ot.team_id
        FROM organization_team ot
        WHERE ot.organization_id = orgId
    ) AND tu.user_id = userId;
    DELETE FROM organization_user WHERE organization_id = orgId AND user_id = userId;
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;

    COMMIT;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Users --
CREATE OR REPLACE FUNCTION organization_user_list(
    IN orgId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, u.email, ou.role
        FROM organization_user ou
        LEFT JOIN users u ON ou.user_id = u.id
        WHERE ou.organization_id = orgId
        ORDER BY ou.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Teams --
CREATE OR REPLACE FUNCTION organization_team_list(
    IN orgId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM organization_team ot
        LEFT JOIN team t ON ot.team_id = t.id
        WHERE ot.organization_id = orgId
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Organization Team --
CREATE OR REPLACE FUNCTION organization_team_create(
    IN orgId UUID,
    IN teamName VARCHAR(256),
    OUT teamId UUID
) AS $$
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING id INTO teamId;
    INSERT INTO organization_team (organization_id, team_id) VALUES (orgId, teamId);
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Team User Role --
CREATE OR REPLACE FUNCTION organization_team_user_role(
    IN userId UUID,
    IN orgId UUID,
    IN teamId UUID
) RETURNS table (
    orgRole VARCHAR(16), teamRole VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT ou.role AS orgRole, COALESCE(tu.role, '') AS teamRole
        FROM organization_user ou
        LEFT JOIN team_user tu ON tu.user_id = userId AND tu.team_id = teamId
        WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

--
-- DEPARTMENTS --
--

-- Get Department --
CREATE OR REPLACE FUNCTION department_get_by_id(
    IN departmentId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT od.id, od.name, od.created_date, od.updated_date
        FROM organization_department od
        WHERE od.id = departmentId;
END;
$$ LANGUAGE plpgsql;

-- Get Department User Role --
CREATE OR REPLACE FUNCTION department_get_user_role(
    IN userId UUID,
    IN orgId UUID,
    IN departmentId UUID
) RETURNS table (
    orgRole VARCHAR(16), departmentRole VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT ou.role AS orgRole, COALESCE(du.role, '') AS departmentRole
        FROM organization_user ou
        LEFT JOIN department_user du ON du.user_id = userId AND du.department_id = departmentId
        WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Departments --
CREATE OR REPLACE FUNCTION department_list(
    IN orgId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT d.id, d.name, d.created_date, d.updated_date
        FROM organization_department d
        WHERE d.organization_id = orgId
        ORDER BY d.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Organization Department --
CREATE OR REPLACE FUNCTION department_create(
    IN orgId UUID,
    IN departmentName VARCHAR(256),
    OUT departmentId UUID
) AS $$
BEGIN
    INSERT INTO organization_department (name, organization_id) VALUES (departmentName, orgId) RETURNING id INTO departmentId;
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Get Department Teams --
CREATE OR REPLACE FUNCTION department_team_list(
    IN departmentId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM department_team dt
        LEFT JOIN team t ON dt.team_id = t.id
        WHERE dt.department_id = departmentId
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Department Team --
CREATE OR REPLACE FUNCTION department_team_create(
    IN departmentId UUID,
    IN teamName VARCHAR(256),
    OUT teamId UUID
) AS $$
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING id INTO teamId;
    INSERT INTO department_team (department_id, team_id) VALUES (departmentId, teamId);
    UPDATE organization_department SET updated_date = NOW() WHERE id = departmentId;
END;
$$ LANGUAGE plpgsql;

-- Get Department Team User Role --
CREATE OR REPLACE FUNCTION department_team_user_role(
    IN userId UUID,
    IN orgId UUID,
    IN departmentId UUID,
    IN teamId UUID
) RETURNS table (
    orgRole VARCHAR(16), departmentRole VARCHAR(16), teamRole VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT ou.role AS orgRole, COALESCE(du.role, '') AS departmentRole, COALESCE(tu.role, '') AS teamRole
        FROM organization_user ou
        LEFT JOIN department_user du ON du.user_id = userId AND du.department_id = departmentId
        LEFT JOIN team_user tu ON tu.user_id = userId AND tu.team_id = teamId
        WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Department Users --
CREATE OR REPLACE FUNCTION department_user_list(
    IN departmentId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, u.email, du.role
        FROM department_user du
        LEFT JOIN users u ON du.user_id = u.id
        WHERE du.department_id = departmentId
        ORDER BY du.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Add User to Department --
CREATE OR REPLACE FUNCTION department_user_add(
    IN departmentId UUID,
    IN userId UUID,
    IN userRole VARCHAR(16)
) RETURNS void AS $$
DECLARE orgId UUID;
BEGIN    
    SELECT organization_id INTO orgId FROM organization_user WHERE user_id = userId;

    IF orgId IS NULL THEN
        RAISE EXCEPTION 'User not in Organization -> %', userId USING HINT = 'Please add user to Organization before department';
    END IF;

    INSERT INTO department_user (department_id, user_id, role) VALUES (departmentId, userId, userRole);
    UPDATE organization_department SET updated_date = NOW() WHERE id = departmentId;
END;
$$ LANGUAGE plpgsql;

-- Remove User from Department --
CREATE OR REPLACE PROCEDURE department_user_remove(departmentId UUID, userId UUID)
AS $$
BEGIN
    DELETE FROM team_user tu WHERE tu.team_id IN (
        SELECT dt.team_id
        FROM department_team dt
        WHERE dt.department_id = departmentId
    ) AND tu.user_id = userId;
    DELETE FROM department_user WHERE department_id = departmentId AND user_id = userId;
    UPDATE organization_department SET updated_date = NOW() WHERE id = departmentId;

    COMMIT;
END;
$$ LANGUAGE plpgsql;

--
-- TEAMS --
--

-- Get Team --
CREATE OR REPLACE FUNCTION team_get_by_id(
    IN teamId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date
        FROM team o
        WHERE o.id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Get Team User Role --
CREATE OR REPLACE FUNCTION team_get_user_role(
    IN userId UUID,
    IN teamId UUID
) RETURNS table (
    role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT tu.role
        FROM team_user tu
        WHERE tu.team_id = teamId AND tu.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Teams --
CREATE OR REPLACE FUNCTION team_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM team t
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Teams by User --
CREATE OR REPLACE FUNCTION team_list_by_user(
    IN userId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM team_user tu
        LEFT JOIN team t ON tu.team_id = t.id
        WHERE tu.user_id = userId
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Team --
CREATE OR REPLACE FUNCTION team_create(
    IN userId UUID,
    IN teamName VARCHAR(256),
    OUT teamId UUID
) AS $$
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING id INTO teamId;
    INSERT INTO team_user (team_id, user_id, role) VALUES (teamId, userId, 'ADMIN');
END;
$$ LANGUAGE plpgsql;

-- Get Team Users --
CREATE OR REPLACE FUNCTION team_user_list(
    IN teamId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, u.email, tu.role
        FROM team_user tu
        LEFT JOIN users u ON tu.user_id = u.id
        WHERE tu.team_id = teamId
        ORDER BY tu.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Add User to Team --
CREATE OR REPLACE FUNCTION team_user_add(
    IN teamId UUID,
    IN userId UUID,
    IN userRole VARCHAR(16)
) RETURNS void AS $$
BEGIN
    INSERT INTO team_user (team_id, user_id, role) VALUES (teamId, userId, userRole);
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Remove User from Team --
CREATE OR REPLACE PROCEDURE team_user_remove(teamId UUID, userId UUID)
AS $$
BEGIN
    DELETE FROM team_user WHERE team_id = teamId AND user_id = userId;
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Get Team Retrospectives --
CREATE OR REPLACE FUNCTION team_retrospective_list(
    IN teamId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256)
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name
        FROM team_retrospective tb
        LEFT JOIN retrospective b ON tb.retrospective_id = b.id
        WHERE tb.team_id = teamId
        ORDER BY tb.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Add Retrospective to Team --
CREATE OR REPLACE FUNCTION team_retrospective_add(
    IN teamId UUID,
    IN retrospectiveId UUID
) RETURNS void AS $$
BEGIN
    INSERT INTO team_retrospective (team_id, retrospective_id) VALUES (teamId, retrospectiveId);
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Remove Retrospective from Team --
CREATE OR REPLACE FUNCTION team_retrospective_remove(
    IN teamId UUID,
    IN retrospectiveId UUID
) RETURNS void AS $$
BEGIN
    DELETE FROM team_retrospective WHERE retrospective_id = retrospectiveId AND team_id = teamId;
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Delete Team --
CREATE OR REPLACE PROCEDURE team_delete(teamId UUID)
AS $$
BEGIN
    DELETE FROM team WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;
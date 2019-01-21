CREATE TABLE "user" (
    handle character varying(128) PRIMARY KEY,
    name character varying(128) NOT NULL,
    login_time timestamp with time zone,
    login_duration bigint,
    created timestamp with time zone NOT NULL DEFAULT now(),
    modified timestamp with time zone,
    ip bytea,
    password bytea NOT NULL,
    enabled bool DEFAULT false
);

CREATE OR REPLACE FUNCTION update_modified_column()	
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified = now();
    RETURN NEW;	
END;
$$ language 'plpgsql';

-- use BEFORE UPDATE if loop
CREATE TRIGGER update_customer_modtime AFTER UPDATE ON "user" FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

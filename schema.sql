CREATE SEQUENCE 
IF NOT EXISTS next_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS individuals (
    id text DEFAULT nextval('next_id') PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    -- cfb_name is the name formatted to match CFB
    -- filings in "Last, First" format
    cfb_name text NOT NULL,
    zip text,
    updated_ts timestamp NOT NULL,
    role text,
    title text,
    twitter text
    -- TODO: add other fields
);

CREATE TABLE IF NOT EXISTS associations (
    id text DEFAULT nextval('next_id') PRIMARY KEY,
    description text NOT NULL,
    category_id text
);

CREATE TABLE IF NOT EXISTS categories (
    id text DEFAULT nextval('next_id') PRIMARY KEY,
    description text NOT NULL
);

CREATE TABLE IF NOT EXISTS individual_associations (
    individual_id text NOT NULL,
    association_id text NOT NULL,
    updated_ts timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS contributions (
    id text DEFAULT nextval('next_id') PRIMARY KEY,
    -- seqno orders contributions by date in ascending order
    -- at time of ingestion
    seqno integer NOT NULL,
    -- amount denominated in cents
    amount integer NOT NULL,
    date timestamp NOT NULL,
    -- contributor_name as reported in CFB filing
    contributor_name text NOT NULL,
    -- optimistically match the contributor to 
    -- an individual
    contributor_id text,
    -- recipient_name as reported in CFB filing
    recipient_name text NOT NULL,
    -- All recipients should have annotations, so we will
    -- always have individual IDs for recipients
    recipient_id text NOT NULL,

    -- Other contribution data from CFB
    address text NOT NULL,
    employer text NOT NULL,
    occupation text NOT NULL,
    type text NOT NULL
);

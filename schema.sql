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
    -- TODO: make a search_name column (?)
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
    -- reference number from CFB
    refno text NOT NULL,
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
    recipient_id text,
    -- CFB gives recipient a recipient ID, so note this.
    cfb_recipient_id text NOT NULL,

    -- Other contribution data from CFB
    election TEXT NOT NULL,
    office_cd TEXT,
    can_class TEXT,
    committee TEXT,
    filing INT,
    schedule TEXT,
    c_code TEXT,

    borough text,
    city text NOT NULL,
    state text NOT NULL,
    zip text NOT NULL,
    occupation text,
    employer_name text
);

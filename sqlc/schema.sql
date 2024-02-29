CREATE TABLE clients (
  id integer PRIMARY KEY AUTOINCREMENT,
  client_id text,
  username text,
  password blob,
  subscribe_acl text,
  publish_acl text,
  created_at integer
);
CREATE UNIQUE INDEX client_cred_idx ON clients(client_id,username);
CREATE TABLE retries (
  id integer PRIMARY KEY AUTOINCREMENT,
  client_id text,
  application_message blob,
  packet_identifier integer,
  qos integer,
  dup integer,
  retries integer,
  ack_status integer,
  created_at integer,
  session_id integer,
  reason_code integer
);
CREATE UNIQUE INDEX client_identifier_idx ON retries(client_id,packet_identifier);
CREATE TABLE retains (
  id integer PRIMARY KEY AUTOINCREMENT,
  client_id text,
  topic text,
  application_message blob,
  created_at integer
);
CREATE UNIQUE INDEX retain_unique_client_topic_idx ON retains(client_id, topic);
CREATE TABLE sessions (
  id integer PRIMARY KEY AUTOINCREMENT,
  last_seen integer,
  last_connect integer,
  expiry_interval integer,
  client_id text,
  connected integer,
  protocol_version integer
);
CREATE UNIQUE INDEX client_unique_session_idx ON sessions(client_id);
CREATE TABLE subscriptions (
  id integer PRIMARY KEY AUTOINCREMENT,
  client_id text,
  topic text,
  retain_handling integer,
  retain_as_published integer,
  no_local integer,
  qos integer,
  protocol_version integer,
  enabled integer,
  created_at integer,
  session_id integer,
  shared integer DEFAULT false,
  share_name text,
  CONSTRAINT fk_sessions_subscriptions FOREIGN KEY (session_id) REFERENCES sessions(id)
);
CREATE UNIQUE INDEX sub_pars_idx ON subscriptions(client_id,topic,share_name);
CREATE TABLE users (
  id integer PRIMARY KEY AUTOINCREMENT,
  username text,
  password blob,
  created_at integer
);
CREATE UNIQUE INDEX username_user_idx ON users(username);
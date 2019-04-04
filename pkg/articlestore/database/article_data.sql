/**Should map id to article_id in article_data **/
CREATE TABLE user_markers (
  id SERIAL PRIMARY KEY,
  user_id int NOT NULL,
  hovered_over int, 
  generated int, 
  clicked int, 
  searched int,
  article_interaction int, 
  created_at TIME,
  updated_at TIME,
  deleted_at TIME
);

CREATE TABLE markers (
  id SERIAL PRIMARY KEY,
  url TEXT,
  info TEXT,
  title TEXT,
  lat double precision,
  lon double precision,
  source TEXT,
  generated int,
  beg_year int NOT NULL DEFAULT -3000,
  end_year int NOT NULL DEFAULT 2019, 
  hovered_over int,
  clicked int,
  searched int,
  created_at TIME,
  updated_at TIME,
);
/**Should map id to article_id in article_data **/
CREATE TABLE user_article_url (
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

CREATE TABLE article_url (
  id SERIAL PRIMARY KEY,
  url TEXT,
  title TEXT,
  lat double precision,
  lon double precision,
  source TEXT,
  generated int,
  beg_year int,
  end_year int, 
  hovered_over int,
  clicked int,
  searched int,
  created_at TIME,
  updated_at TIME,
);

/**Include stores for change in article content.**/
CREATE TABLE article_text (
  id SERIAL PRIMARY KEY, 
  title TEXT, 
  lat double precision, 
  lon double precision,
  generated int,
  beg_year int, 
  end_year int,
  hovered_over int, 
  clicked int, 
  searched int, 
  created_at TIME, 
  updated_at TIME, 
);

/**Should map foreign key to article id number of article_text**/
CREATE TABLE user_article_text (
  id SERIAL PRIMARY KEY, 
  user_id int NOT NULL,
  hovered_over int, 
  generated int, 
  clicked int, 
  searched int,
  source TEXT,
  article_interaction int, 
  created_at TIME,
  updated_at TIME,
  deleted_at TIME
);

CREATE EXTENSION postgis;
/*https://stackoverflow.com/questions/24981784/how-do-i-add-postgis-to-postgresql-pgadmin */

CREATE TABLE markers (
    id BIGINT PRIMARY KEY, 
    info TEXT, 
    url TEXT, 
    title TEXT, 
    lat double precision, 
    lon double precision,
    /*Include field to indicate which source store this information comes from.*/
    source TEXT, 
    beg_year int, 
    end_year int, 
)

SELECT AddGeometryColumn('markers', 'geom', 4326, 'POINT', 2)


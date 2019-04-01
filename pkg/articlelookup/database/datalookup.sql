CREATE EXTENSION postgis;
/*https://stackoverflow.com/questions/24981784/how-do-i-add-postgis-to-postgresql-pgadmin */

CREATE TABLE markers_url (
    id BIGINT PRIMARY KEY, 
    url TEXT, 
    title TEXT, 
    /*Include field to indicate which source store this information comes from.*/
    source TEXT
)

SELECT AddGeometryColumn('markers', 'geom', 4326, 'POINT', 2)


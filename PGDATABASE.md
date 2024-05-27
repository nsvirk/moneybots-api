# Setup postgres database on server

## Install postgres

## Setup main and dev database

```sh
# get psql prompt locally on server without password for postgres user
psql postgresql://postgres:@localhost:5432/postgres

# create role, database and assign privileges in #postgres prompt
CREATE ROLE mb_user WITH LOGIN ENCRYPTED PASSWORD 'your_password_here';
ALTER ROLE mb_user WITH SUPERUSER;
CREATE DATABASE moneybots;
GRANT ALL PRIVILEGES ON DATABASE moneybots to mb_user;

# create dev database and dev user
CREATE ROLE mb_dev WITH LOGIN ENCRYPTED PASSWORD 'your_password_here';
ALTER ROLE mb_dev WITH SUPERUSER;
CREATE DATABASE moneybots_dev;
GRANT ALL PRIVILEGES ON DATABASE moneybots_dev to mb_dev;

```

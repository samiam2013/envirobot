# envirobot

a simple web server and set of go routines for collecting and presenting environmental sensor data

## install the service with 
`sudo ln -s $(pwd)/envirobot.service /etc/systemd/system/`
`sudo systemctl enable envirobot`

## set up the sqlite database
`sqlite3 sqlite.db < migrate_up.sql`

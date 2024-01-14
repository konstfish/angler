#!/bin/sh

REPO="P3TERX/GeoLite.mmdb"
FILE="GeoLite2-City.mmdb"

wget https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-City.mmdb

exec ./main

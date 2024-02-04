#!/bin/sh
set -e

mkdir -p db

REPO="P3TERX/GeoLite.mmdb"
FILE="GeoLite2-City.mmdb"

if [ ! -f "$FILE" ] || find "$FILE" -mtime +7 -print | grep -q .; then
    echo "Redownloading database: $FILE from $REPO"
    wget -O "$FILE" "https://github.com/$REPO/releases/latest/download/$FILE" || { echo "Download failed"; exit 1; }
    echo "Download successful"
else
    echo "Database is up-to-date"
fi

exec ./main
#!/bin/sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
DASHBOARD="$ROOT_DIR/start.html"

open_link() {
  url="$1"
  open "$url"
  sleep 1
}

ask_open() {
  name="$1"
  url="$2"
  printf "Åbn %s? Skriv j for ja: " "$name"
  read answer
  case "$answer" in
    j|J) open_link "$url" ;;
    *) printf "Springer over: %s\n" "$name" ;;
  esac
}

printf "\n== GF2 IT setup-assistent til Mac ==\n"
printf "Assistenten gemmer ingen brugernavne, adgangskoder eller MitID-oplysninger.\n"

if [ -f "$DASHBOARD" ]; then
  printf "\nÅbner dashboardet lokalt.\n"
  open "$DASHBOARD"
else
  printf "\nDashboardet blev ikke fundet ved: %s\n" "$DASHBOARD"
fi

printf "\n== Åbn skole- og programlinks ==\n"
ask_open "Office 365 / skolemail" "https://www.office.com/"
ask_open "Moodle" "https://online.neg.dk/login/index.php"
ask_open "Lectio" "https://www.lectio.dk/lectio/769/default.aspx"
ask_open "PraxisOnline" "https://online.praxis.dk/"
ask_open "OneDrive" "https://www.office.com/launch/onedrive"
ask_open "SketchUp / Trimble" "https://sketchup.trimble.com/"
ask_open "NEG hjemmeside" "https://www.neg.dk/"

printf "\n== SketchUp ==\n"
printf "Hvis SketchUp ikke er installeret, brug SketchUp-linket og følg skolens Trimble-flow fra skolemailen.\n"
printf "Dashboard-genvej på Mac holdes manuel i første version: læg mappen et fast sted og åbn start.html.\n"

printf "\nFærdig. Tryk Enter for at lukke.\n"
read unused

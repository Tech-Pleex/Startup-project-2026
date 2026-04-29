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
ask_open "PraxisOnline" "https://authentication.praxis.dk/Account/Login?ReturnUrl=%2Fconnect%2Fauthorize%3Fclient_id%3DPraxisOnlinev2%26redirect_uri%3Dhttps%253A%252F%252Fonline.praxis.dk%252Fauthentication%252Flogin-callback%26response_type%3Dcode%26scope%3Dopenid%2520profile%2520PraxisOnlineClient%26state%3Debe8a1c6ff5f4f98b3014db0c5dc752d%26code_challenge%3DEa2At8GN59IETq2ud1CQuReFA7oUdSLXkB58eploqic%26code_challenge_method%3DS256%26response_mode%3Dquery"
ask_open "OneDrive" "https://www.office.com/launch/onedrive"
ask_open "SketchUp / Trimble" "https://id.trimble.com/ui/sign_in.html?state=eyJhbGciOiJSUzI1NiIsImtpZCI6IjIiLCJ0eXAiOiJKV1QifQ.eyJvYXV0aF9wYXJhbWV0ZXJzIjp7ImNsaWVudF9pZCI6ImNiMzg4Yzk2LTY2YjUtNDdhMS04MzZmLWFlYzQ0YTdmMGJjYSIsInJlZGlyZWN0X3VyaSI6Imh0dHBzOi8vd3d3LnRyaW1ibGUuY29tL2xvZ2luIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJzY29wZSI6Im9wZW5pZCBpYW0gdHJpbWJsZS1teHAtbG9naW4gVENNaWRkbGV3YXJlIERYLVRyaWFscy1BcHAiLCJzdGF0ZSI6Ii9lbiJ9LCJleHRyYV9wYXJhbWV0ZXJzIjp7fSwiaW50ZXJuYWxfcGFyYW1ldGVycyI6eyJzZW5kX2FjY291bnRfaWRfaW5fY2xhaW1zIjpmYWxzZSwiaXNfaW50ZXJuYWwiOnRydWV9LCJleHAiOiIyMDI2LTA0LTI5IDExOjQ2OjMyLjc4MzIzMCIsIm5iZiI6MTc3NzQ2MjU5MiwiZXhwVHMiOjE3Nzc0NjMxOTIsInJlcV9leHAiOiIyMDI2LTA0LTI5IDExOjM4OjMyLjc4MzI1NyIsInRjcF9yZXF1ZXN0X2lkIjoiOGYxZWI1ZWY2OGU3NGI4MmFhZTdkN2FhM2I3NmRjNmUiLCJjb3JyZWxhdGlvbl9pZCI6IjNkZWE0OWZjZWUxZjQyOGU4ZThhMWZmZGI2MTg2NTA5XzE3Nzc0NDU1OTQiLCJhcHBfZGF0YSI6eyJzaG93X290cF9tYW5kYXRlX2Jhbm5lciI6ZmFsc2UsImlzX2ZlZGVyYXRpb25fZGlzYWxsb3dlZCI6ZmFsc2UsImRpc2FsbG93ZWRfZmVkZXJhdGlvbl9pZHMiOltdfSwic3RhdGVfdG9rZW5faWQiOiI2ZGQ1NzUyNS1iNTFlLTQzZGQtODdmMy0xNGFjOWQ4NTJjOWUiLCJ1c2VyX3R5cGUiOjAsInVhbSI6MSwiaXBtIjpbMiwwLDAsMCwwLDJdfQ.dOKzGl37C4pC_cQBbZsoN9h1Rze0IlpRbkzyofM6ewYnITvDUb2EFcGRjlvq_ukZHuC61rYkDFGpxWqlkXKqrZp7Q2Gr3VkEb61bb5r998mbj1qB30P2ZVPRBglzF9W_bwUmUCLznDUcHf72KPk8HzY55su9Fud3GuQhRap4sanAhkHw5gj-EsRE-qXaG9FXT-3TzPQa-UFh_Wt6zMikDD84tXOFz0y5Cay8cfCxfgDfAFqm3GUaGZJInhDDfLL8OpjuupRwAuWdlyeMiCsfiTcSe9-g2XPLvEUDLflYd62eiaBEYMgss5oVBWlITnIr_tv869nfafOYj_d4lrFkpg"
ask_open "NEG hjemmeside" "https://www.neg.dk/"

printf "\n== SketchUp ==\n"
printf "Hvis SketchUp ikke er installeret, brug SketchUp-linket og følg skolens Trimble-flow fra skolemailen.\n"
printf "Dashboard-genvej på Mac holdes manuel i første version: læg mappen et fast sted og åbn start.html.\n"

printf "\nFærdig. Tryk Enter for at lukke.\n"
read unused

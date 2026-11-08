#!/usr/bin/env python3
"""Fetch GitHub contribution history and write contributions.json for the portfolio.

Usage:  python3 update_contributions.py [output_path]
Token:  GITHUB_TOKEN env var, or the file ~/.github_token (chmod 600). No scopes needed.
Output: output_path, default ./contributions.json (same shape as github-contributions-api).
"""
import json, os, sys, urllib.error, urllib.request
from datetime import datetime, timezone

USER = "08abhinav"
TOKEN = os.environ.get("GITHUB_TOKEN") or open(os.path.expanduser("~/.github_token")).read().strip()
OUT = sys.argv[1] if len(sys.argv) > 1 else "contributions.json"
LEVELS = {"NONE": 0, "FIRST_QUARTILE": 1, "SECOND_QUARTILE": 2, "THIRD_QUARTILE": 3, "FOURTH_QUARTILE": 4}
QUERY = """query($login:String!,$from:DateTime!,$to:DateTime!){
  user(login:$login){createdAt contributionsCollection(from:$from,to:$to){
    contributionCalendar{weeks{contributionDays{date contributionCount contributionLevel}}}}}}"""

def gql(frm, to):
    body = json.dumps({"query": QUERY, "variables": {"login": USER, "from": frm, "to": to}}).encode()
    req = urllib.request.Request("https://api.github.com/graphql", data=body, headers={
        "Authorization": "Bearer " + TOKEN, "Content-Type": "application/json", "User-Agent": "portfolio-updater"})
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            d = json.load(r)
    except urllib.error.HTTPError as e:
        if e.code == 401:
            sys.exit("401 Unauthorized: GitHub rejected the token (wrong, expired, revoked, or has extra characters/quotes). Generate a new one.")
        raise
    if "errors" in d:
        sys.exit("GitHub error: %s" % d["errors"])
    return d["data"]["user"]

now = datetime.now(timezone.utc)
fmt = "%Y-%m-%dT%H:%M:%SZ"
created = datetime.fromisoformat(gql("%d-01-01T00:00:00Z" % now.year, now.strftime(fmt))["createdAt"].replace("Z", "+00:00"))

days = {}
for year in range(created.year, now.year + 1):
    frm = "%d-01-01T00:00:00Z" % year
    to = now.strftime(fmt) if year == now.year else "%d-12-31T23:59:59Z" % year
    for week in gql(frm, to)["contributionsCollection"]["contributionCalendar"]["weeks"]:
        for d in week["contributionDays"]:
            c = {"date": d["date"], "count": d["contributionCount"], "level": LEVELS[d["contributionLevel"]]}
            if d["date"] not in days or c["count"] > days[d["date"]]["count"]:  # windows overlap at year edges
                days[d["date"]] = c

contribs = [days[k] for k in sorted(days)]
total = {}
for c in contribs:
    total[c["date"][:4]] = total.get(c["date"][:4], 0) + c["count"]

tmp = OUT + ".tmp"
with open(tmp, "w") as f:
    json.dump({"generated": now.strftime(fmt), "total": total, "contributions": contribs}, f)
os.chmod(tmp, 0o644)
os.replace(tmp, OUT)  # atomic, so visitors never see a half-written file
print("wrote", len(contribs), "days")

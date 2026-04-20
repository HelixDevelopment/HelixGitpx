.PHONY: baseline git ws ai all soak

K6 ?= k6
OUT ?= /tmp/k6-out

baseline:
	$(K6) run --out json=$(OUT)-baseline.json scenarios/api_baseline.js

git:
	$(K6) run --out json=$(OUT)-git.json scenarios/git_push_pull.js

ws:
	$(K6) run --out json=$(OUT)-ws.json scenarios/websocket_fanout.js

ai:
	$(K6) run --out json=$(OUT)-ai.json scenarios/ai_chat.js

all: baseline git ws ai

soak:
	argo submit ../../tools/perf/soak-7d.yaml --watch

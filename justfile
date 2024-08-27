protoc:
	protoc --proto_path=proto proto/*.proto --go_out=./ --go-grpc_out=.

test: test-cleanup test-setup test-run

test-run:
  TESTCONTAINERS_RYUK_DISABLED=true go test ./... -v

test-setup:
  scripts/migrate.test.sh

test-cleanup:
  scripts/cleanup.test.sh

air DIR:
  docker compose -f compose.yml -f air.compose.yml {{DIR}}{{ if DIR == "up" { " --detach" } else { "" } }}

sqlite:
  docker run --rm -it -v template_user_db:/var/db -w /var/db keinos/sqlite3 sqlite3 /var/db/user.db

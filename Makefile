.PHONY: migrate-up

# Run all new migrations (only migrations not yet applied).
migrate-up:
	go run . migrate-up

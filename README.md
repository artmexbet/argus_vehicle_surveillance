# argus_vehicle_surveillance
## Generating docs
```bash
cd microservices/gateway
swag init --parseInternal --pd -d "cmd/app" --parseFuncBody --pdl 3 --parseVendor
```
Then docs will be generated in `docs` folder.
You can see it in `http://localhost:8080/swagger`
## Running the project
```bash
docker-compose up
```

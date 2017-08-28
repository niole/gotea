# gotea
Aggregates information on different types of tea from a set of sites.

## install

In root directory, run `sh setup.sh`. Follow instructions, keep pressing enter, et c.. Have a MySQL data base running on your machine.


Make sure you have a `secrets.py` and a `secrets.go` file. These hold relevant information for the data base.

## run in dev (no prod setup)

In root directory run `go run db.go secrets.go main.go`, `export FLASK_APP=route.py; flask run --reload`.

In `viewtea/` run `npm start`.

The app is running on `localhost:8080`, the python service is running on `localhost:5000`.

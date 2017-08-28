# gotea
Aggregates information on different types of tea from a set of sites.

## install

In root directory, run `sh setup.sh`. Follow instructions, keep pressing enter, et c.. Have a MySQL data base running on your machine.


Make sure you have a `secrets.py` and a `secrets.go` file. These hold relevant information for the data base.

## run in dev (no prod setup)

In root directory run `go run db.go secrets.go main.go`. This script crawls sites that sell tea and aggregates information. It is a webcrawler, so it will take some time to complete. It populates the data base with information that is used to reccommend teas to users based on their text queries. It is best to let it finish before moving on.

To run the python service, run the following in the root directory: `export FLASK_APP=route.py; flask run --reload`.

In order to serve the GUI through which the user can query for tea, in `viewtea/` run `npm start`.

The user facing application is running on `localhost:8080`, the python service is running on `localhost:5000`.

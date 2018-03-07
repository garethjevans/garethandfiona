# kirstenandchris website
A website using templating, gorilla/mux router and newrelic monitoring. Which utilises the [New Age](https://startbootstrap.com/template-overviews/new-age/) Start Bootstrap theme, with thanks!

## Build the docker container
From the root of the project run:
```
docker build -t kirstenandchris .
```

## Run the docker container
From the root of the project run:
```
docker run -d --name kirstenandchris \
              --publish 8080:8080 \
              --env NEWRELIC_LICENSE_KEY="<YOUR_LICENSE_KEY_HERE>" \
              --env NEWRELIC_APP_NAME="kirstenandchris" \
              kirstenandchris
```

## Run go app
From the root of the project run:
```
make
go run main.go
```

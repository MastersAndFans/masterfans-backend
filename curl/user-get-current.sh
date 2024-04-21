curl --header "Content-Type: application/json" \
  --request GET \
  --cookie "auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJpc3MiOiJtYXN0ZXJmYW5zIiwiZXhwIjoxNzEzODA1MDM1LCJpYXQiOjE3MTM3MTg2MzV9.rglGr7u_AwzIVJJ1umq4v4j7sp3ios0x8MyQlYszhqc" \
  http://localhost:5000/api/current-user

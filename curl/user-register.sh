curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"email":"john@gmail.com","password":"123456789Aa@", "repeat_pass":"123456789Aa@", "name":"John", "surname":"Johnny", "birth_date":"2020-01-01"}' \
  http://localhost:5000/api/register

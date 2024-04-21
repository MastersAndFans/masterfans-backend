curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"email":"antanas.matinis@gmail.com","password":"Antanelis5@", "repeat_pass":"Antanelis5@", "name":"John", "surname":"Johnny", "birth_date":"2020-01-01", "phone_number":"+37061235648"}' \
  http://localhost:5000/api/register

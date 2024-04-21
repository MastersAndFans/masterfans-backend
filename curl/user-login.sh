curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"email":"antanas.matinis@gmail.com","password":"Antanelis5@"}' \
  --cookie-jar "curl-login-cookies.txt" \
  http://localhost:5000/api/login

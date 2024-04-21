curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"proposer_id":1, "starting_price":1000, "start_date":"2019-01-24", "end_date":"2020-04-01", "active":true, "title":"Ieskomas interjero dizaineris!", "description":"Reikia suprojektuoti namo kambariu dizaina", "city":"Vilnius", "category":5}' \
  http://localhost:5000/api/auction

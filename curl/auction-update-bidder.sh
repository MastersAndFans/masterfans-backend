curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"auction_id":1, "user_id":2, "bid":900}' \
  http://localhost:5000/api/auction/update-bidder

#Create Account
curl -v --header "Content-Type: application/json" --request POST --data '{"first_name":"Bob","last_name":"White", "balance": 100}' http://localhost:8080/account
#Get balance
curl -v http://localhost:8080/account/{id}/balance
#Move cash
curl -v --header "Content-Type: application/json" --request PUT --data '{"sender_id":1,"recipient_id":2, "Amount": 50}' http://localhost:8080/paymentHandler

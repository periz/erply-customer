
**Swagger URL** : http://localhost:8080/swagger/

Sample Login details:

  Username: persrc@gmail.com
  
  Password: Aa@1234567
  
  Customer Code: 602569
  
  
 Use /auth to get session key and use the sessionkey for all the subsequent add/get customer api calls
 
 /add uses SaveCustomerBulk to create new customer
 
  **Sample Input** : Send array of customers.Customer details to create new customers
 
 /get uses  GetCustomers to get matching customers
 
  **Sample Input** :  { "customerID": "28" }
 
 
 Instructions to build and run the middleware:
 
  -  git clone https://github.com/periz/erply-customer.git
 
  - cd erply-customer
 
  -  go build 
 
  -  ./client
  
  - Access http://localhost:8080/swagger/

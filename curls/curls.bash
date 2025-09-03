curl http://localhost:8080/tasks

#Создание Task
curl -X POST -H "Content-Type: application/json" -d "{\"title\":\"Купить молоко\"}" http://localhost:8080/tasks

#Удаление Task
curl -X DELETE http://localhost:8080/tasks/123

#Частичное изменение Task'а (выборочно)
curl -X PATCH -H "Content-Type: application/json" -d "{\"title\":\"Помыться\"}" http://localhost:8080/tasks/4
curl -X PATCH -H "Content-Type: application/json" -d "{\"title\":\"new2\", \"completed\":true}" http://localhost:8080/tasks/4


curl -X GET "http://localhost:8080/tasks?page=2&limit=10"

curl -X GET "http://localhost:8080/tasks?page=1&limit=5" -H "Content-Type: application/json"



curl -X POST http://localhost:8080/register \  -H "Content-Type: application/json" \  -d "{\"email\": \"user@example.com\", \"password\": \"password123\"}"

curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d "{\"email\": \"user@example.com\", \"password\": \"password123\"}"
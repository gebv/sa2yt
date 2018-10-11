# Запуск
1. необходимо создать .env из env.example заполнив все переменные реальными данными
2. создать билд приложения
  $ GOOS=linux buffalo build -o build/heroku
3. создать коммит
  $ git add .
  $ git commit -m "Commit message"
4. задеплоить на heroku
  $ git push heroku master

# Запуск
1. необходимо создать .env из env.example заполнив все переменные реальными данными
2. создать билд приложения
  $ GOOS=linux buffalo build -o build/heroku
3. создать коммит
  $ git add .
  $ git commit -m "Commit message"
4. задеплоить на heroku
  $ git push heroku master

# Создание приложения в Slack https://api.slack.com/apps
## Добавить actions:
1. В правом меню Features -> Interactive Components, включить этот блок.
2. В Request URL прописать https://example.com/slack_actions
3. В Actions добавить два элемента с CallbackID new_task(диалоговое окно для создания задачи) и
  new_comment(диалоговое окно для добавления коментария)
4. В Options Load URL прописать https://example.com/slack_options

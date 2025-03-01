# vinted_go
DIscord bot to notify about newly added items to Vinted

- discord bot ui for user to ignore seller_currency
- setup database SQLC??
- the SQLC should contain ids of posted items on discord and also all of the other json info
    - 3 tables
      - ids of posted items
      - ids, all the other info about items
      - urls of user's chosen vinted urls + timedata + currency_seller
- once it commits in db, the new message should be sent to discord channel embedded
- 

so at first, the bot starts  running
then the user should add the url in message by some command, maybe !watchUrl url:https://vinted.sk/... seller_currency=CZK seller_currency=EUR (to filter out polish ones)
the agent saves the Url to db and could every one minute check for items on the url by <creating thread?channel?>
everytime he does that, he could check the IDs of items against the database of already sent and if there's nonmatching id, then send that ID's url of item to discord so it can be seen that it is a new item never added before.
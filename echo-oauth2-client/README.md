
OAuth2 for authenticate user
`./echo-oauth2-client -authurl 'https://discord.com/api/oauth2/authorize' -client-id 'XXXXXXXXXXXXXXXXXX' -client-secret 'XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX' -tokenurl 'https://discord.com/api/oauth2/token' -server-id '719506617686687764' -scope 'identify'`

OAuth2 for invite bot user to your server
`./echo-oauth2-client -authurl 'https://discord.com/api/oauth2/authorize' -client-id 'XXXXXXXXXXXXXXXXXX' -client-secret 'XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX' -tokenurl 'https://discord.com/api/oauth2/token' -server-id '719506617686687764' -scope 'bot' -permission=X`

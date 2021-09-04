# cookit-backend

It is a browser-based application aimed at managing recipes. 
This project contains the backend part written by me. 
The application is user based.
- Link to the app:
https://cookit0.herokuapp.com
- API specification available on folder **/api_documentation**


## Authors

- [@daniorocket](https://www.github.com/daniorocket)
- [@rafal21-ops](https://www.github.com/rafal21-ops)

  
## Tech Stack

**Frontend:** Angular

**Backend:** Go 1.15

  
## Features

- Deployment on Heroku
I decided to deploy this app on heroku, because heroku makes it easy. 
- MongoDB
I choosed non-relational database mongoDB, because its interesting alternative for classing relationals DB's like SQL, Postresql etc. Mongo is really fast database, based on documents and collections.
- JWT Authorization
I decided to use JSON Web Token for authorize app users. This solution is more secure than the session-based examples like in PHP.
- UUID
This project contains ID's which are not increment values, but random strings 128-bit length. This solution prevents ID's duplication.



  
## Run Locally

Clone the project

```bash
  git clone https://github.com/Daniorocket/cookit-backend.git
```

Go to the project directory

```bash
  cd path-to-my-project
```

Install heroku locally

```bash
  https://devcenter.heroku.com/articles/getting-started-with-go
```
Start the server

```bash
  heroku local web
```


  
## Environment Variables

To run this project, you will need to add the following environment variables.

`PORT` - Defines numbert of port used in app.

`MONGODB_URI` -  Defines your database server, comes from mongo Atlas.

`DBName` - Defines name of used database

`Email_HOST` - Defines name of used email host. 

`EMAIL_LOGIN` - Defines email login. 

`EMAIL_PASSWORD` - Defines email password. 

`EMAIL_PORT` - Defines email port. 
  
`JWT_KEY` - Defines jwt key, used for creating and verify JSON Web Tokens. 
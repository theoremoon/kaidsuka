require "kemal"
require "kemal-session"
require "json"
require "crypto/bcrypt/password"
require "sqlite3"

Kemal::Session.config do |config|
  config.secret = "aiueo"
end

class User
  JSON.mapping(
    username: String,
    password: String,
  )
end

before_all do |env|
  env.response.headers.add("Access-Control-Allow-Origin", "http://localhost:1234")
  env.response.headers.add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
  env.response.headers.add("Access-Control-Allow-Credentials", "true")
  env.response.headers.add("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Origin, Accept")
  env.response.headers.add("Access-Control-Max-Age", "86400")
end

get "/user" do |env|
  username = env.session.string?("username")
  if username

    env.response.content_type = "application/json"
    {username: username}.to_json
  else
    halt env, status_code: 400, response: "do not logged in"
  end
end

post "/logout" do |env|
  env.session.destroy
  ""
end

post "/login" do |env|
  user = User.from_json env.request.body.not_nil!
  if !user.username || !user.password
    halt env, status_code: 400, response: "Invalid Request"
  end

  DB.open "sqlite3://./database.sqlite" do |db|
    password = db.query_one "select password from users where username = ?", user.username, &.read(String)
    if !Crypto::Bcrypt::Password.new(password).verify( user.password)
      halt env, status_code: 400, response: "Wrong username/password"
    end
  end

  env.session.string("username", user.username)
  ""
end

post "/register" do |env|
  user = User.from_json env.request.body.not_nil!
  if !user.username || !user.password
    halt env, status_code: 400, response: "Invalid Request"
  end

  DB.open "sqlite3://./database.sqlite" do |db|
    db.exec "insert into users(username, password) values (?, ?)", user.username, Crypto::Bcrypt::Password.create(user.password).to_s
  end
  ""
end

options "/*" do |env|
end

Kemal.run do |config|
  server = config.server.not_nil!
  server.bind_tcp "0.0.0.0", 5000, reuse_port: true
end

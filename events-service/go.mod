module github.com/chandanghosh/events-management/events-service

go 1.16

require (
	github.com/chandanghosh/events-management/contracts v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/streadway/amqp v1.0.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/yaml.v2 v2.4.0 // indirect

)

replace github.com/chandanghosh/events-management/contracts => ../contracts

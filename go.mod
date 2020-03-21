module github.com/BabyEngine/Backend

go 1.14

require (
	github.com/BabyEngine/UnityConnector v0.0.0-00010101000000-000000000000
	github.com/DGHeroin/golua v1.0.1
	github.com/boltdb/bolt v1.3.1
	github.com/gorilla/websocket v1.4.1
)

replace github.com/BabyEngine/UnityConnector => ../UnityConnector

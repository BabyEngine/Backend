PBDefine =  [[
syntax = "proto2";

message GameMessage {
  required string action = 1;
  optional bytes  data   = 2;
}

message RespMessage {
  required int64 code = 1;
  required string msg   = 2;
}

message RequestLoginMessage {
  required string token = 1;
}

message ListRoomResp {
    repeated ListRoomInfo rooms = 1;
}

message ListRoomInfo {
    required int64 id = 1;
    repeated int64 players = 2;
}
message ResponseCommon {
  required int64  code = 1;
  required string msg  = 2;
}

message JoinRoom {
    required int64 roomId   = 1;
    required int64 playerId = 2;
}

message PlayingMessage {
    required int64 cmd   = 1;
}


]]
